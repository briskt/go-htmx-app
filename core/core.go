package core

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"maps"

	"github.com/briskt/go-htmx-app/app"
	"github.com/briskt/go-htmx-app/data"
	"github.com/briskt/go-htmx-app/email"
	"github.com/briskt/go-htmx-app/email/message"
	"github.com/briskt/go-htmx-app/log"
)

func SendPeriodicMessages(ctx context.Context, tx *sql.Tx, svc email.Service) {
}

func commonEmailImages() map[string]string {
	return map[string]string{
		"logo": "assets/img/logo.png",
	}
}

func commonEmailFields() message.Fields {
	return message.Fields{
		"BrandColor":     template.CSS(app.Env.BrandColor),
		"EmailSignature": template.HTML(app.Env.EmailSignature),
		"HelpCenterURL":  template.URL(app.Env.HelpCenterURL),
		"AppName":        app.Env.AppName,
		"SupportEmail":   app.Env.SupportEmail,
		"SupportName":    app.Env.SupportName,
	}
}

// sendBatch sends an identical message to a list of users, customized only by the user's email and DisplayName
func sendBatch(ctx context.Context, tx *sql.Tx, svc email.Service, users []data.User, params message.Params) (int, error) {
	if len(users) == 0 {
		return 0, nil
	}

	numSent := 0
	for _, user := range users {
		params.Fields["DisplayName"] = user.GetDisplayName()
		params.To = message.NewAddress(user.GetEmail(), user.GetDisplayName())
		err := sendMessage(ctx, tx, svc, int(user.ID), params)
		if err != nil {
			log.Error(err)
			continue
		}
		numSent++
	}
	if numSent == 0 {
		return 0, fmt.Errorf("none of the %d emails in the '%s' batch were sent", len(users), params.Template)
	}

	logEntry := log.WithFields(log.Fields{"numSent": numSent, "numFailed": len(users) - numSent, "template": params.Template})
	if numSent < len(users) {
		logEntry.Error("errors in email batch")
		return numSent, nil // only return error for a complete failure
	}
	logEntry.Info("sent email batch")
	return numSent, nil
}

// sendMessage sends a message described by params, using the email service, to a single user
func sendMessage(ctx context.Context, tx *sql.Tx, svc email.Service, userID int, params message.Params) error {
	params.From = message.NewAddress(app.Env.AppName, app.Env.FromEmail)

	if params.Images == nil {
		params.Images = commonEmailImages()
	} else {
		maps.Copy(params.Images, commonEmailImages())
	}

	maps.Copy(params.Fields, commonEmailFields())

	msg, err := message.New(params)
	if err != nil {
		return fmt.Errorf("failed to create %s message: %w", params.Template, err)
	}

	if err = svc.Send(ctx, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if err = data.CreateEmailLog(ctx, tx, userID, params.Template); err != nil {
		log.Errorf("failed to create email log: %s", err.Error())
	}
	return nil
}
