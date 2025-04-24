package email_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/briskt/go-htmx-app/app"
	"github.com/briskt/go-htmx-app/email"
	"github.com/briskt/go-htmx-app/email/message"
)

// TestSendMailgun can be used in a local environment for development. Add credentials to the appropriate
// environment variables, and change the "To" and "From" email addresses to valid addresses.
func TestSendMailgun(t *testing.T) {
	t.Skip("only for use in local environment if configured with credentials")

	service := email.NewMailgun(email.MailgunConfig{
		Domain:       app.Env.MailgunDomain,
		PrivateKey:   app.Env.MailgunAPIKey,
		SandboxEmail: "me@example.com",
	})

	params := message.Params{
		Template: message.Welcome,
		From:     message.NewAddress("name", "from@example.com"),
		To:       message.NewAddress("name", "to@example.com"),
		Fields:   map[string]any{},
		Images:   map[string]string{},
	}
	msg, _ := message.New(params)

	err := service.Send(context.Background(), msg)
	require.NoError(t, err)
}
