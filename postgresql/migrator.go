package postgresql

import (
	"github.com/trueifnotfalse/golang-migrate/driver/postgresql"
	migratorI "github.com/trueifnotfalse/golang-migrate/interface/migrator"
	"github.com/trueifnotfalse/golang-migrate/migrator"
)

func New(dsn string) migratorI.Interface {
	return migrator.New(dsn).SetDriver(postgresql.New())
}
