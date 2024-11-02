package storage

import (
	"context"
	"database/sql"
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
		price 	INT
		);`

	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("create table: %v", err)
	}

	s.lg.Info("Migration is succeed...")

	return nil
}

func (s *Storage) PostNewUser(query string, args ...any) error {
	_, err := s.db.ExecContext(context.Background(), query, args...)
	if err != nil {
		slog.Error("err", err)
	}
	return nil
}

//запрос на добавление нового пользователя
//query := `INSERT INTO products (name, hobby, price) VALUES ($1, $2, $3)`
//name := "jorg"
//hobby := "coffemania"
//price := 430
//err = psql.PostNewUser(query, name, hobby, price)
//if err != nil {
//	lg.Error("Failed to create new user")
//}

func (s *Storage) GetQuery(query string, args ...any) (*sql.Rows, error) {
	return s.db.QueryContext(context.Background(), query, args...)
}

type Product struct {
	ID    int
	Name  string
	Hobby string
	Price int
}

func (s *Storage) GetZZZ() ([]Product, error) {

	if s == nil {
		slog.Error("Storage is nil")
		return nil, fmt.Errorf("storage is nil")
	}

	rows, err := s.GetQuery(`SELECT * FROM products`)
	if err != nil {
		slog.Error("Failed to execute query in GetZZZ", err)
		return nil, err
	}
	if rows == nil {
		slog.Error("GetQuery returned nil rows")
		return nil, fmt.Errorf("no rows returned")
	}
	defer rows.Close()

	var product []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Hobby, &p.Price); err != nil {
			slog.Error("Failed to scan row in GetZZZ", err)
			return nil, err
		}
		product = append(product, p)
	}
	if err := rows.Err(); err != nil {
		slog.Error("Error iterating rows in GetZZZ", err)
		return nil, err
	}
	return product, nil
}

//запрос в таблицу на получение данных
//rows, err := psql.GetQuery("SELECT * FROM products WHERE id=$1", "1")
//if err != nil {
//lg.Error("Failed to get products")
//return
//}
//defer rows.Close()
//
//for rows.Next() {
//var id int
//var name string
//var hobby string
//var price int
//if err = rows.Scan(&id, &name, &hobby, &price); err != nil {
//lg.Error("Failed to scan", err)
//}
//fmt.Print(id, name, hobby, price)
//}

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
