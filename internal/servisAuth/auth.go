package servisAuth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/AndreySirin/Friends/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Authenticate interface {
	SingUp(context.Context, *storage.Registration) error
	SingIn(ctx context.Context, password, email string) (string, string, error)
	newRefreshToken() (string, error)
	generateTokens(ctx context.Context, userId int) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	ParseToken(ctx context.Context, Token string) (int64, error)
}

type Hasher interface {
	Hash(string) (string, error)
}

type Auth struct {
	auth storage.StorageUser
	hash Hasher
	Salt []byte
}

func NewAuth(auth storage.StorageUser, hash Hasher, SecretSalt []byte) *Auth {
	return &Auth{
		auth: auth,
		hash: hash,
		Salt: SecretSalt,
	}
}

func (a *Auth) SingUp(ctx context.Context, r *storage.Registration) error {
	password, err := a.hash.Hash(r.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := storage.User{
		Name:         r.Name,
		Email:        r.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	if err = a.auth.CreatUser(ctx, &user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (a *Auth) SingIn(ctx context.Context, email, password string) (string, string, error) {
	user, err := a.auth.GetUser(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", fmt.Errorf("invalid password")
	}

	return a.generateTokens(ctx, user.ID)
}

func (a *Auth) newRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func (a *Auth) generateTokens(ctx context.Context, userId int) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(int64(userId), 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	})
	token, err := t.SignedString(a.Salt)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate token: %w", err)
	}
	refreshToken, err := a.newRefreshToken()
	if err != nil {
		return "", "", err
	}

	RT := storage.RefreshToken{
		IdUser:       userId,
		RefreshToken: refreshToken,
		ExpiresAt:    jwt.NewNumericDate(time.Now().Add(24 * time.Hour)).Time,
	}
	err = a.auth.CreateRefreshToken(ctx, &RT)
	if err != nil {
		return "", "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return refreshToken, token, nil
}

func (a *Auth) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	RT, err := a.auth.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to get refresh token: %w", err)
	}
	if RT.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("refresh token is expired")
	}
	return a.generateTokens(ctx, RT.IdUser)
}

func (a *Auth) ParseToken(ctx context.Context, Token string) (int64, error) {
	token, err := jwt.Parse(Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.Salt, nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.ParseInt(subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid subject: %w", err)
	}

	return id, nil
}
