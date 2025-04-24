package email

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"strings"

	"jaytaylor.com/html2text"
)

// rawEmail generates a multipart MIME email message with a plain text, html text, and inline image attachments
// The images should be provided as a map, where the keys are the image tag and the values are the filenames.
// Any image that doesn't map to a corresponding `src="cid:tag"` in the body will be omitted from the inline
// attachments. The included filenames must be in the static.EFS() filesystem. TODO: make this file system more
// portable by adding it as an input to this package's initialization.
// The generated email can be summarized as follows:
//
//   - multipart/alternative
//
//   - text/plain
//
//   - multipart/related
//
//   - text/html
//
//   - image/png
//
//     Abbreviated example of the generated email message:
//     From: from@example.com
//     To: to@example.com
//     Subject: subject text
//     Content-Type: multipart/alternative; boundary="boundary_alternative"
//
//     --boundary_alternative
//     Content-Type: text/plain; charset=utf-8
//
//     Plain text body
//     --boundary_alternative
//     Content-type: multipart/related; boundary="boundary_related"
//
//     --boundary_related
//     Content-Type: text/html; charset=utf-8
//
//     HTML body
//     --boundary_related
//     Content-Type: image/png
//     Content-Transfer-Encoding: base64
//     Content-ID: <logo>
//     --boundary_related--
//     --boundary_alternative--
func rawEmail(to, from, subject, body string, images map[string]string, fileSys *embed.FS) ([]byte, error) {
	if images == nil {
		images = map[string]string{}
	}

	tbody, err := html2text.FromString(body)
	if err != nil {
		return nil, fmt.Errorf("error converting html email to plain text: %q", err)
	}

	b := &bytes.Buffer{}

	b.WriteString("From: " + from + "\n")
	b.WriteString("To: " + to + "\n")
	b.WriteString("Subject: " + subject + "\n")
	b.WriteString("MIME-Version: 1.0\n")

	alternativeWriter := multipart.NewWriter(b)
	b.WriteString(`Content-Type: multipart/alternative; type="text/plain"; boundary="` +
		alternativeWriter.Boundary() + `"` + "\n\n")

	w, err := alternativeWriter.CreatePart(textproto.MIMEHeader{
		"Content-Type":        {"text/plain; charset=utf-8"},
		"Content-Disposition": {"inline"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MIME text part: %q", err)
	} else {
		_, _ = fmt.Fprint(w, tbody)
	}

	relatedWriter := multipart.NewWriter(b)
	_, err = alternativeWriter.CreatePart(textproto.MIMEHeader{
		"Content-Type": {`multipart/related; type="text/html"; boundary="` + relatedWriter.Boundary() + `"`},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MIME related part: %q", err)
	}

	w, err = relatedWriter.CreatePart(textproto.MIMEHeader{
		"Content-Type":        {"text/html; charset=utf-8"},
		"Content-Disposition": {"inline"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MIME html part: %q", err)
	} else {
		_, _ = fmt.Fprint(w, body)
	}

	cids := findImagesInBody(body, images)
	if err = attachImages(relatedWriter, b, cids, fileSys); err != nil {
		return nil, fmt.Errorf("failed to attach images: %q", err)
	}

	if err = relatedWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close MIME related part: %q", err)
	}

	if err = alternativeWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close MIME alternative part: %q", err)
	}

	return b.Bytes(), nil
}

func findImagesInBody(body string, images map[string]string) map[string]string {
	imagesFound := map[string]string{}
	for cid, filename := range images {
		if strings.Contains(body, fmt.Sprintf(`src="cid:%s"`, cid)) {
			imagesFound[cid] = filename
		}
	}
	return imagesFound
}

func attachImages(relatedWriter *multipart.Writer, b *bytes.Buffer, images map[string]string, fileSys *embed.FS) error {
	for cid, filename := range images {
		_, err := relatedWriter.CreatePart(textproto.MIMEHeader{
			"Content-Type":              {"image/png"},
			"Content-Disposition":       {"inline"},
			"Content-ID":                {"<" + cid + ">"},
			"Content-Transfer-Encoding": {"base64"},
		})
		if err != nil {
			return fmt.Errorf("failed to create MIME image part for '%s': %w", cid, err)
		}

		if err = encodeFile(fileSys, filename, b); err != nil {
			return fmt.Errorf("failed to encode file for '%s': %w", cid, err)
		}
	}
	return nil
}

// encodeFile reads a file from an embedded file system, base64 encodes it, and streams into a bytes.Buffer
func encodeFile(fileSys *embed.FS, filename string, buffer *bytes.Buffer) error {
	file, err := fileSys.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read '%s' file: %w", filename, err)
	}

	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	_, err = encoder.Write(file)
	if err != nil {
		return fmt.Errorf("failed to encode file '%s': %w", filename, err)
	}

	err = encoder.Close()
	if err != nil {
		return fmt.Errorf("failed to close '%s' base64 encoder: %w", filename, err)
	}

	return nil
}
