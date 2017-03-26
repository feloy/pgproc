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
	base     *PgProc
)

func TestMain(m *testing.M) {
	base, _ = connect()
	m.Run()
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
	var res int
	err := base.Call(&res, "tests", "unknown_function")
	if err == nil {
		t.Errorf("Error calling unknown function")
	}
}

func TestCallReturnsInteger(t *testing.T) {
	var res int
	err := base.Call(&res, "tests", "test_returns_integer")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_integer")
	}
	if res != 42 {
		t.Errorf("Error expected %d value is %d", 42, res)
	}
}

func TestCallReturnsSetofInteger(t *testing.T) {
	ch := make(chan int64) // Must be int64, Value.Int() returns int64???
	go base.Call(ch, "tests", "test_returns_setof_integer")
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

func TestCallReturnsIntegerAsString(t *testing.T) {
	var res string
	err := base.Call(&res, "tests", "test_returns_integer_as_string")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_integer_as_string")
	}
	if res != "42" {
		t.Errorf("Error expected '%s' value is '%s'", "42", res)
	}
}

func TestCallReturnsString(t *testing.T) {
	var res string
	err := base.Call(&res, "tests", "test_returns_string")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_string")
	}
	if res != "hello" {
		t.Errorf("Error expected '%s' value is '%s'", "hello", res)
	}
}

func TestCallReturnsSetofString(t *testing.T) {
	ch := make(chan string)
	go base.Call(ch, "tests", "test_returns_setof_string")
	a := <-ch
	b := <-ch
	c := <-ch
	if a != "hello"|| b != "world" || c != "!" {
		t.Errorf("Error expected values")
	}
}

func TestCallReturnsNumeric(t *testing.T) {
	var res float32
	err := base.Call(&res, "tests", "test_returns_numeric")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_numeric")
	}
	if res != 3.14159 {
		t.Errorf("Error expected %f value is %f", 3.14159, res)
	}
}

func TestCallReturnsReal(t *testing.T) {
	var res float32
	err := base.Call(&res, "tests", "test_returns_real")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_real")
	}
	if res != 3.14 {
		t.Errorf("Error expected %f value is %f", 3.14, res)
	}
}

func TestCallReturnsBoolTrue(t *testing.T) {
	var res bool
	err := base.Call(&res, "tests", "test_returns_bool_true")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_bool_true")
	}
	if res != true {
		t.Errorf("Error expected %t value is %t", true, res)
	}
}

func TestCallReturnsBoolFalse(t *testing.T) {
	var res bool
	err := base.Call(&res, "tests", "test_returns_bool_false")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_bool_false")
	}
	if res != false {
		t.Errorf("Error expected %t value is %t", false, res)
	}
}

func TestCallReturnsDate(t *testing.T) {
	var res time.Time
	err := base.Call(&res, "tests", "test_returns_date")
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

func TestCallReturnsInfinityDate(t *testing.T) {
	var res time.Time
	err := base.Call(&res, "tests", "test_returns_infinity_date")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_infinity_date")
	}
	if res != DateInfinity {
		t.Errorf("Error date infinity")
	}
}

func TestCallReturnsMinusInfinityDate(t *testing.T) {
	var res time.Time
	err := base.Call(&res, "tests", "test_returns_minus_infinity_date")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_minus_infinity_date")
	}
	if res != DateMinusInfinity {
		t.Errorf("Error date infinity")
	}
}

func TestCallReturns64bitsDate(t *testing.T) {
	var res time.Time
	err := base.Call(&res, "tests", "test_returns_64bits_date")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_64bits_date")
	}
	y, m, d := res.Date()
	if y != 2040 || m != 1 || d != 1 {
		t.Errorf("Error 64bits date")
	}
}

func TestCallReturnsTimestamp(t *testing.T) {
	var res time.Time
	err := base.Call(&res, "tests", "test_returns_timestamp")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_timestamp")
	}
	yExp, mExp, dExp := res.Date()
	yNow, mNow, dNow := time.Now().Date()
	if yExp != yNow || mExp != mNow || dExp != dNow {
		t.Errorf("Error expected value")
	}
}


func TestCallReturnsTime(t *testing.T) {
	var res time.Time
	err := base.Call(&res, "tests", "test_returns_time")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_time")
	}
	y, m, d := res.Date()
	if y != 0 || m != 1 || d != 1 {
		t.Errorf("Error expected date 0/1/1")
	}
}

func TestCallReturnsComposite(t *testing.T) {
	var res struct {
		A int
		B string
	}
	err := base.Call(&res, "tests", "test_returns_composite")
	if err != nil {
		fmt.Println(err)
		t.Errorf("Error calling tests.test_returns_composite")
	}
	if res.A != 1 || res.B != "hello" {
		t.Errorf("Error expected value")
	}
}

// The order of fields of the struct doo not need to respect 
// order of fields in PostgreSQL composite type
func TestCallReturnsCompositeRandomOrder(t *testing.T) {
	var res struct {
		B string
		A int
	}
	err := base.Call(&res, "tests", "test_returns_composite")
	if err != nil {
		fmt.Println(err)
		t.Errorf("Error calling tests.test_returns_composite")
	}
	if res.A != 1 || res.B != "hello" {
		t.Errorf("Error expected value")
	}
}
