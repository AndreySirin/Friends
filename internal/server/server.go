package server

import (
	"context"
	"fmt"
	"github.com/AndreySirin/Friends/internal/htmlFile"
	"github.com/AndreySirin/Friends/internal/servisAuth"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/AndreySirin/Friends/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	log          *slog.Logger
	server       *http.Server
	methodFriend storage.MethFriend
	auth         servisAuth.Authenticate
}

func NewServer(log *slog.Logger, addr string, friend storage.MethFriend, au *servisAuth.Auth) *Server {
	s := &Server{
		log:          log.With("module", "server"),
		methodFriend: friend,
		auth:         au,
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/registration", s.singUn)
			r.Post("/authentication", s.singIn)
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
	path, err := htmlFile.PathHtml("main.html")
	if err != nil {
		http.Error(w, "error loading path home.html", http.StatusInternalServerError)
	}
	home, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "error loading home", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err = home.Execute(w, nil); err != nil {
		http.Error(w, "error rendering home", http.StatusInternalServerError)
	}
}
