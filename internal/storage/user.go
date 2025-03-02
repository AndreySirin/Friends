package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

type Registration struct {
	Name     string `json:"name" validate:"required,gte=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (r Registration) Validate() error { return validate.Struct(r) }

type Auth struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (a Auth) Validate() error { return validate.Struct(a) }

func (s *Storage) CreatUser(ctx context.Context, user *User) error {

	_, err := s.db.ExecContext(ctx, "INSERT INTO users (name,email,password,registeredAt)values($1,$2,$3,$4)",
		user.Name, user.Email, user.Password, user.RegisteredAt)
	if err != nil {
		s.lg.Error("Error inserting user into database")
		return err
	}
	return nil
}

func (s *Storage) GetUser(ctx context.Context, email string) (*User, error) {
	var us User
	err := s.db.QueryRowContext(ctx, "SELECT id,name,email,password,registeredAt FROM users WHERE email=$1", email).
		Scan(&us.ID, &us.Name, &us.Email, &us.Password, &us.RegisteredAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &us, nil
}
