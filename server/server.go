package server

import (
	"Friends/storage"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log/slog"
	"net/http"
)

type zeter interface {
	GetZZZ() ([]storage.Product, error)
}

type Server struct {
	log    *slog.Logger
	server *http.Server
	z      zeter
}

func NewServer(log *slog.Logger, addr string, z zeter) *Server {
	s := &Server{
		log: log.With("component", "server"),
		z:   z,
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/friends", s.friendHandler)
			r.Get("/menu", s.priceHandler)
			//	r.Post("/friend",//fixme )
			//TODO написать метод который добавляет данные в базу данных
		})
	})

	s.server = &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return s
}

func (s *Server) friendHandler(w http.ResponseWriter, r *http.Request) {
	home, err := template.ParseFiles("htmlFile/home.html")
	if err != nil {
		http.Error(w, "error loading home", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err = home.Execute(w, nil); err != nil {
		http.Error(w, "error rendering home", http.StatusInternalServerError)
	}
}

func (s *Server) priceHandler(w http.ResponseWriter, r *http.Request) {

	data, err := s.z.GetZZZ()
	if err != nil {
		http.Error(w, "error getting data from storage", http.StatusInternalServerError)
		return
	}
	menu, err := template.ParseFiles("htmlFile/price.html")
	if err != nil {
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

	s.log.Info(fmt.Sprintf("Listening on %s", s.server.Addr))
	return s.server.ListenAndServe()
}
