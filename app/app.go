package app

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"

	"github.com/briskt/go-htmx-app/email"
	"github.com/briskt/go-htmx-app/log"
)

const (
	EnvDevelopment = "dev"  // EnvDevelopment enables various debugging aids
	EnvTest        = "test" // EnvTest is for automated tests, during which some things are disabled
)

type ContextKey string

const (
	ContextKeyCurrentUser = ContextKey("current_user")
	ContextKeyTx          = ContextKey("tx")
	ContextKeyTokenAuth   = ContextKey("token_auth")
)

func (c ContextKey) Set(ctx echo.Context, value any) {
	log.Tracef("setting %s in echo context", string(c))
	ctx.Set(string(c), value)
}

func (c ContextKey) Get(ctx echo.Context) any {
	return ctx.Get(string(c))
}

const AccessTokenLifetime = 30 * time.Minute

func init() {
	readEnv()
}

func OpenDatabase() (*sql.DB, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s?sslmode=%s",
		Env.PostgresUser, Env.PostgresPassword, Env.PostgresHost, Env.PostgresDB, Env.PostgresSSLMode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, err
}

func NewEmailService() (email.Service, error) {
	emailService := email.NewFake()
	switch Env.EmailService {
	case "mailgun":
		log.WithFields(log.Fields{"domain": Env.MailgunDomain}).Info("using Mailgun")
		emailService = email.NewMailgun(email.MailgunConfig{
			Domain:       Env.MailgunDomain,
			PrivateKey:   Env.MailgunAPIKey,
			SandboxEmail: Env.SandboxEmail,
		})
	case "ses":
		log.WithFields(log.Fields{"region": Env.AWSRegion, "accessKeyID": Env.AWSAccessKeyID}).
			Infof("using AWS SES")
		var err error
		emailService, err = email.NewSES(Env.SandboxEmail)
		if err != nil {
			return nil, fmt.Errorf("error creating SES email service: %w", err)
		}
	}
	return emailService, nil
}
