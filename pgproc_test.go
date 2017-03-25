package pgproc

import (
	"testing"
	_ "github.com/lib/pq"
)

func TestNewPgProc(t *testing.T) {
	_, err := NewPgProc("user=pgproctest password=pgproctest host=localhost dbname=pgproctest sslmode=disabled")
	if err != nil {
		t.Errorf("Error connecting to db")
	}
}
