// Package action contains http handlers for HTMX templates
//
// Go-HTMX-App
//
//	Terms Of Service:
//	  There are no TOS at this moment, use at your own risk. We take no responsibility.
//
//	 Schemes: https
//	 Host: localhost
//	 BasePath: /
//	 Version: 0.0.1
//	 License: private/none
//
//	 Consumes:
//	 - application/json
//
//	 Produces:
//	 - application/json
//	 - text/html
//
//	 Security:
//	 - saml2:
//
//	 SecurityDefinitions:
//	 saml2:
//	   type: saml2
//	   authorizationUrl: /auth/login
//	   callbackUrl: /auth/callback
//
// swagger:meta
package action

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/briskt/go-htmx-app/app"
	"github.com/briskt/go-htmx-app/data"
	"github.com/briskt/go-htmx-app/email"
	"github.com/briskt/go-htmx-app/log"
	"github.com/briskt/go-htmx-app/public"
	"github.com/briskt/go-htmx-app/saml"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	*echo.Echo
	store        sessions.Store
	samlProvider *saml.Provider
}

type Config struct {
	DB           *sql.DB
	EmailService email.Service
	Store        sessions.Store
}

var a *App

var errorNotAuthenticated = errors.New("not authenticated")

var emailService = email.NewFake()

func NewApp(config *Config) *App {
	if a == nil {
		a = &App{
			Echo: echo.New(),
		}

		a.Binder = &Binder{}
		a.Debug = app.Env.AppEnv == app.EnvDevelopment
		a.HTTPErrorHandler = customHTTPErrorHandler

		if config.Store != nil {
			a.store = config.Store
		} else {
			a.store = newCookieStore()
		}
		a.Use(session.Middleware(a.store))

		a.Renderer = &public.TemplRenderer{}

		a.Use(requestLogger())

		if app.Env.AppEnv == app.EnvDevelopment {
			a.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
				LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
					// log stack dump on multiple lines
					log.Error(err.Error() + string(stack))
					return nil
				},
			}))
		} else {
			a.Use(middleware.Recover())
		}

		a.Use(transactionMiddleware(config.DB))

		switch app.Env.EmailService {
		case "mailgun":
			log.WithFields(log.Fields{"domain": app.Env.MailgunDomain}).Info("using Mailgun")
			emailService = email.NewMailgun(email.MailgunConfig{
				Domain:       app.Env.MailgunDomain,
				PrivateKey:   app.Env.MailgunAPIKey,
				SandboxEmail: app.Env.SandboxEmail,
			})
		case "ses":
			log.WithFields(log.Fields{"region": app.Env.AWSRegion, "accessKeyID": app.Env.AWSAccessKeyID}).
				Infof("using AWS SES")
			var err error
			emailService, err = email.NewSES(app.Env.SandboxEmail)
			if err != nil {
				log.Fatalf("Error creating SES email service: %v", err)
			}
		}

		a.samlProvider = initSAML()
		a.Use(authenticationMiddleware())

		a.Group("/assets").Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       "assets",
			Browse:     true,
			Filesystem: http.FS(public.EFS()),
		}))

		// Authentication endpoints for UI
		a.GET("/auth/login", a.authLogin)
		a.POST("/auth/callback", a.authCallback)
		a.GET("/auth/logout", a.authLogout)
		a.GET("/auth/logout-callback", a.authLogoutCallback)

		// HTML endpoints for UI
		a.GET("/", home)
		a.PUT("/card", cardItem)

		// for ECS healthcheck
		a.GET("/site/status", siteStatus)

		routes := a.Routes()
		for _, r := range routes {
			log.Tracef("%s %s\n", r.Method, r.Path)
		}
	}

	return a
}

// CurrentUser retrieves the current user from the context.
func CurrentUser(c echo.Context) data.User {
	user, _ := app.ContextKeyCurrentUser.Get(c).(data.User)
	return user
}

// setReturnToInSession looks for a returnTo value in the query string and sets it in the session
func setReturnToInSession(c echo.Context) {
	returnTo := c.QueryParam(ReturnToParam)
	if returnTo != "" {
		if err := sessionSetValue(c, ReturnToSessionKey, returnTo); err != nil {
			log.Errorf("failed to set %s in session: %s", ReturnToSessionKey, err)
		}
	}
}

// getReturnTo gets the returnTo from the session. If not found, it returns the home path.
func getReturnTo(c echo.Context) string {
	returnTo, err := sessionGetString(c, ReturnToSessionKey)
	if err != nil {
		return app.Env.AppURL
	}
	return returnTo
}

// Binder is a custom request body binder. It throws an error if any unknown fields are provided.
type Binder struct{}

// Bind satisfies the Echo Binder interface
func (cb *Binder) Bind(i any, c echo.Context) (err error) {
	dec := json.NewDecoder(c.Request().Body)
	dec.DisallowUnknownFields()
	return dec.Decode(i)
}

// toCtx converts an echo context to a standard context
func toCtx(c echo.Context) context.Context {
	return c.Request().Context()
}

// Tx extracts the database connection from the echo request context
func Tx(c echo.Context) *sql.Tx {
	tx, ok := app.ContextKeyTx.Get(c).(*sql.Tx)
	if !ok {
		log.Fatal("no DB connection in context")
	}
	return tx
}

// newCookieStore creates a new cookie store for session storage
func newCookieStore() sessions.Store {
	store := sessions.NewCookieStore([]byte(app.Env.SessionSecret))

	store.Options.SameSite = http.SameSiteDefaultMode
	store.Options.HttpOnly = true

	if !app.Env.DisableTLS {
		// Cookies will be sent in all contexts, i.e. in responses to both first-party and cross-origin requests.
		// This appears to be required to work with Firefox default cookie blocking setting.
		store.Options.SameSite = http.SameSiteNoneMode
		store.Options.Secure = true
	}

	return store
}

// initSAML initializes the SAML provider
func initSAML() *saml.Provider {
	if app.Env.SamlIdpMetadataURL == "" {
		return nil
	}

	samlConfig := saml.Config{
		SPEntityID:                  app.Env.SamlSpEntityID,
		AudienceURI:                 app.Env.SamlSpEntityID,
		AssertionConsumerServiceURL: app.Env.SamlAssertionConsumerServiceURL,
		SPPublicCert:                app.Env.SamlSpCert,
		SPPrivateKey:                app.Env.SamlSpPrivateKey,
		IDPMetadataURL:              app.Env.SamlIdpMetadataURL,
	}
	p, err := saml.New(samlConfig)
	if err != nil {
		log.Errorf("failed to init SAML Provider: %s", err.Error())
	}
	return p
}
