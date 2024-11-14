package server

import (
	"Friends/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

type zeter interface {
	GetZZZ() ([]storage.ProductFriend, error)
	AddProductFriend(context.Context, *storage.ProductFriend) error
	DeleteProductFriend(int) error
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
			r.Get("/main", s.mainHandler)
			r.Get("/price", s.priceHandler)
			r.Post("/PostNewUser", s.PostUserHandler)
			r.Delete("/DeleteUser/{id}", s.DeleteUserHandler)
		})
	})

	s.server = &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return s
}

func (s *Server) mainHandler(w http.ResponseWriter, r *http.Request) {
	home, err := template.ParseFiles("htmlFile/main.html")
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

func (s *Server) PostUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	var prod *storage.ProductFriend
	if err := json.NewDecoder(r.Body).Decode(&prod); err != nil {
		http.Error(w, "error decoding product", http.StatusBadRequest)
	}
	if err := s.z.AddProductFriend(context.Background(), prod); err != nil {
		http.Error(w, "error adding product", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode((map[string]int{
		"id": prod.ID,
	}))

}
func (s *Server) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(productID)
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	err = s.z.DeleteProductFriend(id)
	if err != nil {
		http.Error(w, "error deleting product", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("successful delete"))
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) Run() error {

	s.log.Info(fmt.Sprintf("Listening on %s", s.server.Addr))
	return s.server.ListenAndServe()
}
