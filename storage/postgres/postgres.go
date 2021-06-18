package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(dbstring string) (*Storage, error) {
	db, err := sqlx.Connect("postgres", dbstring)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to postgres '%s'", dbstring)
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}
