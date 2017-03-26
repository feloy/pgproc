package pgproc

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"time"
	"reflect"
	"strings"
)

type PgProc struct {
	db *sql.DB
}

type returnType struct {
	scalar bool
	setof bool
	scalarType string
	compositeNames pq.StringArray
	compositeTypes pq.StringArray
}

var (
	DateMinusInfinity = time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
	DateInfinity      = time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC)
)

// NewPgProc creates a new connection to a PostgreSQL database
func NewPgProc(conninfo string) (*PgProc, error) {
	var pgproc = PgProc{}
	var err error
	pgproc.db, err = sql.Open("postgres", conninfo)
	if err != nil {
		return nil, err
	}
	pq.EnableInfinityTs(DateMinusInfinity, DateInfinity)
	return &pgproc, nil
}

// Call calls a PostgreSQL procedure and stores the result
func (p *PgProc) Call(result interface{}, schema string, proc string, params ...interface{}) error {

	rt, err := p.getReturnType(schema, proc, len(params))
	if err != nil {
		return err
	}
	query := fmt.Sprintf("SELECT * FROM %s.%s(%s)",
		pq.QuoteIdentifier(schema),
		pq.QuoteIdentifier(proc),
		paramsString(len(params)))

	row := p.db.QueryRow(query, params...)
	if rt.scalar {
		err = row.Scan(result)
	} else {
		err = ScanCompositeRow(row, rt, result)
	}
	return err
}

func ScanCompositeRow(row *sql.Row, rt *returnType, result interface{}) error {
	v := reflect.ValueOf(result).Elem()
	var vs []interface{}

	for _, name := range rt.compositeNames {
		field := v.FieldByName(strings.Title(name)).Addr().Interface()
		vs = append(vs, field)
	}

	err := row.Scan(vs...)
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

// getReturnType gives the type returned by a postgreSQL procedure
func (p *PgProc) getReturnType(schema string, proc string, nargs int) (*returnType, error) {
	rt, err := p.getScalarReturnType(schema, proc, nargs)
	if err == sql.ErrNoRows {
		return p.getCompositeReturnType(schema, proc, nargs)
	} else {
		return rt, nil
	}
}

// getScalarReturnType gives the scalar type returned by a postgreSQL procedure 
// or returns a ErrNoRows error if the return type is not scalar
func (p *PgProc) getScalarReturnType(schema string, proc string, nargs int) (*returnType, error) {
	query := `
SELECT
  pg_type_ret.typname, 
  proretset
FROM pg_proc
INNER JOIN pg_type pg_type_ret ON pg_type_ret.oid = pg_proc.prorettype
INNER JOIN pg_namespace pg_namespace_ret ON pg_namespace_ret.oid = pg_type_ret.typnamespace
INNER JOIN pg_namespace pg_namespace_proc ON pg_namespace_proc.oid = pg_proc.pronamespace
WHERE 
  pg_namespace_proc.nspname = $1 AND 
  proname = $2 AND 
  pronargs = $3 AND 
  typtype IN ('b', 'p') AND
  pg_namespace_ret.nspname = 'pg_catalog'`

	row := p.db.QueryRow(query, schema, proc, nargs)
	var (
		name string
		setof bool
	)
	err := row.Scan(&name, &setof)
	if err == sql.ErrNoRows {
		return nil, err
	} else {
		return &returnType{scalar: true, setof: setof, scalarType: name}, nil
	}
}

// getCompositeReturnType gives the compiste type returned by a postgreSQL procedure 
func (p *PgProc) getCompositeReturnType(schema string, proc string, nargs int) (*returnType, error) {
	query := `
SELECT 
  (SELECT array_agg(attname ORDER BY attnum) FROM pg_attribute 
   WHERE attrelid = pg_type_ret.typrelid),
  (SELECT array_agg(typname ORDER BY attnum) FROM pg_attribute 
   INNER JOIN pg_type ON pg_attribute.atttypid = pg_type.oid 
   WHERE attrelid = pg_type_ret.typrelid),
  proretset
FROM pg_proc
INNER JOIN pg_type pg_type_ret ON pg_type_ret.oid = pg_proc.prorettype
INNER JOIN pg_namespace pg_namespace_proc ON pg_namespace_proc.oid = pg_proc.pronamespace
WHERE 
  pg_namespace_proc.nspname = $1 AND 
  proname = $2 AND 
  pronargs = $3 AND
  pg_type_ret.typtype IN ('c')`

	row := p.db.QueryRow(query, schema, proc, nargs)
	var (
		names pq.StringArray
		types pq.StringArray
		setof bool
	)
	err := row.Scan(&names, &types, &setof)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else {
		return &returnType{scalar: false, setof: setof, compositeNames: names, compositeTypes: types}, nil
	}
	
}
