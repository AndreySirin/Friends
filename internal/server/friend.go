package server

import (
	"context"
	"encoding/json"
	"github.com/AndreySirin/Friends/internal/htmlFile"
	"github.com/AndreySirin/Friends/internal/storage"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"strconv"
)

func (s *Server) priceHandler(w http.ResponseWriter, r *http.Request) {
	data, err := s.methodFriend.GetQueryDB()
	if err != nil {
		http.Error(w, "error getting data from storage", http.StatusInternalServerError)
		return
	}
	path, err := htmlFile.PathHtml("price.html")
	if err != nil {
		http.Error(w, "error loading path price.html", http.StatusInternalServerError)
	}
	menu, err := template.ParseFiles(path)
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
	if err := s.methodFriend.AddProductFriend(context.Background(), prod); err != nil {
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
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	var prod *storage.ProductFriend
	if err := json.NewDecoder(r.Body).Decode(&prod); err != nil {
		http.Error(w, "error decoding product", http.StatusBadRequest)
	}
	if err := s.methodFriend.UpdateProductFriend(context.Background(), prod); err != nil {
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

	err = s.methodFriend.DeleteProductFriend(id)
	if err != nil {
		http.Error(w, "error deleting product", http.StatusInternalServerError)
		return
	}

	if _, err = w.Write([]byte("successful delete")); err != nil {
		http.Error(w, "error deleting product", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}
