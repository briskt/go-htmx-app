package action

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"github.com/briskt/go-htmx-app/app"
	"github.com/briskt/go-htmx-app/data"
	"github.com/briskt/go-htmx-app/public/view"
	"github.com/briskt/go-htmx-app/public/view/card"
)

var enabled bool

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
	profileData := data.ProfileView{
		AppName:       app.Env.AppName,
		DisplayName:   user.GetDisplayName(),
		Enabled:       enabled,
		HelpCenterURL: templ.URL(app.Env.HelpCenterURL),
		LastLogin:     formatDate(user.LastLoginAt),
		Username:      user.Username,
		UserID:        strconv.Itoa(int(user.ID)),
	}

	component := view.Home(profileData)
	return c.Render(http.StatusOK, "", component)
}

// card responds to the button on "card"
func cardItem(c echo.Context) error {
	enabled = !enabled
	return c.Render(http.StatusOK, "", card.Card(enabled))
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
