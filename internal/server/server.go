package server

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/AndreySirin/Friends/internal/servisAuth"
	"github.com/AndreySirin/Friends/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	log          *slog.Logger
	server       *http.Server
	methodFriend storage.StorageFriend
	auth         servisAuth.Authenticate
}

func NewServer(log *slog.Logger, addr string, friend storage.StorageFriend, au *servisAuth.Auth) *Server {
	s := &Server{
		log:          log.With("module", "server"),
		methodFriend: friend,
		auth:         au,
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/registration", s.signUn)
			r.Post("/authentication", s.signIn)
			r.Get("/refreshToken", s.refresh)
			r.Get("/main", s.mainHandler)

			r.With(s.authMiddleware).Group(func(r chi.Router) {
				r.Get("/price", s.priceHandler)
				r.Post("/AddNewUser", s.AddUserHandler)
				r.Put("/UpdateUser", s.UpdateUser)
				r.Delete("/DeleteUser/{id}", s.DeleteUserHandler)
			})
		})
	})

	s.server = &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return s
}

func (s *Server) Run() error {
	s.log.Info(fmt.Sprintf("Listening on %s", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *Server) ShutDown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.log.Error("error shutting down server", "error", err)
		return err
	}
	return nil
}

func (s *Server) mainHandler(w http.ResponseWriter, r *http.Request) {
	//path, err := config.PathHtml("main.html")
	//if err != nil {
	//	s.writeJSONError(w, "Error loading path main.html", err, http.StatusInternalServerError)
	//	return
	//}
	home, err := template.ParseFiles("/root/htmlFile/main.html")
	if err != nil {
		s.writeJSONError(w, "Error loading main template", err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err = home.Execute(w, nil); err != nil {
		s.writeJSONError(w, "Error rendering main template", err, http.StatusInternalServerError)
	}
}

func (s *Server) writeJSONError(w http.ResponseWriter, message string, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := map[string]string{
		"error":   message,
		"details": err.Error(),
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		s.log.Error("error writing response", "error", err)
	}
}

func (s *Server) writeJSONResponse(w http.ResponseWriter, data map[string]string, refreshToken string) {
	w.Header().Set("Content-Type", "application/json")
	if refreshToken != "" {
		w.Header().Add("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	}
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		s.log.Error("error writing response", "error", err)
	}
}
