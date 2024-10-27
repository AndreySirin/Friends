package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AndreySirin/Friends/storage"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type zeter interface {
	GetZZZ() ([]storage.ProductFriend, error)
	AddProductFriend(context.Context, *storage.ProductFriend) error
	DeleteProductFriend(int) error
	UpdateProductFriend(context.Context, *storage.ProductFriend) error
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
			r.Post("/AddNewUser", s.AddUserHandler)
			r.Post("/UpdateUser", s.UpdateUser)
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

func (s *Server) AddUserHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      prod.ID,
		"message": "Product added successfully",
	}); err != nil {
		http.Error(w, "error encoding product", http.StatusInternalServerError)
	}
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	var prod *storage.ProductFriend
	if err := json.NewDecoder(r.Body).Decode(&prod); err != nil {
		http.Error(w, "error decoding product", http.StatusBadRequest)
	}
	if err := s.z.UpdateProductFriend(context.Background(), prod); err != nil {
		http.Error(w, "error updating product", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      prod.ID,
		"message": "Product updated successfully",
	}); err != nil {
		http.Error(w, "error encoding product", http.StatusInternalServerError)
	}
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
	if _, err = w.Write([]byte("successful delete")); err != nil {
		http.Error(w, "error deleting product", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
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
