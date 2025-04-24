package action

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/briskt/go-htmx-app/app"
	"github.com/briskt/go-htmx-app/data"
)

var enabled bool

type ProfileView struct {
	DisplayName   string
	Enabled       bool
	HelpCenterURL template.URL
	AppName       string
	LastLogin     string
	UserID        string
	Username      string
}

// home renders the home page
func home(c echo.Context) error {
	user := CurrentUser(c)
	if user.ID == 0 {
		return c.Redirect(http.StatusFound, "/auth/login")
	}

	return renderHome(c, user)
}

// renderHome renders the "home.gohtml" template
func renderHome(c echo.Context, user data.User) error {
	profileData := ProfileView{
		AppName:       app.Env.AppName,
		DisplayName:   user.GetDisplayName(),
		Enabled:       enabled,
		HelpCenterURL: template.URL(app.Env.HelpCenterURL),
		LastLogin:     formatDate(user.LastLoginAt),
		Username:      user.Username,
		UserID:        strconv.Itoa(int(user.ID)),
	}

	return c.Render(http.StatusOK, "home.gohtml", profileData)
}

// card responds to the button on "card"
func card(c echo.Context) error {
	enabled = !enabled
	return c.Render(http.StatusOK, "card.gohtml", map[string]any{"Enabled": enabled})
}

// formatNullDate returns a long-form, user-friendly date string from a valid sql.NullTime. If invalid, it returns "-"
func formatNullDate(d sql.NullTime) string {
	if d.Valid {
		return formatDate(d.Time)
	}
	return "-"
}

// formatDate returns a long-form, user-friendly date string from a time.Time
func formatDate(d time.Time) string {
	return d.Format("Monday, January 2, 2006")
}
