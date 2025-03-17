# golang-migrate

Database migrations written in Go.

# Example usage:

```go
package main

import (
	"database/sql"
	"fmt"
	"github.com/trueifnotfalse/golang-migrate/interface/migration"
	"github.com/trueifnotfalse/golang-migrate/postgresql"
)

type Migration struct {
	conn *sql.Tx
}

func New() *Migration {
	return &Migration{}
}

func (r *Migration) SetDriver(val *sql.Tx) migration.Interface {
	r.conn = val
	return r
}

func (r *Migration) Name() string {
	return "m202502011005AddPostgisExtension"
}

func (r *Migration) Up() error {
	_, err := r.conn.Exec(`CREATE EXTENSION IF NOT EXISTS postgis;`)

	return err
}

func (r *Migration) Down() error {
	return nil
}

func main() {
	m := postgresql.New("postgres://postgres:postgres@localhost:5432/test")
	m.Register(
		New(),
	)
	err := m.Up()
	if err != nil {
		fmt.Println(err)
	}
}
```
