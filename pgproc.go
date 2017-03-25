package pgproc

import (
	"database/sql"
)

type PgProc struct {
	db *sql.DB
}

func NewPgProc(conninfo string) (*PgProc, error) {
	var pgproc = PgProc{}
	var err error
	pgproc.db, err = sql.Open("postgres", conninfo)
	if err != nil {
		return nil, err
	}
	return &pgproc, nil
}
