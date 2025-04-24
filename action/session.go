package action

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/briskt/go-htmx-app/app"
)

const sessionName = "caisson"

func sessionSetValue(c echo.Context, key, value interface{}) error {
	sess, err := getSession(c)
	if err != nil {
		return fmt.Errorf("error getting session in sessionSetValue(): %w", err)
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   !app.Env.DisableTLS,
	}
	sess.Values[key] = value
	if err = sess.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}
	return nil
}

func sessionGetString(c echo.Context, key interface{}) (string, error) {
	if i, err := sessionGetValue(c, key); err != nil {
		return "", err
	} else {
		return i.(string), nil
	}
}

func sessionGetValue(c echo.Context, key interface{}) (interface{}, error) {
	sess, err := getSession(c)
	if err != nil {
		return nil, fmt.Errorf("failed to get value from session: %w", err)
	}
	if v, ok := sess.Values[key]; !ok {
		return nil, fmt.Errorf("key '%s' not found in session", key)
	} else {
		return v, nil
	}
}

func getSession(c echo.Context) (*sessions.Session, error) {
	if sess, err := session.Get(sessionName, c); err != nil {
		return nil, fmt.Errorf("error getting session: %w", err)
	} else {
		return sess, nil
	}
}

func clearSession(c echo.Context) error {
	sess, err := getSession(c)
	if err != nil {
		return fmt.Errorf("failed to get session to clear: %w", err)
	}
	sess.Options.MaxAge = -1
	sess.Options.Secure = false
	sess.Values = make(map[any]any)
	if err = sess.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("error clearing session: %w", err)
	}
	return nil
}
