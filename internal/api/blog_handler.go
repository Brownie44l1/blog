package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Brownie44l1/blog/internal/service"

	"github.com/go-chi/chi/v5"
)

type BlogHandler struct {
	blogService *service.BlogService
}

func NewBlogHandler(blogService *service.BlogService) *BlogHandler {
	return &BlogHandler{blogService: blogService}
}

func (h *BlogHandler) Publish(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"user_id"`
		Title string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		log.Fatalf("Request error: %v", err)
		return 
	}

	user, err := h.blogService.Publish(req.UserID, req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		log.Fatalf("JSON encoding error: %v", err)
		return
	}
}


func (h *BlogHandler) GetBlog(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	blog, err := h.blogService.GetBlog(id)
	if err != nil {
		http.Error(w, "invalid blog id", http.StatusNotFound)
		log.Fatalf("Blog not found: %v", err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(blog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("JSON encoding error: %v", err)
		return 
	}
}
