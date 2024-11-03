package storage

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	_ "errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"net/url"
)

const moduleName = "storage"

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
		lg: lg.With("module", moduleName),
		db: sqlDB,
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) DummyMigration(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS products(
		id 		SERIAL PRIMARY KEY,
		name 	VARCHAR (255) NOT NULL,
    	hobby	VARCHAR (255) NOT NULL,
		price 	INT,
    	image_data BYTEA NOT NULL
		);`

	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("create table: %v", err)
	}

	s.lg.Info("Migration is succeed...")

	return nil
}

func (s *Storage) AddProductFriend(ctx context.Context, productFriend ProductFriend) error {
	query := `INSERT INTO products (name, hobby, price, image_data) VALUES ($1, $2, $3, $4)`
	if _, err := s.db.ExecContext(
		ctx,
		query,
		productFriend.Name,
		productFriend.Hobby,
		productFriend.Price,
		productFriend.ImageData,
	); err != nil {
		return fmt.Errorf("add product friend: %v", err)
	}
	return nil
}

func (s *Storage) GetQuery(query string, args ...any) (*sql.Rows, error) {
	return s.db.QueryContext(context.Background(), query, args...)
}

type Product struct {
	ID          int
	Name        string
	Hobby       string
	Price       int
	ImageData   []byte
	ImageBase64 string
}

func (s *Storage) GetZZZ() ([]Product, error) {

	rows, err := s.GetQuery(`SELECT * FROM products`)
	if err != nil {
		slog.Error("Failed to execute query in GetZZZ", err)
		return nil, err
	}

	defer rows.Close()

	var product []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Hobby, &p.Price, &p.ImageData); err != nil {
			slog.Error("Failed to scan row in GetZZZ", err)
			return nil, err
		}
		p.ImageBase64 = base64.StdEncoding.EncodeToString(p.ImageData)
		product = append(product, p)
	}
	if err := rows.Err(); err != nil {
		slog.Error("Error iterating rows in GetZZZ", err)
		return nil, err
	}
	return product, nil
}

func (s *Storage) EncodeImageToBase64(img []byte) string {
	return base64.StdEncoding.EncodeToString(img)
}

func (s *Storage) Getdelete(id int) error {
	result, err := s.db.ExecContext(context.Background(), `DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		slog.Error("err")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("err")
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
