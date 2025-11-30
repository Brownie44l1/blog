package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Brownie44l1/blog/internal/middleware"
	"github.com/Brownie44l1/blog/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile handles GET /users/{id}
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract ID from path: /users/123
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	idStr := strings.Split(path, "/")[0]

	// Check if it's not a sub-path like /users/123/blogs
	if strings.Contains(path, "/") {
		respondWithError(w, http.StatusNotFound, "Not found")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	profile, err := h.userService.GetUserProfile(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		log.Printf("Error retrieving user profile: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user profile")
		return
	}

	respondWithJSON(w, http.StatusOK, profile)
}

// GetMe handles GET /users/me (authenticated user's profile)
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
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

	profile, err := h.userService.GetUserProfile(userID)
	if err != nil {
		log.Printf("Error retrieving user profile: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user profile")
		return
	}

	respondWithJSON(w, http.StatusOK, profile)
}
