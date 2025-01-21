package migrate

import (
	"github.com/trueifnotfalse/golang-migrate/database/clickhouse"
	"github.com/trueifnotfalse/golang-migrate/database/postgres"
)

func Postgres(dsn string) *postgres.Migrator {
	return postgres.New(dsn)
}

func Clickhouse(dsn string) *clickhouse.Migrator {
	return clickhouse.New(dsn)
}
