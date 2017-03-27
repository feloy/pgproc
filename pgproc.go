package pgproc

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"github.com/lib/pq"
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

	if proc[0] == '_' {
		return errors.New("function not callable")
	}

	rt, err := p.getReturnType(schema, proc, len(params))
	if err != nil {
		return err
	}
	query := fmt.Sprintf("SELECT * FROM %s.%s(%s)",
		pq.QuoteIdentifier(schema),
		pq.QuoteIdentifier(proc),
		paramsString(len(params)))

	if rt.scalar {
		if !rt.setof {
			row := p.db.QueryRow(query, params...)
			err = row.Scan(result)
		} else {
			rows, _ := p.db.Query(query, params...)
			defer rows.Close()
			c := reflect.ValueOf(result) // the channel we have to send to
			// val is a zero element of the same type of the channel type
			val := reflect.Zero(reflect.TypeOf(result).Elem()).Interface()
			for rows.Next() {
				if err := rows.Scan(&val); err != nil {
					return err
				}
				c.Send(reflect.ValueOf(val))
			}
		}
	} else {
		if !rt.setof {
			row := p.db.QueryRow(query, params...)
			err = ScanCompositeRow(row, rt, result)
		} else {
			rows, _ := p.db.Query(query, params...)
			defer rows.Close()
			for rows.Next() {
				if err := ScanCompositeRows(rows, rt, result); err != nil {
					return err
				}
			}
		}
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


func ScanCompositeRows(rows *sql.Rows, rt *returnType, result interface{}) error {
	c := reflect.ValueOf(result) // the channel we have to send to
	v := reflect.New(reflect.TypeOf(result).Elem()).Elem()
	var vs []interface{}
	
	for _, name := range rt.compositeNames {
		field := v.FieldByName(strings.Title(name)).Addr().Interface()
		vs = append(vs, field)
	}

	err := rows.Scan(vs...)
	c.Send(v)
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
  typtype IN ('b', 'p', 'e')`

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
   WHERE attrelid = pg_type_ret.typrelid AND attnum > 0),
  (SELECT array_agg(typname ORDER BY attnum) FROM pg_attribute 
   INNER JOIN pg_type ON pg_attribute.atttypid = pg_type.oid 
   WHERE attrelid = pg_type_ret.typrelid AND attnum > 0),
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
