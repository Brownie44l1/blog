package api

import (
	"encoding/json"
	"net/http"
)

// respondWithJSON sends a JSON response with the given status code and data
func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// respondWithError sends a JSON error response with the given status code and message
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// getUserIDFromContext extracts the user ID from the request context (set by auth middleware)
func getUserIDFromContext(r *http.Request) (int64, error) {
	userID := r.Context().Value("userID")
	if userID == nil {
		return 0, nil
	}
	id, ok := userID.(int64)
	if !ok {
		return 0, nil
	}
	return id, nil
}