package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AndreySirin/Friends/internal/htmlFile"
	"github.com/AndreySirin/Friends/internal/storage"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"strconv"
)

func (s *Server) priceHandler(w http.ResponseWriter, r *http.Request) {
	data, err := s.methodFriend.GetProductFriend()
	if err != nil {
		s.writeJSONError(w, "Error getting data from storage", err, http.StatusInternalServerError)
		return
	}
	path, err := htmlFile.PathHtml("price.html")
	if err != nil {
		s.writeJSONError(w, "Error loading path price.html", err, http.StatusInternalServerError)
		return
	}
	menu, err := template.ParseFiles(path)
	if err != nil {
		s.writeJSONError(w, "Error loading price template", err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = menu.Execute(w, data)
	if err != nil {
		s.writeJSONError(w, "Error rendering price", err, http.StatusInternalServerError)
	}
}

func (s *Server) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var prod *storage.ProductFriend
	if err := json.NewDecoder(r.Body).Decode(&prod); err != nil {
		s.writeJSONError(w, "Error decoding product", err, http.StatusBadRequest)
		return
	}
	if err := s.methodFriend.AddProductFriend(context.Background(), prod); err != nil {
		s.writeJSONError(w, "Error adding product", err, http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"id":      fmt.Sprintf("%d", prod.ID),
		"message": "Product added successfully",
	}
	s.writeJSONResponse(w, data, "")
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	var prod *storage.ProductFriend
	if err := json.NewDecoder(r.Body).Decode(&prod); err != nil {
		s.writeJSONError(w, "Error decoding product", err, http.StatusBadRequest)
		return
	}
	if err := s.methodFriend.UpdateProductFriend(context.Background(), prod); err != nil {
		s.writeJSONError(w, "Error updating product", err, http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"id":      fmt.Sprintf("%d", prod.ID),
		"message": "Product updated successfully",
	}
	s.writeJSONResponse(w, data, "")
}

func (s *Server) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(productID)
	if err != nil {
		s.writeJSONError(w, "Invalid product ID", err, http.StatusBadRequest)
		return
	}

	err = s.methodFriend.DeleteProductFriend(context.Background(), id)
	if err != nil {
		s.writeJSONError(w, "Error deleting product", err, http.StatusInternalServerError)
		return
	}

	s.writeJSONResponse(w, map[string]string{
		"message": "Product deleted successfully",
	}, "")
}
