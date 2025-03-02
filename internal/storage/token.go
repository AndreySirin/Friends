package storage

import (
	"context"
	"time"
)

type RefreshToken struct {
	ID           int    `json:"id"`
	IdUser       int    `json:"id_user"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    time.Time
}

func (s *Storage) CreateRefreshToken(ctx context.Context, r *RefreshToken) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO refresh_tokens (id_user,refresh_token,expires_at) values ($1,$2,$3)",
		r.IdUser, r.RefreshToken, r.ExpiresAt)
	if err != nil {
		s.lg.Error("ошибка при сохранении токена ", err)
		return err
	}
	return nil
}

func (s *Storage) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	var t RefreshToken

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.lg.Error("Ошибка при создании транзакции:", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = tx.QueryRowContext(ctx,
		"SELECT id_user, refresh_token, expires_at FROM refresh_tokens WHERE refresh_token=$1", token).
		Scan(&t.IdUser, &t.RefreshToken, &t.ExpiresAt)
	if err != nil {
		s.lg.Error("Ошибка при получении refresh-токена:", err)
		return nil, err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE id_user=$1", t.IdUser)
	if err != nil {
		s.lg.Error("Ошибка при удалении refresh-токена:", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		s.lg.Error("Ошибка при коммите транзакции:", err)
		return nil, err
	}

	return &t, nil
}
