package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/goaferlx/go-core/log"
)

type DB struct {
	*sql.DB
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{
		now: time.Now(),
		Tx:  tx,
	}, nil
}

type Tx struct {
	now time.Time
	*sql.Tx
}

func (tx *Tx) Now() time.Time {
	return tx.now
}

var ErrNoRows = sql.ErrNoRows

func Open(cfg Config) (*DB, error) {
	db, err := sql.Open(cfg.Dialect, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("sql: %w", err)
	}
	if err := db.PingContext(context.TODO()); err != nil {
		return nil, fmt.Errorf("sql: %w", err)
	}

	log.WithFields(log.Fields{
		"dialect": cfg.Dialect,
		"dbname":  cfg.DBName,
		"host":    cfg.Host,
		"port":    cfg.Port,
	}).Info("Connected to database")

	return &DB{db}, nil
}
