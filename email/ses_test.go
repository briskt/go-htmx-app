package email

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/briskt/go-htmx-app/email/message"
)

// TestSendRaw can be used in a local environment for development. Add AWS credentials according to AWS
// SDK documentation and change the "To" and "From" email addresses to valid addresses.
func TestSendRaw(t *testing.T) {
	t.Skip("only for use in local environment if configured with credentials")

	data, err := rawEmail(
		"me@example.com",
		"from@example.com",
		"test subject",
		`<h4>body</h4><img src="cid:logo"><p>End of body</p>`,
		map[string]string{"logo": "logo.png"},
		&files)
	require.NoError(t, err)

	ses, err := NewSES("me@example.com")
	require.NoError(t, err)

	err = ses.(SES).SendRaw(context.Background(), data)
	require.NoError(t, err)
}

// TestSendSES can be used in a local environment for development. Add credentials according to AWS
// SDK documentation and change the "To" and "From" email addresses to valid addresses.
func TestSendSES(t *testing.T) {
	t.Skip("only for use in local environment if configured with credentials")

	service, err := NewSES("me@example.com")
	require.NoError(t, err)

	params := message.Params{
		Template: message.Welcome,
		From:     message.NewAddress("name", "from@example.com"),
		To:       message.NewAddress("name", "to@example.com"),
		Fields:   map[string]any{},
		Images:   map[string]string{},
	}
	msg, _ := message.New(params)

	err = service.Send(context.Background(), msg)
	require.NoError(t, err)
}
