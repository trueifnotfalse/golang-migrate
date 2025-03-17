package clickhouse

import (
	"github.com/trueifnotfalse/golang-migrate/driver/clickhouse"
	migratorI "github.com/trueifnotfalse/golang-migrate/interface/migrator"
	"github.com/trueifnotfalse/golang-migrate/migrator"
)

func New(dsn string) migratorI.Interface {
	return migrator.New(dsn).SetDriver(clickhouse.New())
}
