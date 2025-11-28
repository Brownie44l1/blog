package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Brownie44l1/blog/internal/auth"
	"github.com/Brownie44l1/blog/internal/service"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

type AuthHandler struct {
	userService service.UserService
	jwtSecret   string
}

func NewAuthHandler(userService service.UserService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Username == "" || req.Password == ""{
		respondWithError(w, http.StatusBadRequest, "Username and password cannot be empty")
		return
	}

	user, err := h.userService.RegisterUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUsernameTaken) {
			respondWithError(w, http.StatusConflict, "Username already taken")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "An internal server error occurred")
			return
	}

	tokenString, err := auth.GenerateToken(user.ID, h.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate authentication token")
		return
	}

	response := AuthResponse{
		Token:    tokenString,
		Username: user.Username,
	}
	respondWithJSON(w, http.StatusCreated, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Username == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Username and password cannot be empty")
		return
	}

	user, err := h.userService.Authenticate(req.Username, req.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	tokenString, err := auth.GenerateToken(user.ID, h.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate authentication token")
		return
	}

	response := AuthResponse{
		Token:    tokenString,
		Username: user.Username,
	}

	respondWithJSON(w, http.StatusOK, response)
}