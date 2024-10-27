package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log/slog"
	"net/http"
)

type Server struct {
	log    *slog.Logger
	server *http.Server
}

func NewServer(log *slog.Logger, addr string) *Server {
	s := &Server{
		log: log.With("component", "server"),
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

	menu, err := template.ParseFiles("htmlFile/price.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "error loading price", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = menu.Execute(w, nil)
	if err != nil {
		http.Error(w, "error rendering price", http.StatusInternalServerError)
	}

}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}
