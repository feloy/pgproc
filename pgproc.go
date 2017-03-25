package pgproc

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type PgProc struct {
	db *sql.DB
}

// NewPgProc creates a new connection to a PostgreSQL database
func NewPgProc(conninfo string) (*PgProc, error) {
	var pgproc = PgProc{}
	var err error
	pgproc.db, err = sql.Open("postgres", conninfo)
	if err != nil {
		return nil, err
	}
	return &pgproc, nil
}

// Call calls a PostgreSQL procedure and stores the result
func (p *PgProc) Call(result interface{}, schema string, proc string, params ...interface{}) error {

	query := fmt.Sprintf("SELECT * FROM %s.%s(%s)",
		pq.QuoteIdentifier(schema),
		pq.QuoteIdentifier(proc),
		paramsString(len(params)))

	row := p.db.QueryRow(query, params...)
	err := row.Scan(result)
	return err
}

//
// Local static functions
//

// paramsString returns a string $1,$2,...,$len
func paramsString(len int) string {
	if len == 0 {
		return ""
	}
	result := "$1"
	for i := 2; i <= len; i++ {
		result += fmt.Sprintf(",$%d", i)
	}
	return result
}
