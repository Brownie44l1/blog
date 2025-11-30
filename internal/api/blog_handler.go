package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Brownie44l1/blog/internal/middleware"
	"github.com/Brownie44l1/blog/internal/models"
	"github.com/Brownie44l1/blog/internal/service"
)

type BlogHandler struct {
	blogService service.BlogService
}

func NewBlogHandler(blogService service.BlogService) *BlogHandler {
	return &BlogHandler{
		blogService: blogService,
	}
}

type CreateBlogRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CreateBlog handles POST /blogs/create
func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("❌ Failed to get user ID from context")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateBlogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	blog := &models.Blog{
		UserId:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := h.blogService.Create(blog); err != nil {
		log.Printf("Error creating blog: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, blog)
}

// GetBlog handles GET /blogs/{id}
func (h *BlogHandler) GetBlog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract ID from path: /blogs/123
	path := strings.TrimPrefix(r.URL.Path, "/blogs/")
	idStr := strings.Split(path, "/")[0]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid blog ID")
		return
	}

	blog, err := h.blogService.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			respondWithError(w, http.StatusNotFound, "Blog not found")
			return
		}
		log.Printf("Error retrieving blog: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve blog")
		return
	}

	respondWithJSON(w, http.StatusOK, blog)
}

// GetMyBlogs handles GET /blogs/me (get blogs for the authenticated user)
func (h *BlogHandler) GetMyBlogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("❌ Failed to get user ID from context")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	blogs, err := h.blogService.GetByUserID(userID)
	if err != nil {
		log.Printf("Error retrieving user blogs: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve blogs")
		return
	}

	respondWithJSON(w, http.StatusOK, blogs)
}

// GetUserBlogs handles GET /users/{userId}/blogs
func (h *BlogHandler) GetUserBlogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract userID from path: /users/123/blogs
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	userIDStr := strings.Split(path, "/")[0]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	blogs, err := h.blogService.GetByUserID(userID)
	if err != nil {
		log.Printf("Error retrieving user blogs: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve blogs")
		return
	}

	respondWithJSON(w, http.StatusOK, blogs)
}

// DeleteBlog handles DELETE /blogs/{id}
func (h *BlogHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("❌ Failed to get user ID from context")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Extract ID from path: /blogs/123
	path := strings.TrimPrefix(r.URL.Path, "/blogs/")
	idStr := strings.Split(path, "/")[0]

	blogID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid blog ID")
		return
	}

	if err := h.blogService.Delete(blogID, userID); err != nil {
		if strings.Contains(err.Error(), "no blog found") {
			respondWithError(w, http.StatusNotFound, "Blog not found or unauthorized")
			return
		}
		log.Printf("Error deleting blog: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to delete blog")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Blog deleted successfully"})
}

// ListBlogs handles GET /blogs
func (h *BlogHandler) ListBlogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var limit, offset int64 = 10, 0

	if limitStr != "" {
		parsedLimit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil || parsedLimit < 1 {
			respondWithError(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
		limit = parsedLimit
	}

	if offsetStr != "" {
		parsedOffset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil || parsedOffset < 0 {
			respondWithError(w, http.StatusBadRequest, "Invalid offset parameter")
			return
		}
		offset = parsedOffset
	}

	blogs, err := h.blogService.ListAll(limit, offset)
	if err != nil {
		log.Printf("Error listing blogs: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve blogs")
		return
	}

	respondWithJSON(w, http.StatusOK, blogs)
}

// SearchBlogs handles GET /blogs/search?q=query
func (h *BlogHandler) SearchBlogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query().Get("q")

	if strings.TrimSpace(query) == "" {
		respondWithError(w, http.StatusBadRequest, "Search query is required")
		return
	}

	blogs, err := h.blogService.Search(query)
	if err != nil {
		log.Printf("Error searching blogs: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to search blogs")
		return
	}

	respondWithJSON(w, http.StatusOK, blogs)
}
