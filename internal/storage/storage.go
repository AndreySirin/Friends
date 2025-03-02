package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AndreySirin/Friends/internal/dirMigrite"
	"log/slog"
	"net/url"

	_ "github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
)

type StorageMethod interface {
	CreatUser(context.Context, *User) error
	GetUser(context.Context, string) (*User, error)
	CreateRefreshToken(context.Context, *RefreshToken) error
	GetRefreshToken(context.Context, string) (*RefreshToken, error)
}

type Storage struct {
	lg *slog.Logger
	db *sql.DB
}

func New(
	lg *slog.Logger,
	username string,
	password string,
	address string,
	database string,
) (*Storage, error) {
	dsn := (&url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(username, password),
		Host:   address,
		Path:   database,
	}).String()

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("init db: %v", err)
	}
	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %v", err)
	}

	return &Storage{
		lg: lg.With("module", "storage"),
		db: sqlDB,
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) MigriteUP() (int, error) {
	path, err := dirMigrite.PathMigrite()
	if err != nil {
		return 0, err
	}
	migrations := &migrate.FileMigrationSource{
		Dir: path,
	}
	n, err := migrate.Exec(s.db, "postgres", migrations, migrate.Up)
	if err != nil {
		s.lg.Error("ошибка", "error", err)
		return n, err
	}
	return n, nil
}

func (s *Storage) MigriteDOWN() (int, error) {
	path, err := dirMigrite.PathMigrite()
	if err != nil {
		return 0, err
	}
	migrations := &migrate.FileMigrationSource{
		Dir: path,
	}
	n, err := migrate.Exec(s.db, "postgres", migrations, migrate.Down)
	if err != nil {
		s.lg.Error("ошибка", "error", err)
		return n, err
	}
	return n, nil
}
