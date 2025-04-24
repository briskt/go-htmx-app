package email

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/briskt/go-htmx-app/email/message"
	"github.com/briskt/go-htmx-app/log"
	"github.com/briskt/go-htmx-app/public"
)

type FakeEmailService struct {
	sentMessages []FakeMessage
}

type FakeMessage struct {
	Subject, Body, From, To string
}

func NewFake() Service {
	return &FakeEmailService{}
}

func (t *FakeEmailService) Send(_ context.Context, msg message.Message) error {
	to := msg.To()
	from := msg.From()
	subject := msg.Subject()
	body := msg.Body()

	if body == "ERROR" {
		return errors.New("mock error for testing")
	}

	rawMessage, err := rawEmail(to, from, subject, body, msg.Images(), public.EFS())
	if err != nil {
		return fmt.Errorf("failed to MIME encode email body: %w", err)
	}
	t.sentMessages = append(t.sentMessages,
		FakeMessage{
			Subject: subject,
			Body:    string(rawMessage),
			From:    from,
			To:      to,
		})

	log.WithFields(log.Fields{"subject": subject, "to": to}).Info("(fake) message sent")

	err = os.MkdirAll("email"+string(os.PathSeparator)+"tmp", 0o777)
	if err != nil {
		return fmt.Errorf("failed to create tmp folder: %w", err)
	}
	filename := "email/tmp/" + time.Now().Format("2006_01_02_15_04_05.eml")
	if err = os.WriteFile(filename, rawMessage, 0o666); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetNumberOfMessagesSent returns the number of messages sent since initialization or the last call to
// DeleteSentMessages
func (t *FakeEmailService) GetNumberOfMessagesSent() int {
	return len(t.sentMessages)
}

// DeleteSentMessages erases the store of sent messages
func (t *FakeEmailService) DeleteSentMessages() {
	t.sentMessages = []FakeMessage{}
}

func (t *FakeEmailService) GetLastToEmail() string {
	if len(t.sentMessages) == 0 {
		return ""
	}

	return t.sentMessages[len(t.sentMessages)-1].To
}

func (t *FakeEmailService) GetToEmailByIndex(i int) string {
	if len(t.sentMessages) <= i {
		return ""
	}

	return t.sentMessages[i].To
}

func (t *FakeEmailService) GetAllToAddresses() []string {
	emailAddresses := make([]string, len(t.sentMessages))
	for i := range t.sentMessages {
		emailAddresses[i] = t.sentMessages[i].To
	}
	return emailAddresses
}

func (t *FakeEmailService) GetLastBody() string {
	if len(t.sentMessages) == 0 {
		return ""
	}

	return t.sentMessages[len(t.sentMessages)-1].Body
}

func (t *FakeEmailService) GetSentMessages() []FakeMessage {
	return t.sentMessages
}
