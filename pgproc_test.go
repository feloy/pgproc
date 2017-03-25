package pgproc

import (	
	"fmt"
	_ "github.com/lib/pq"
	"testing"
	"time"
)

var (
	user     = "pgproctest"
	password = "pgproctest"
	host     = "localhost"
	dbname   = "pgproctest"
)

func TestNewPgProc(t *testing.T) {
	conninfo := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
		user, password, host, dbname)
	_, err := NewPgProc(conninfo)
	if err != nil {
		t.Errorf("Error connecting to db")
	}
}

func connect() (*PgProc, error) {
	conninfo := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
		user, password, host, dbname)
	base, err := NewPgProc(conninfo)
	return base, err
}

func TestParamsString(t *testing.T) {
	args0 := paramsString(0)
	wanted0 := ""
	if args0 != wanted0 {
		t.Errorf("paramsString with 0 params should be '%s' but is '%s'",
			wanted0, args0)
	}

	args1 := paramsString(1)
	wanted1 := "$1"
	if args1 != wanted1 {
		t.Errorf("paramsString with 1 params should be '%s' but is '%s'",
			wanted1, args1)
	}

	args4 := paramsString(4)
	wanted4 := "$1,$2,$3,$4"
	if args4 != wanted4 {
		t.Errorf("paramsString with 4 params should be '%s' but is '%s'",
			wanted4, args4)
	}
}

func TestCallUnknown(t *testing.T) {
	base, err := connect()
	var res int
	err = base.Call(&res, "tests", "unknown_function")
	if err == nil {
		t.Errorf("Error calling unknown function")
	}
}

func TestCallReturnsInteger(t *testing.T) {
	base, err := connect()
	var res int
	err = base.Call(&res, "tests", "test_returns_integer")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_integer")
	}
	if res != 42 {
		t.Errorf("Error expected %d value is %d", 42, res)
	}
}

func TestCallReturnsIntegerAsString(t *testing.T) {
	base, err := connect()
	var res string
	err = base.Call(&res, "tests", "test_returns_integer_as_string")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_integer_as_string")
	}
	if res != "42" {
		t.Errorf("Error expected '%s' value is '%s'", "42", res)
	}
}

func TestCallReturnsString(t *testing.T) {
	base, err := connect()
	var res string
	err = base.Call(&res, "tests", "test_returns_string")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_string")
	}
	if res != "hello" {
		t.Errorf("Error expected '%s' value is '%s'", "hello", res)
	}
}

func TestCallReturnsNumeric(t *testing.T) {
	base, err := connect()
	var res float32
	err = base.Call(&res, "tests", "test_returns_numeric")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_numeric")
	}
	if res != 3.14159 {
		t.Errorf("Error expected %f value is %f", 3.14159, res)
	}
}

func TestCallReturnsReal(t *testing.T) {
	base, err := connect()
	var res float32
	err = base.Call(&res, "tests", "test_returns_real")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_real")
	}
	if res != 3.14 {
		t.Errorf("Error expected %f value is %f", 3.14, res)
	}
}

func TestCallReturnsBoolTrue(t *testing.T) {
	base, err := connect()
	var res bool
	err = base.Call(&res, "tests", "test_returns_bool_true")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_bool_true")
	}
	if res != true {
		t.Errorf("Error expected %t value is %t", true, res)
	}
}

func TestCallReturnsBoolFalse(t *testing.T) {
	base, err := connect()
	var res bool
	err = base.Call(&res, "tests", "test_returns_bool_false")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_bool_false")
	}
	if res != false {
		t.Errorf("Error expected %t value is %t", false, res)
	}
}

func TestCallReturnsDate(t *testing.T) {
	base, err := connect()
	var res time.Time
	err = base.Call(&res, "tests", "test_returns_date")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_date")
	}
	yExp, mExp, dExp := res.Date()
	yNow, mNow, dNow := time.Now().Date()
	if yExp != yNow || mExp != mNow || dExp != dNow {
		t.Errorf("Error expected value")
	}
	h, m, s := res.Clock()
	if h != 0 || m != 0 || s != 0 {
		t.Errorf("Error expected clock 0")
	}
}

