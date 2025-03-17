package migrator

import "github.com/trueifnotfalse/golang-migrate/interface/migration"

type Interface interface {
	SetTableName(val string) Interface
	Register(values ...migration.Interface) Interface
	Up() error
	Down() error
}
