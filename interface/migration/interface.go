package migration

import "database/sql"

type Interface interface {
	SetDriver(val *sql.Tx) Interface
	Name() string
	Up() error
	Down() error
}
