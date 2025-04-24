package core

import (
	"context"
	"database/sql"

	"github.com/briskt/go-htmx-app/data"
	"github.com/briskt/go-htmx-app/email"
	"github.com/briskt/go-htmx-app/email/message"
)

func sendWelcomeMessage(ctx context.Context, tx *sql.Tx, svc email.Service, user data.User) error {
	fields := map[string]any{
		"DisplayName": user.GetDisplayName(),
		"Username":    user.Username,
	}
	params := message.Params{
		Template: message.Welcome,
		To:       message.NewAddress(user.GetEmail(), user.GetDisplayName()),
		Fields:   fields,
	}
	return sendMessage(ctx, tx, svc, int(user.ID), params)
}
