package server

import (
	"Friends/storage"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log/slog"
	"net/http"
)

const (
	address  = "127.0.0.1:5432"
	username = "myuser"
	password = "mypassword"
	database = "mydatabase"
)

type Server struct {
	log     *slog.Logger
	server  *http.Server
	storage *storage.Storage
}

func NewServer(log *slog.Logger, addr string) *Server {
	stor, _ := storage.New(log, username, password, address, database)

	s := &Server{
		log:     log.With("component", "server"),
		storage: stor,
	}
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/friends", s.frienHandler)
			r.Get("/menu", s.priceHandler)
		})
	})

	s.server = &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return s
}

func (s *Server) frienHandler(w http.ResponseWriter, r *http.Request) {

	home, err := template.ParseFiles("htmlFile/home.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "error loading home", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = home.Execute(w, nil)
	if err != nil {
		http.Error(w, "error rendering home", http.StatusInternalServerError)
	}

}

func (s *Server) priceHandler(w http.ResponseWriter, r *http.Request) {

	data, err := s.storage.GetZZZ()
	if err != nil {
		http.Error(w, "error getting data from storage", http.StatusInternalServerError)
		return
	}

	menu, err := template.ParseFiles("htmlFile/price.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "error loading price", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = menu.Execute(w, data)
	if err != nil {
		http.Error(w, "error rendering price", http.StatusInternalServerError)
	}

}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}
