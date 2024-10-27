package base

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
)

type base struct {
	lg *slog.Logger
	db *sql.DB
}

func NewBase(
	lg *slog.Logger,
	userName string,
	passward string,
	address string,
	nameBase string,
) (*base, error) {
	dsn := (&url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(userName, passward),
		Host:   address,
		Path:   nameBase,
	}).String()

	DB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := DB.Ping(); err != nil {
		return nil, err
	}
	return &base{
		lg: lg.With("base"),
		db: DB,
	}, nil
}

func (b *base) Close() error { return b.db.Close() }

func (b *base) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := b.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *base) DummyMigration(ctx context.Context) error {

	query := `CREATE TABLE IF NOT EXISTS mens (
        id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
        name VARCHAR NOT NULL,
        description VARCHAR,
    	price INTEGER,
        deleted BOOL DEFAULT FALSE NOT NULL
    );`

	if _, err := b.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("create table: %v", err)
	}
	b.lg.Info("Migration is succeed...")

	return nil

}
