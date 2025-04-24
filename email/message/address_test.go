package message_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/briskt/go-htmx-app/email/message"
)

func TestAddress(t *testing.T) {
	tests := []struct {
		name string
		addr message.Address
		want string
	}{
		{
			name: "address only",
			addr: message.NewAddress("", "john@example.com"),
			want: "john@example.com",
		},
		{
			name: "both",
			addr: message.NewAddress("John", "john@example.com"),
			want: "John <john@example.com>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.addr.String())
		})
	}
}
