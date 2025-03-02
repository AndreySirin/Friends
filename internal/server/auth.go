package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AndreySirin/Friends/internal/storage"
	"io"
	"net/http"
	"strings"
)

type CtxValue int

const (
	ctxUserID CtxValue = iota
)

func (s *Server) singUn(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		s.log.Error("singUp error")
		return
	}
	var inp storage.Registration
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		s.log.Error("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = inp.Validate(); err != nil {
		s.log.Error("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = s.auth.SingUp((r.Context()), &inp); err != nil {
		s.log.Error("singUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) singIn(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		s.log.Error("singUp error")
		return
	}
	var inp storage.Auth
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]string{
			"error":   "Invalid request format",
			"details": err.Error(),
		}
		respBytes, _ := json.Marshal(resp)
		w.Write(respBytes)
		return
	}

	if err = inp.Validate(); err != nil {
		s.log.Error("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]string{
			"error":   "Invalid request format",
			"details": err.Error(),
		}
		respBytes, _ := json.Marshal(resp)
		w.Write(respBytes)
		return
	}
	refreshToken, token, err := s.auth.SingIn(r.Context(), inp.Email, inp.Password)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]string{
			"error":   "Invalid request format",
			"details": err.Error(),
		}
		respBytes, _ := json.Marshal(resp)
		w.Write(respBytes)
		return
	}
	response, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		s.log.Error("singIn", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)

}

func (s *Server) refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		s.log.Error("refresh", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	refreshToken, token, err := s.auth.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		s.log.Error("refresh", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		s.log.Error("refresh", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)
		if err != nil {
			s.log.Error("authMiddleware", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, err := s.auth.ParseToken(r.Context(), token)
		if err != nil {
			s.log.Error("authMiddleware", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserID, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
