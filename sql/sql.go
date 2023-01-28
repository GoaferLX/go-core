package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Default values for configuring the DB connection pool.  Values taken from
// https://www.alexedwards.net/blog/configuring-sqldb - accessed 2023/01/28.
const (
	DefaultMaxOpenConns int           = 25
	DefaultMaxIdleConns int           = 5
	DefaultMaxLifetime  time.Duration = 5 * time.Minute
	DefaultMaxIdleTime  time.Duration = 5 * time.Minute
)

type DB struct {
	*sql.DB
	path string // path to migration files.
}

// BeginTx wraps the sql.BeginTx and sets a tx time.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{
		now: time.Now().UTC(),
		Tx:  tx,
	}, nil
}

type Tx struct {
	now time.Time
	*sql.Tx
}

// Now returns the time that the Tx started at.
func (tx *Tx) Now() time.Time {
	return tx.now.UTC()
}

var ErrNoRows = sql.ErrNoRows

// Open is a convenience function that wraps sql.Open to establish a connection to the DB
// and verifies the connection, in one step, as well as setting sensible default values
// for the connection pool.
func Open(cfg Config) (*DB, error) {
	db, err := sql.Open(cfg.Dialect, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("sql: %w", err)
	}
	if err := db.PingContext(context.TODO()); err != nil {
		return nil, fmt.Errorf("sql: %w", err)
	}

	// TODO: Change these to values from config.
	db.SetMaxOpenConns(DefaultMaxOpenConns)
	db.SetMaxIdleConns(DefaultMaxIdleConns)
	db.SetConnMaxLifetime(DefaultMaxLifetime)
	db.SetConnMaxIdleTime(DefaultMaxIdleTime)

	return &DB{DB: db}, nil
}
