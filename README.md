[![Build Status](https://travis-ci.org/feloy/pgproc.svg?branch=master)](https://travis-ci.org/feloy/pgproc)

# pgproc - Easily call PostgreSQL procedures from Go (golang)

If you are using PostgreSQL server programming functionalities (user-defined 
procedures and data types essentially) and want to access these functions 
from your Go programs, this library is made for you.

Thanks to both PostgreSQL and Go reflection APIs, it is possible to automatically
map types of parameters and returned values between both worlds.

The main advantage of the `pgproc` library is to be able to call in a single pair of lines
a PostgreSQL stored procedure and get its result in a Go native of `struct` variable.

**This library is at an early development stage, please have a look at the tests files 
to see the covered cases.**

## Install

```sh
$ go get github.com/feloy/pgproc
```

## Usage

```go
package main

import (
        "fmt"
        "github.com/feloy/pgproc"
)

var (
        user     = "myuser"
        password = "mysecret"
        host     = "localhost"
        dbname   = "mydb"
        base     *pgproc.PgProc
)

func main() {
        conninfo := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
                user, password, host, dbname)
        base, err := pgproc.NewPgProc(conninfo)
        if err != nil {
                panic("cannot connect to database")
        }

        // Call an SQL procedure returning an integer value
        var val int
        base.Call(&val, "public", "age_of_captain")
        fmt.Printf("The captain is %d years old\n", val)

        // Call an SQL procedures returning a composite value
	var str struct {
	        Age  int
                Name string
        }
	base.Call(&str, "public", "get_captain_info")
	fmt.Println(str)

}
```

results in:
```sh
$ go run main.go
The captain is 42 years old
{42 Ford Prefect}
```
