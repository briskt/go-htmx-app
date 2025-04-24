package email

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/briskt/go-htmx-app/email/message"
	"github.com/briskt/go-htmx-app/log"
	"github.com/briskt/go-htmx-app/public"
)

// MailgunConfig stores required configuration parameters for the Mailgun SDK
type MailgunConfig struct {
	Domain       string
	PrivateKey   string
	SandboxEmail string
}

// Mailgun sends email using Amazon Simple Email Service (Mailgun)
type Mailgun struct {
	sandbox string

	*mailgun.MailgunImpl
}

func NewMailgun(config MailgunConfig) Service {
	svc := mailgun.NewMailgun(config.Domain, config.PrivateKey)
	return Mailgun{MailgunImpl: svc, sandbox: config.SandboxEmail}
}

// Send a message
func (s Mailgun) Send(ctx context.Context, msg message.Message) error {
	to := msg.To()
	if s.sandbox != "" {
		to = s.sandbox
	}
	log.WithFields(log.Fields{"to": to, "subject": msg.Subject()}).Debug("sending message using Mailgun")

	rawBody, err := rawEmail(to, msg.From(), msg.Subject(), msg.Body(), msg.Images(), public.EFS())
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	mgMsg := s.NewMIMEMessage(io.NopCloser(bytes.NewReader(rawBody)), to) // TODO: make a streaming version of rawEmail for this
	status, id, err := s.MailgunImpl.Send(ctx, mgMsg)
	if err != nil {
		return fmt.Errorf("failed to send using Mailgun: %w", err)
	}

	log.WithFields(log.Fields{"to": to, "status": status, "id": id}).Info("sent message using Mailgun")
	return nil
}
