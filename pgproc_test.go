package pgproc

import (
	"fmt"
	"math"
	"testing"
	"time"
	"github.com/lib/pq"
)

const (
	TOLERANCE = 1e-6
)

var (
	user     = "pgproctest"
	password = "pgproctest"
	host     = "localhost"
	dbname   = "pgproctest"
	base     *PgProc
)

type Content struct {
	Cnt_id int
	Cnt_name string
}

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

func dateEqual(date1 time.Time, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func clockEqual(date1 time.Time, date2 time.Time) bool {
	h1, m1, s1 := date1.Clock()
	h2, m2, s2 := date2.Clock()
	return h1 == h2 && m1 == m2 && s1 == s2
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
	a := <-ch
	b := <-ch
	c := <-ch
	if a != 42 || b != 43 || c != 44 {
		t.Errorf("Error expected values")
	}
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

/*
func TestCallReturnsSetofNumeric(t *testing.T) {
	ch := make(chan float32)
	go base.Call(ch, "tests", "test_returns_setof_numeric")
	a := <- ch
	b := <- ch
	if a != 3.14159 || b != 4.49 {
		t.Errorf("Error expected value")
	}
}
*/

func TestCallReturnsReal(t *testing.T) {
	var res float32
	err := base.Call(&res, "tests", "test_returns_real")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_real")
	}
	if math.Abs(float64(res - 3.14)) > TOLERANCE {
		t.Errorf("Error expected %f value is %f", 3.14, res)
	}
}

func TestCallReturnsSetofReal(t *testing.T) {
	ch := make(chan float64) // Must be float64, Value.Float() returns float64???
	go base.Call(ch, "tests", "test_returns_setof_real")
	a := <- ch
	b := <- ch
	if math.Abs(a - 3.14) > TOLERANCE || math.Abs(b - 4.49) > TOLERANCE {
		t.Errorf("Error expected values %f %f are %f %f", 3.14, 4.49, a, b)
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

func TestCallReturnsSetofBool(t *testing.T) {
	ch := make(chan bool)
	go base.Call(ch, "tests", "test_returns_setof_bool")
	a := <- ch
	b := <- ch
	c := <- ch
	d := <- ch
	if a != false || b != true || c != true || d != false {
		t.Errorf("Error expected values")

	}
}

func TestCallReturnsDate(t *testing.T) {
	var res time.Time
	err := base.Call(&res, "tests", "test_returns_date")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_date")
	}
	if !dateEqual(res, time.Now()) {
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

func TestCallReturnsSetofDate(t *testing.T) {
	ch := make(chan time.Time)
	go base.Call(ch, "tests", "test_returns_setof_date")
	a := <- ch
	b := <- ch
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	if y1 != 2015 || m1 != 1 || d1 != 1 ||
		y2 != 2016 || m2 != 2 || d2 != 2 {
		t.Errorf("Error expected values")		
	}
}

func TestCallReturnsTimestamp(t *testing.T) {
	var res time.Time
	err := base.Call(&res, "tests", "test_returns_timestamp")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_timestamp")
	}
	if !dateEqual(res, time.Now()) {
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

func TestCallReturnsSetofComposite(t *testing.T) {
	type T struct {
		A int
		B string
	}
	ch := make (chan T)
	go base.Call(ch, "tests", "test_returns_setof_composite")
	res1 := <-ch
	res2 := <-ch
	if res1.A != 1 || res1.B != "hello" {
		t.Errorf("Error expected value")
	}
	if res2.A != 2 || res2.B != "bye" {
		t.Errorf("Error expected value")
	}
}

// TODO pass res as a user-defined type with constants 
// type Enumtype int
// const (
//     Val1 MyEnum = iota
//     Val2 
//     Val3
// )
func TestCallReturnsEnum(t *testing.T) {
	var res string
	err := base.Call(&res, "tests", "test_returns_enum")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_enum")
		t.Error(err)
	}
	if res != "val1" {
		t.Errorf("Error expected %s value is %s", "val1", res)
	}
}

// TODO get an array and not an []uint8
func TestCallReturnsEnumArray(t *testing.T) {
	var res pq.StringArray
	err := base.Call(&res, "tests", "test_returns_enum_array")
	if err != nil {
		t.Errorf("Error calling tests.test_returns_enum_array")
		t.Error(err)
	}
	if len(res) != 2 || res[0] != "val1" || res[1] != "val2" {
		t.Errorf("Error expected values")
	}
}

func TestHiddenFunction(t *testing.T) {
	var res bool
	err := base.Call(&res, "tests", "_hidden_function")
	if err == nil {
		t.Errorf("Error calling hidden function")
	}
}

func TestFunctionRaisingException(t *testing.T) {
	var res bool
	err := base.Call(&res, "tests", "function_raising_exception")
	if err == nil {
		t.Errorf("Error calling function_raising_exception")
	}
}

func TestReturnsIncrementedInteger(t *testing.T) {
	var res int
	input := int(5)
	err := base.Call(&res, "tests", "test_returns_incremented_integer", input)
	if err != nil {
		t.Errorf("Error calling test_returns_incremented_integer")
	}
	if res != input + 1 {
		t.Errorf("Error expected %d value is %d\n", input + 1, res)
	}
}

func TestReturnsIncrementedNumeric(t *testing.T) {
	var res float32
	input := float32(5.4)
	err := base.Call(&res, "tests", "test_returns_incremented_numeric", input)
	if err != nil {
		t.Errorf("Error calling test_returns_incremented_numeric")
	}
	if res != input + 1.5 {
		t.Errorf("Error expected %f value is %f\n", input + 1.5, res)
	}
}

func TestReturnsIncrementedReal(t *testing.T) {
	var res float32
	input := float32(5.4)
	err := base.Call(&res, "tests", "test_returns_incremented_real", input)
	if err != nil {
		t.Errorf("Error calling test_returns_incremented_real")
	}
	if res != input + 1.42 {
		t.Errorf("Error expected %f value is %f\n", input + 1.42, res)
	}
}

func TestReturnsCatString(t *testing.T) {
	var res string
	input := string("hello")
	err := base.Call(&res, "tests", "test_returns_cat_string", input)
	if err != nil {
		t.Errorf("Error calling test_returns_cat_string")
	}
	if res != input + "." {
		t.Errorf("Error expected '%s' value is '%s'\n", input + ".", res)
	}
}

func TestReturnsSameBool(t *testing.T) {
	var res bool
	input := false
	err := base.Call(&res, "tests", "test_returns_same_bool", input)
	if err != nil {
		t.Errorf("Error calling test_returns_same_bool")
	}
	if res != input {
		t.Errorf("Error expected '%t' value is '%t'\n", input, res)
	}
	
	input = true
	err = base.Call(&res, "tests", "test_returns_same_bool", input)
	if err != nil {
		t.Errorf("Error calling test_returns_same_bool")
	}
	if res != input {
		t.Errorf("Error expected '%t' value is '%t'\n", input, res)
	}
}

func TestReturnsSameDate(t *testing.T) {
	var res time.Time
	input := time.Now()
	err := base.Call(&res, "tests", "test_returns_same_date", input)
	if err != nil {
		t.Errorf("Error calling test_returns_same_date")
	}
	if !dateEqual(input, res) {
		t.Errorf("Error expected '%s' value is '%s'\n", input, res)
	}
}

func TestReturnsSameTimestamp(t *testing.T) {
	var res time.Time
	input := time.Now()
	err := base.Call(&res, "tests", "test_returns_same_timestamp", input)
	if err != nil {
		t.Errorf("Error calling test_returns_same_timestamp")
	}
	if !dateEqual(input, res) || !clockEqual(input, res) {
		t.Errorf("Error expected '%s' value is '%s'\n", input, res)
	}
}

func TestReturnsSameTime(t *testing.T) {
	var res time.Time
	input := time.Now()
	err := base.Call(&res, "tests", "test_returns_same_time", input)
	if err != nil {
		t.Errorf("Error calling test_returns_same_time")
	}
	if !clockEqual(input, res) {
		t.Errorf("Error expected '%s' value is '%s'\n", input, res)
	}
}

// TODO compute pqVal in pgproc.Call
func TestIntegerArrayArg(t *testing.T) {
	var ch = make(chan int64)
	var input pq.Int64Array = pq.Int64Array{1, 3, 7}
	pqVal, _ := input.Value()
	go base.Call(ch, "tests", "test_integer_array_arg", pqVal)
	a := <-ch
	b := <-ch
	c := <-ch
	if a != 1 || b != 3 || c != 7 {
		t.Errorf("Error expected values")
	}
}

// TODO compute pqVal in pgproc.Call
func TestVarcharArrayArg(t *testing.T) {
	var ch = make(chan string)
	var input pq.StringArray = pq.StringArray{"foo", "bar"}
	pqVal, _ := input.Value()
	go base.Call(ch, "tests", "test_varchar_array_arg", pqVal)
	a := <-ch
	b := <-ch
	if a != "foo" || b != "bar" {
		t.Errorf("Error expected values")
	}
}

func TestEnumArg(t *testing.T) {
	var res string
	input := "val1"
	err := base.Call(&res, "tests", "test_enum_arg", input)
	if err != nil {
		t.Errorf("Error calling test_enum_arg")
	}
	if res != input {
		t.Errorf("Error expected %s value is %s", input, res)
	}
}

// TODO why result as []uint8?
func TestEnumArrayArg(t *testing.T) {
	var ch = make(chan []uint8)
	var input pq.StringArray = pq.StringArray{"val3", "val1"}
	pqVal, _ := input.Value()
	go base.Call(ch, "tests", "test_enum_array_arg", pqVal)
	a := <-ch
	b := <-ch
	if string(a) != "val3" || string(b) != "val1" {
		t.Errorf("Error expected values")
	}
}

func TestAccentedString(t *testing.T) {
	var res string
	err := base.Call(&res, "tests", "test_returns_accented_string")
	if err != nil {
		t.Errorf("Error calling test_returns_accented_string")
	}
	expected := "héllo"
	if res != expected {
		t.Errorf("Error expected %s value is %s", expected, res)
	}
}

func TestEmptyArray(t *testing.T) {
	var res pq.Int64Array
	go base.Call(&res, "tests", "test_returns_empty_array")
	if len(res) != 0 {
		t.Errorf("Error expected empty array")
	}
}

func TestContentAddAndGet(t *testing.T) {
	var res int
	aName := "a name"
	base.Call(&res, "tests", "content_add", aName)
	if res <= 0 {
		t.Errorf("Expected positive value is %d", res)
	}

	var item Content
	base.Call(&item, "tests", "content_get", res)
	if item.Cnt_id != res || item.Cnt_name != aName {
		t.Errorf("Error Expected value")
	}
	
}
