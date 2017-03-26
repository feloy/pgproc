[![Build Status](https://travis-ci.org/feloy/pgproc.svg?branch=master)](https://travis-ci.org/feloy/pgproc)

# pgproc - Easily call PostgreSQL procedures from Go (golang)

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

		// Call an SQL procedure returning an integer
        var val int
        base.Call(&val, "public", "age_of_captain")
        fmt.Printf("The captain is %d years old\n", val)
}
```

results in:
```sh
$ go run main.go
The captain is 42 years old
```
