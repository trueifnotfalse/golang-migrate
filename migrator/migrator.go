package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/trueifnotfalse/golang-migrate/interface/driver"
	"github.com/trueifnotfalse/golang-migrate/interface/migration"
	"github.com/trueifnotfalse/golang-migrate/interface/migrator"
	"slices"
)

const defaultTableName = "migrations"

type Migrator struct {
	conn          *sql.DB
	migrationList []migration.Interface
	dsn           string
	tableName     string
	driver        driver.Interface
}

func New(dsn string) *Migrator {

	return &Migrator{
		dsn:       dsn,
		tableName: defaultTableName,
	}
}

func (r *Migrator) SetDriver(v driver.Interface) migrator.Interface {
	r.driver = v

	return r
}

func (r *Migrator) SetTableName(val string) migrator.Interface {
	r.tableName = val

	return r
}

func (r *Migrator) Register(values ...migration.Interface) migrator.Interface {
	for i := 0; i < len(values); i++ {
		r.migrationList = append(r.migrationList, values[i])
	}

	return r
}

func (r *Migrator) Up() error {
	err := r.openConnection()
	if err != nil {
		return fmt.Errorf("cannot connect to DB: %s", err.Error())
	}
	defer r.closeConnection()
	err = r.createMigrationTable()
	if err != nil {
		return err
	}
	appliedMigrationList, err := r.getAppliedMigrations()
	if err != nil {
		return err
	}
	migrationList := r.getNotAppliedMigrations(appliedMigrationList)
	err = r.upMigrations(migrationList)
	if err != nil {
		return err
	}

	return nil
}

func (r *Migrator) Down() error {
	err := r.openConnection()
	if err != nil {
		return fmt.Errorf("cannot connect to DB: %s", err.Error())
	}
	defer r.closeConnection()
	err = r.createMigrationTable()
	if err != nil {
		return err
	}
	appliedMigrationList, err := r.getAppliedMigrations()
	if err != nil {
		return err
	}
	indexedMigrationList := r.indexRegisteredMigrationsByName()
	err = r.downMigrations(appliedMigrationList, indexedMigrationList)
	if err != nil {
		return err
	}

	return nil
}

func (r *Migrator) openConnection() error {
	var err error
	r.conn, err = sql.Open(r.driver.GetName(), r.dsn)
	if err != nil {
		return err
	}

	return nil
}

func (r *Migrator) closeConnection() error {
	return r.conn.Close()
}

func (r *Migrator) downMigrations(appliedMigrationList []string, registeredMigrations map[string]migration.Interface) error {
	var (
		err error
		tx  *sql.Tx
		ok  bool
	)
	for i := 0; i < len(appliedMigrationList); i++ {
		_, ok = registeredMigrations[appliedMigrationList[i]]
		if !ok {
			continue
		}
		tx, err = r.conn.BeginTx(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("error opening transaction: %s", err.Error())
		}
		registeredMigrations[appliedMigrationList[i]].SetDriver(tx)
		err = registeredMigrations[appliedMigrationList[i]].Down()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error on rollback migration: %s, %s", registeredMigrations[appliedMigrationList[i]].Name(), err.Error())
		}
		err = r.removeRollbackMigration(registeredMigrations[appliedMigrationList[i]].Name())
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("error on applying commit: %s, %s", registeredMigrations[appliedMigrationList[i]].Name(), err.Error())
		}
	}

	return nil
}

func (r *Migrator) upMigrations(migrationList []migration.Interface) error {
	var (
		err error
		tx  *sql.Tx
	)
	for i := 0; i < len(migrationList); i++ {
		tx, err = r.conn.BeginTx(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("error opening transaction: %s", err.Error())
		}
		migrationList[i].SetDriver(tx)
		err = migrationList[i].Up()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error on applying migration: %s, %s", migrationList[i].Name(), err.Error())
		}
		err = r.saveAppliedMigration(migrationList[i].Name())
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("error on applying commit: %s, %s", migrationList[i].Name(), err.Error())
		}
	}

	return nil
}

func (r *Migrator) indexRegisteredMigrationsByName() map[string]migration.Interface {
	result := make(map[string]migration.Interface, len(r.migrationList))
	for i := 0; i < len(r.migrationList); i++ {
		result[r.migrationList[i].Name()] = r.migrationList[i]
	}

	return result
}

func (r *Migrator) getNotAppliedMigrations(migrationList []string) []migration.Interface {
	if len(migrationList) == 0 {
		return r.migrationList
	}
	var result []migration.Interface
	for i := 0; i < len(r.migrationList); i++ {
		if slices.Contains(migrationList, r.migrationList[i].Name()) {
			continue
		}
		result = append(result, r.migrationList[i])
	}

	return result
}

func (r *Migrator) saveAppliedMigration(name string) error {
	query := r.driver.GetSaveSQL(r.tableName, name)
	_, err := r.conn.Exec(query)
	if err != nil {
		return fmt.Errorf("error on saving migration: %s", err.Error())
	}

	return nil
}

func (r *Migrator) removeRollbackMigration(name string) error {
	query := r.driver.GetDeleteSQL(r.tableName, name)
	_, err := r.conn.Exec(query)
	if err != nil {
		return fmt.Errorf("error on removing migration: %s", err.Error())
	}

	return nil
}

func (r *Migrator) getAppliedMigrations() ([]string, error) {
	query := r.driver.GetAppliedSQL(r.tableName)
	rows, err := r.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error retrieving applied migrations: %s", err.Error())
	}
	defer rows.Close()
	var migrationList []string
	var name string
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("error retrieving applied migrations: %s", err.Error())
		}
		migrationList = append(migrationList, name)
	}

	return migrationList, nil
}

func (r *Migrator) createMigrationTable() error {
	query := r.driver.GetCreateSQL(r.tableName)
	_, err := r.conn.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating migration table: %s", err.Error())
	}

	return nil
}
