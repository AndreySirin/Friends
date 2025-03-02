package storage

import (
	"context"
	"errors"
	"fmt"
)

type ProductFriend struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Hobby string `json:"hobby"`
	Price int    `json:"price"`
}

type MethFriend interface {
	GetQueryDB() ([]ProductFriend, error)
	AddProductFriend(context.Context, *ProductFriend) error
	DeleteProductFriend(int) error
	UpdateProductFriend(context.Context, *ProductFriend) error
}

func (s *Storage) AddProductFriend(ctx context.Context, productFriend *ProductFriend) error {
	query := `INSERT INTO products (name, hobby, price)  VALUES ($1, $2, $3) RETURNING id`
	err := s.db.QueryRowContext(
		ctx,
		query,
		productFriend.Name,
		productFriend.Hobby,
		productFriend.Price,
	).Scan(&productFriend.ID)
	if err != nil {
		return fmt.Errorf("add product friend to db: %v", err)
	}
	return nil
}

func (s *Storage) UpdateProductFriend(ctx context.Context, productFriend *ProductFriend) error {
	query := `UPDATE products SET name=$1, hobby=$2, price=$3 WHERE id=$4`
	_, err := s.db.ExecContext(
		ctx,
		query,
		productFriend.Name,
		productFriend.Hobby,
		productFriend.Price,
		productFriend.ID)
	if err != nil {
		return fmt.Errorf("update product friend to db: %v", err)
	}
	return nil
}

func (s *Storage) DeleteProductFriend(id int) error {
	result, err := s.db.ExecContext(context.Background(), `DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		s.lg.Error("err")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.lg.Error("err")
	}
	if rowsAffected == 0 {
		return errors.New("no products found")
	}
	return nil
}

func (s *Storage) GetQueryDB() ([]ProductFriend, error) {
	rows, err := s.db.Query(`SELECT * FROM products`)
	if err != nil {
		s.lg.Error("Failed to execute query in GetZZZ", "error", err)
		return nil, err
	}

	defer rows.Close()

	var product []ProductFriend
	for rows.Next() {
		var p ProductFriend
		if err = rows.Scan(&p.ID, &p.Name, &p.Hobby, &p.Price); err != nil {
			s.lg.Error("Failed to scan row in GetZZZ", "error", err)
			return nil, err
		}
		product = append(product, p)
	}
	if err = rows.Err(); err != nil {
		s.lg.Error("Error iterating rows in GetZZZ", "error", err)
		return nil, err
	}
	return product, nil
}
