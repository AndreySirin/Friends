package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/AndreySirin/Friends/internal/storage"
)

type CtxValue int

const (
	ctxUserID CtxValue = iota
)

func (s *Server) signUn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeJSONError(w, "Method Not Allowed", nil, http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var inp storage.Registration
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		s.writeJSONError(w, "Invalid request format", err, http.StatusBadRequest)
		return
	}

	if err := inp.Validate(); err != nil {
		s.writeJSONError(w, "Validation failed", err, http.StatusBadRequest)
		return
	}
	if err := s.auth.SingUp(r.Context(), &inp); err != nil {
		s.writeJSONError(w, "Registration failed", err, http.StatusBadRequest)
		return
	}

	s.writeJSONResponse(w, map[string]string{"message": "Registration successful"}, "")
}

func (s *Server) signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeJSONError(w, "Method Not Allowed", nil, http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		s.writeJSONError(w, "Failed to read request body", err, http.StatusBadRequest)
		return
	}
	var inp storage.Auth
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		s.writeJSONError(w, "Invalid request format", err, http.StatusBadRequest)
		return
	}

	if err = inp.Validate(); err != nil {
		s.writeJSONError(w, "Validation failed", err, http.StatusBadRequest)
		return
	}
	refreshToken, token, err := s.auth.SingIn(r.Context(), inp.Email, inp.Password)
	if err != nil {
		s.writeJSONError(w, "Authentication failed", err, http.StatusUnauthorized)
		return
	}
	response := map[string]string{"token": token}
	s.writeJSONResponse(w, response, refreshToken)
}

func (s *Server) refresh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		s.writeJSONError(w, "Missing or invalid refresh token", err, http.StatusBadRequest)
		return
	}
	refreshToken, token, err := s.auth.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		s.writeJSONError(w, "Failed to refresh token", err, http.StatusUnauthorized)
		return
	}
	s.writeJSONResponse(w, map[string]string{"token": token}, refreshToken)
}

func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", fmt.Errorf("authorization header is empty")
	}

	headerParts := strings.Fields(header)
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	if len(headerParts[1]) == 0 {
		return "", fmt.Errorf("token is empty")
	}

	return headerParts[1], nil
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		token, err := getTokenFromRequest(r)
		if err != nil {
			s.writeJSONError(w, "Invalid or missing token", err, http.StatusUnauthorized)
			return
		}

		userId, err := s.auth.ParseToken(r.Context(), token)
		if err != nil {
			s.writeJSONError(w, "Failed to parse token", err, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserID, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
