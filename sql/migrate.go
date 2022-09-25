package sql

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (db *DB) MigrateUp() error {
	m, err := db.migrator("file://../sql")
	if err != nil {
		return fmt.Errorf("sql: creating migrator: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("sql: migrating up: %w", err)
	}
	return nil
}

func (db *DB) MigrateDown() error {
	m, err := db.migrator("file://../sql")
	if err != nil {
		return fmt.Errorf("sql: creating migrator: %w", err)
	}

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("sql: migrating down: %w", err)
	}

	return nil
}

func (db *DB) DestructiveReset() error {
	m, err := db.migrator("file://../sql")
	if err != nil {
		return fmt.Errorf("sql: creating migrator: %w", err)
	}

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("sql: migrating down: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("sql: migrating up: %w", err)
	}
	return nil
}

func (db *DB) migrator(filePath string) (*migrate.Migrate, error) {
	instance, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		return nil, fmt.Errorf("sql: creating db instance: %w", err)
	}
	return migrate.NewWithDatabaseInstance(filePath, "mysql", instance)

}
