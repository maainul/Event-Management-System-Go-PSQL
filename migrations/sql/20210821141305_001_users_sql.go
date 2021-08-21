package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up001UsersSql, down001UsersSql)
}

func up001UsersSql(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func down001UsersSql(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
