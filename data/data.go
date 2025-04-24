package data

import (
	"database/sql"
	"errors"

	"github.com/briskt/go-htmx-app/data/sqlc"
)

var ErrorRowNotUpdated = errors.New("row not updated")

func q(db sqlc.DBTX) *sqlc.Queries {
	return sqlc.New(db)
}

func newNullInt32(i int32) sql.NullInt32 {
	return sql.NullInt32{Int32: i, Valid: true}
}

func DestroyTables(db *sql.DB) {
	resultMust(db.Exec("DELETE FROM email_logs"))
	resultMust(db.Exec("DELETE FROM tokens"))
	resultMust(db.Exec("DELETE FROM users"))
}

func resultMust(_ sql.Result, err error) {
	if err != nil {
		panic(err)
	}
}
