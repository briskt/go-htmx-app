package email

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"

	"github.com/briskt/go-htmx-app/email/message"
	"github.com/briskt/go-htmx-app/log"
	"github.com/briskt/go-htmx-app/public"
)

// SES sends email using Amazon Simple Email Service (SES)
type SES struct {
	sandbox string

	*ses.Client
}

// NewSES returns an SES service provider for the Service interface
func NewSES(sandboxEmail string) (Service, error) {
	defaultConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load SES configuration: %w", err)
	}
	client := ses.NewFromConfig(defaultConfig)
	return SES{Client: client, sandbox: sandboxEmail}, nil
}

// Send a message
func (s SES) Send(ctx context.Context, msg message.Message) error {
	to := msg.To()
	if s.sandbox != "" {
		to = s.sandbox
	}
	rawBody, err := rawEmail(to, msg.From(), msg.Subject(), msg.Body(), msg.Images(), public.EFS())
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return s.SendRaw(ctx, rawBody)
}

// SendRaw sends a message using SES, given a pre-built raw byte stream
func (s SES) SendRaw(ctx context.Context, data []byte) error {
	input := ses.SendRawEmailInput{
		RawMessage: &types.RawMessage{Data: data},
	}
	output, err := s.Client.SendRawEmail(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to send using SES: %w", err)
	}

	log.WithFields(log.Fields{"messageID": *output.MessageId}).Info("message sent using SES")
	return nil
}
