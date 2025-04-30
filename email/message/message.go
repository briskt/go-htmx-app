package message

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/briskt/go-htmx-app/email/templates"
)

const (
	Welcome = "welcome"
)

type Message struct {
	from    Address
	to      Address
	subject string
	body    string

	// Images returns a map where the keys are the image tag (cid) and the values are the filenames.
	// Any image that doesn't map to a corresponding `src="cid:tag"` in the body will be omitted from the inline
	// attachments. The included filenames must be in the static.EFS() filesystem. TODO: make this file system more
	// portable by adding it as an input to this package's initialization.
	images map[string]string
}

var bodyTemplates = make(map[string]*template.Template)

type Fields map[string]any

type Params struct {
	Template string
	From     Address
	To       Address
	Fields   Fields
	Images   map[string]string
}

func init() {
	for name := range subjects {
		parsed, err := template.ParseFS(templates.EFS(), "base.gohtml", name+".gohtml")
		if err != nil {
			panic(fmt.Sprintf("failed to parse subject template '%s': %v", name, err))
		}
		bodyTemplates[name] = parsed
	}
}

func New(params Params) (Message, error) {
	body, err := execBodyTemplate(params.Template, params.Fields)
	if err != nil {
		return Message{}, err
	}

	subject, err := execSubjectTemplate(params.Template, params.Fields)
	if err != nil {
		return Message{}, err
	}

	m := Message{
		body:    body,
		from:    params.From,
		to:      params.To,
		subject: subject,
		images:  params.Images,
	}
	return m, nil
}

func execBodyTemplate(templateName string, fields Fields) (string, error) {
	var sb strings.Builder
	err := bodyTemplates[templateName].ExecuteTemplate(&sb, "base.gohtml", fields)
	if err != nil {
		return "", fmt.Errorf("failed to execute template '%s': %w", templateName, err)
	}
	return sb.String(), nil
}

func (m Message) Body() string {
	return m.body
}

func (m Message) From() string {
	return m.from.String()
}

func (m Message) To() string {
	return m.to.String()
}

func (m Message) Subject() string {
	return m.subject
}

func (m Message) Images() map[string]string {
	return m.images
}