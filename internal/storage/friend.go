package storage

import (
	"context"
	"errors"
	"fmt"
)

type ProductFriend struct {
	ID    int     `json:"id" validate:"required"`
	Name  string  `json:"name" validate:"required"`
	Hobby string  `json:"hobby" validate:"required,gte=6"`
	Price float64 `json:"price" validate:"required,gt=0"`
}

func (s *Storage) AddProductFriend(ctx context.Context, productFriend *ProductFriend) error {
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO products (name, hobby, price)  VALUES ($1, $2, $3) RETURNING id`,
		productFriend.Name,
		productFriend.Hobby,
		productFriend.Price,
	).Scan(&productFriend.ID)
	if err != nil {
		return fmt.Errorf("failed to insert product friend into db: %w", err)
	}
	return nil
}

func (s *Storage) UpdateProductFriend(ctx context.Context, productFriend *ProductFriend) error {
	var err error
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for UpdateProductFriend: %w", err)
	}

	defer func() {
		if err != nil {
			if RollbackErr := tx.Rollback(); RollbackErr != nil {
				err = fmt.Errorf("error rolling back transaction: %w", err)
			}
		}
	}()

	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM products WHERE id = $1)`,
		productFriend.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to select products by id from db: %w", err)
	}
	if !exists {
		return fmt.Errorf("product with id %d not found", productFriend.ID)
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE products SET name=$1, hobby=$2, price=$3 WHERE id=$4`,
		productFriend.Name,
		productFriend.Hobby,
		productFriend.Price,
		productFriend.ID)
	if err != nil {
		return fmt.Errorf("failed to update product friend in db: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction for UpdateProductFriend: %w", err)
	}
	return nil
}

func (s *Storage) DeleteProductFriend(ctx context.Context, id int) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete product friend from db: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows after delete: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("no product found with the given ID")
	}
	return nil
}

func (s *Storage) GetProductFriend() ([]ProductFriend, error) {
	rows, err := s.db.Query(`SELECT id, name, hobby, price FROM products`)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()

	var product []ProductFriend
	for rows.Next() {
		var p ProductFriend
		if err = rows.Scan(&p.ID, &p.Name, &p.Hobby, &p.Price); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		product = append(product, p)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}
	return product, nil
}
