package message_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/briskt/go-htmx-app/email/message"
)

func TestMethodVerify(t *testing.T) {
	fields := message.Fields{
		"AppName":        "Test",
		"DisplayName":    "X Smith",
		"EmailSignature": "this is the message signature",
		"HelpCenterURL":  "http://example.com",
		"SupportEmail":   "support@example.com",
		"SupportName":    "Support Team",
		"ToAddress":      "foo@example.com",
	}
	params := message.Params{
		Template: message.Welcome,
		From:     message.NewAddress("no_reply@example.com", "No Reply"),
		To:       message.NewAddress("X Smith", "foo@example.com"),
		Fields:   fields,
		Images:   map[string]string{"logo": "logo.png"},
	}
	msg, err := message.New(params)
	require.NoError(t, err)
	require.Equal(t, "X Smith <foo@example.com>", msg.To())
	require.Equal(t, "no_reply@example.com <No Reply>", msg.From())
	require.Equal(t, "Important information about your Test account", msg.Subject())
	require.Contains(t, msg.Body(), "message signature")
	require.Equal(t, map[string]string{"logo": "logo.png"}, msg.Images())
}
