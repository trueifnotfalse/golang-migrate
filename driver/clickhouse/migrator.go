package clickhouse

import (
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/trueifnotfalse/golang-migrate/interface/driver"
	"time"
)

type Driver struct {
}

func New() driver.Interface {
	return &Driver{}
}

func (r *Driver) GetName() string {
	return "clickhouse"
}

func (r *Driver) GetSaveSQL(tableName, migrationName string) string {
	return `insert into ` + tableName + `(id, name) values(` + fmt.Sprintf("%d", time.Now().UnixNano()) + `, '` + migrationName + `')`
}

func (r *Driver) GetDeleteSQL(tableName, migrationName string) string {
	return `alter table ` + tableName + ` delete where name='` + migrationName + `'`
}

func (r *Driver) GetAppliedSQL(tableName string) string {
	return `select name from ` + tableName + ` order by id desc`
}

func (r *Driver) GetCreateSQL(tableName string) string {
	return `create table if not exists ` + tableName + ` (
id         UInt64,
name       String,
apply_time DateTime('UTC') default now()
) Engine=TinyLog`
}
