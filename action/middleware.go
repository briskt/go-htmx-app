package action

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/briskt/go-htmx-app/api"
	"github.com/briskt/go-htmx-app/app"
	"github.com/briskt/go-htmx-app/core"
	"github.com/briskt/go-htmx-app/email"
	"github.com/briskt/go-htmx-app/log"
)

// authenticationMiddleware supports both bearer token and session-based user token authentication. If a valid bearer
// token is found, the token-auth flag is set in context. If a valid user session is present, the user record is added
// to context.
func authenticationMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if authnSkipper(c) {
				return next(c)
			}

			if hasValidBearerToken(c.Request().Header) {
				app.ContextKeyTokenAuth.Set(c, true)
				return next(c)
			}

			token, err := sessionGetString(c, AccessTokenSessionKey)
			if err != nil {
				log.Debugf("failed to retrieve access token from session: %s", err)
				return next(c)
			}

			if token == "" {
				err = errors.New("no access token provided")
				return api.NewAppError(err, api.ErrorNotAuthenticated, http.StatusUnauthorized)
			}

			user, err := core.FindUserByToken(toCtx(c), Tx(c), token)
			if err != nil {
				return api.NewAppError(err, api.ErrorNotAuthenticated, http.StatusUnauthorized)
			}

			log.SetUser(toCtx(c), user.EmployeeID, email.MaskString(user.GetDisplayName()),
				email.MaskEmail(user.GetEmail()))

			app.ContextKeyCurrentUser.Set(c, user)
			return next(c)
		}
	}
}

// isTokenAuth returns true if the token-auth flag was set in the context by the middleware
func isTokenAuth(c echo.Context) bool {
	if _, ok := app.ContextKeyTokenAuth.Get(c).(bool); !ok {
		return false
	}
	return true
}

// hasValidBearerToken compares a provided token against the know list of tokens and returns true if there's a match
func hasValidBearerToken(h http.Header) bool {
	authHeader := h.Get(echo.HeaderAuthorization)
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && parts[0] == "Bearer" && slices.Contains(app.Env.APIAccessKeys, parts[1]) {
		log.WithFields(log.Fields{"tokenPrefix": parts[1][0:3]}).Debug("authenticated with bearer token")
		return true
	}
	return false
}

// authnSkipper is the skipper for the authentication middleware
func authnSkipper(c echo.Context) bool {
	method := c.Request().Method
	if method == http.MethodOptions {
		return true
	}
	path := c.Request().URL.Path
	if strings.HasPrefix(path, "/assets") {
		log.Tracef("authn skipping %s %q", method, path)
		return true
	}
	// TODO: expand this to differentiate token APIs from user APIs
	skipURLs := []struct{ method, expression string }{
		{"POST", "/auth/callback"},
		{"GET", "/auth/login"},
		{"GET", "/auth/logout"},
		{"GET", "/auth/logout-callback"},
		{"GET", "/robots.txt"},
		{"GET", "/site/status"},
	}
	for _, skip := range skipURLs {
		if match, _ := regexp.MatchString(skip.expression, path); match && skip.method == method {
			log.Tracef("authn skipping %s %s", method, path)
			return true
		}
	}
	log.Tracef("authn not skipping %s %s", method, path)
	return false
}

// transactionMiddleware starts a database transaction and rolls back if status is 400 or higher
func transactionMiddleware(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tx, err := db.Begin()
			if err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			app.ContextKeyTx.Set(c, tx)

			if err = next(c); err != nil {
				_ = tx.Rollback()
				return err
			}

			if c.Response().Status > 399 {
				return tx.Rollback()
			}
			return tx.Commit()
		}
	}
}

func requestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		HandleError: true,
		LogError:    true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			err := values.Error
			if err == nil {
				log.WithFields(log.Fields{
					"employeeID": CurrentUser(c).EmployeeID,
					"method":     c.Request().Method,
					"uri":        c.Request().RequestURI,
				}).Info("request")
				return nil
			}

			fields := log.Fields{}
			var appErr *api.AppError
			if errors.As(err, &appErr) {
				fields = appErr.Extras
			}

			fields["error"] = err.Error()
			logEntry := log.WithFields(fields).WithContext(c.Request().Context())

			switch values.Status {
			case http.StatusUnauthorized, http.StatusBadRequest:
				logEntry.Warning("http error")
			default:
				logEntry.Error("http error")
			}
			return nil
		},
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Request().URL.Path, "/assets")
		},
	})
}
