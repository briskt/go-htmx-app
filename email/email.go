package email

import (
	"context"
	"fmt"
	"strings"

	"github.com/briskt/go-htmx-app/email/message"
	"github.com/briskt/go-htmx-app/log"
)

type Service interface {
	Send(ctx context.Context, msg message.Message) error
}

// SendBatch sends a batch of email messages using the given service.
func SendBatch(ctx context.Context, svc Service, messages []message.Message) error {
	if len(messages) == 0 {
		return nil
	}

	var err error
	var errors []string
	for _, msg := range messages {
		if err = svc.Send(ctx, msg); err != nil {
			errors = append(errors, err.Error())
		}
	}
	if len(errors) == len(messages) {
		return fmt.Errorf("failed to send all %d messages in the batch: %w", len(messages), err)
	}
	if len(errors) > 0 {
		log.Errorf("failed to send one or more messages: %s", strings.Join(errors, "; "))
	}
	return nil
}

func addressWithName(name, address string) string {
	if name == "" {
		return address
	}
	return fmt.Sprintf("%s <%s>", name, address)
}
