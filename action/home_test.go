package action

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/briskt/go-htmx-app/data"
)

func (s *Suite) TestHome() {
	user, err := data.CreateUser(s.ctx, s.db, data.UserCreateInput{})
	s.NoError(err)
	saveToken(s.db, int(user.ID), testToken)

	_, status := s.request("GET", "/", "invalid", nil)
	s.Equal(http.StatusUnauthorized, status)

	response, status := s.request("GET", "/", testToken, nil)
	s.Equal(status, http.StatusOK)
	s.Contains(string(response), "<h1>My Go HTMX App</h1>")
}

func (s *Suite) TestFormatNullDate() {
	date := time.Date(2024, 1, 1, 1, 1, 1, 0, time.UTC)

	got := formatNullDate(sql.NullTime{Valid: true, Time: date})
	s.Equal("Monday, January 1, 2024", got)

	got = formatNullDate(sql.NullTime{Valid: false, Time: date})
	s.Equal("-", got)
}

func (s *Suite) TestFormatDate() {
	got := formatDate(time.Date(2024, 1, 1, 1, 1, 1, 0, time.UTC))
	s.Equal("Monday, January 1, 2024", got)
}
