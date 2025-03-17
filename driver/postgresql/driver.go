package postgresql

import (
	_ "github.com/lib/pq"
	"github.com/trueifnotfalse/golang-migrate/interface/driver"
)

type Driver struct {
}

func New() driver.Interface {
	return &Driver{}
}

func (r *Driver) GetName() string {
	return "postgres"
}

func (r *Driver) GetSaveSQL(tableName, migrationName string) string {
	return `insert into ` + tableName + `(name) values('` + migrationName + `')`
}

func (r *Driver) GetDeleteSQL(tableName, migrationName string) string {
	return `delete from ` + tableName + ` where name='` + migrationName + `'`
}

func (r *Driver) GetAppliedSQL(tableName string) string {
	return `select name from ` + tableName + ` order by id desc`
}

func (r *Driver) GetCreateSQL(tableName string) string {
	return `create table if not exists ` + tableName + ` (
id         serial                      primary key,
name       varchar(255)                not null,
apply_time timestamp(0) with time zone not null default CURRENT_TIMESTAMP,
CONSTRAINT ` + tableName + `_name_unique UNIQUE (name)
)`
}
