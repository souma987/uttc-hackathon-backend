package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"uttc-hackathon-backend/internal/service"
)

// HandleGetProfile returns a user's profile (public info).
//
// Route: GET /users/{userId}/profile
//
// Success Response:
//   - 200 OK
//   - Content-Type: application/json
//   - Body: UserProfile (id, name, avatar_url)
//
// Error Responses:
//   - 400 Bad Request: missing user id
//   - 404 Not Found: user not found
//   - 500 Internal Server Error
func (h *UserHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("userId")
	if id == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	profile, err := h.svc.GetUserProfile(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		log.Printf("failed to get user profile: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
