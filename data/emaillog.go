package data

import (
	"context"
	"strings"

	"github.com/briskt/go-htmx-app/data/sqlc"
	"github.com/briskt/go-htmx-app/log"
)

func CreateEmailLog(ctx context.Context, tx sqlc.DBTX, userID int, template string) error {
	template = strings.ReplaceAll(template, "_", "-")
	log.WithFields(log.Fields{"userID": userID, "template": template}).Debug("creating email log")
	return q(tx).CreateEmailLog(ctx, int32(userID), template)
}

// HasReceivedMessageRecently searches the email log for the given template and returns true if at least one such
// message has been sent to the user recently
func (u User) HasReceivedMessageRecently(ctx context.Context, tx sqlc.DBTX, template string) bool {
	template = strings.ReplaceAll(template, "_", "-")
	n, err := q(tx).CountRecentEmails(ctx, u.ID, template)
	if err != nil {
		log.Errorf("failed to count recent emails for template %s: %v", template, err)
	}
	log.Debugf("userID: %d , template: %s", u.ID, template)
	return n > 0
}
