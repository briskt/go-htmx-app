package data

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/briskt/go-htmx-app/app"
	"github.com/briskt/go-htmx-app/data/sqlc"
	"github.com/briskt/go-htmx-app/log"
)

func init() {
	log.Init()
}

type Suite struct {
	suite.Suite
	*require.Assertions
	ctx context.Context
	db  *sql.DB
}

func (s *Suite) SetupTest() {
	s.Assertions = require.New(s.T())
	DestroyTables(s.db)
}

// TestSuite runs the test suite
func TestSuite(t *testing.T) {
	dsn := fmt.Sprintf("postgresql://%s:%s@test_db:5432/%s?sslmode=disable",
		app.Env.PostgresUser, app.Env.PostgresPassword, app.Env.PostgresDB)
	db, err := sql.Open("pgx", dsn)
	must(err)
	s := &Suite{
		ctx: context.Background(),
		db:  db,
	}
	suite.Run(t, s)
}

func insertUser(db *sql.DB) sqlc.User {
	user, err := q(db).CreateUser(context.Background(), sqlc.CreateUserParams{
		EmployeeID:  "10001",
		FirstName:   "John",
		LastName:    "Doe",
		DisplayName: "Jonny Doe",
		Username:    "john_doe",
		Email:       "john_doe@example.com",
	})
	must(err)
	return user
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
