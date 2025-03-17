package driver

type Interface interface {
	GetName() string
	GetSaveSQL(tableName, migrationName string) string
	GetDeleteSQL(tableName, migrationName string) string
	GetAppliedSQL(tableName string) string
	GetCreateSQL(tableName string) string
}
