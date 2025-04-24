package action

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/briskt/go-htmx-app/app"
	"github.com/briskt/go-htmx-app/core"
	"github.com/briskt/go-htmx-app/data"
	"github.com/briskt/go-htmx-app/log"
)

const testToken = "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0"

func init() {
	log.Init()
}

type Suite struct {
	suite.Suite
	*require.Assertions
	app     *App
	ctx     context.Context
	db      *sql.DB
	session *sessions.Session
}

func (s *Suite) SetupTest() {
	s.session, _ = s.app.store.New(nil, sessionName) // testSessionStore doesn't require an http.Request
	s.Assertions = require.New(s.T())
	data.DestroyTables(s.db)
}

// TestSuite runs the test suite
func TestSuite(t *testing.T) {
	dsn := fmt.Sprintf("postgresql://%s:%s@test_db:5432/%s?sslmode=disable",
		app.Env.PostgresUser, app.Env.PostgresPassword, app.Env.PostgresDB)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	s := &Suite{
		app: NewApp(&Config{
			DB:           db,
			EmailService: emailService,
			Store:        newTestSessionStore(),
		}),
		ctx: context.Background(),
		db:  db,
	}
	s.app.samlProvider = initSAML()
	suite.Run(t, s)
}

// requestResponse submits a test request and captures the response. The provided token is set as the Bearer token (for
// API calls) and as the session token (for user calls). If input is a string, it is assumed to be URL-encoded.
// Otherwise, it will be json encoded.
func (s *Suite) requestResponse(method, path, token string, input any) *httptest.ResponseRecorder {
	var r io.Reader
	var contentType string
	if input != nil {
		if str, ok := input.(string); ok {
			r = strings.NewReader(str)
			contentType = "application/x-www-form-urlencoded"
		} else {
			j, _ := json.Marshal(&input)
			r = bytes.NewReader(j)
			contentType = "application/json"
		}
	}
	req := httptest.NewRequest(method, path, r)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	req.Header.Set("Content-Type", contentType)

	s.session.Values[AccessTokenSessionKey] = token

	res := httptest.NewRecorder()
	s.app.ServeHTTP(res, req)
	return res
}

// request submits a test request and captures the response. The provided token is set as the Bearer token (for API
// calls) and as the session token (for user calls). If input is a string, it is assumed to be URL-encoded. Otherwise,
// it will be json encoded.
func (s *Suite) request(method, path, token string, input any) ([]byte, int) {
	res := s.requestResponse(method, path, token, input)
	body, err := io.ReadAll(res.Body)
	s.NoError(err)
	return body, res.Code
}

// testSessionStore is a limited session store that satisfies the gorilla/sessions Store interface
type testSessionStore struct {
	sessions map[string]*sessions.Session
}

func (s *testSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	if s, ok := s.sessions[name]; ok {
		return s, nil
	}
	return s.New(r, name)
}

func (s *testSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	sess := sessions.NewSession(s, name)
	s.sessions[name] = sess
	return sess, nil
}

func (s *testSessionStore) Save(r *http.Request, w http.ResponseWriter, sess *sessions.Session) error {
	if s.sessions == nil {
		s.sessions = map[string]*sessions.Session{}
	}
	s.sessions[sess.Name()] = sess
	return nil
}

func newTestSessionStore() sessions.Store {
	return &testSessionStore{
		sessions: map[string]*sessions.Session{},
	}
}

func saveToken(db *sql.DB, userID int, token string) {
	_, err := data.CreateAccessToken(context.Background(), db, userID, core.HashAccessToken(token))
	if err != nil {
		panic(err)
	}
}
