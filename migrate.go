package migrate

import (
	"github.com/trueifnotfalse/golang-migrate/database/postgres"
)

func Postgres(dsn string) *postgres.Migrator {
	return postgres.New(dsn)
}
