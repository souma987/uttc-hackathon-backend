package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"uttc-hackathon-backend/internal/middleware"
)

// HandleMe returns the current authenticated user.
//
// Route
//   - GET /me
//
// Required Headers
//   - Authorization: Bearer <Firebase ID token>
//
// Request
//   - No body
//
// Success Response
//   - 200 OK
//   - Content-Type: application/json
//   - Body:
//     {
//     "id": string,
//     "name": string,
//     "email": string,
//     "avatar_url": string
//     }
//
// Errors
//   - 401 Unauthorized: missing or invalid Authorization header, or invalid/expired token
//   - 404 Not Found: valid token but user missing in DB
//   - 500 Internal Server Error: failed to encode response
func (h *UserHandler) HandleMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	user, err := h.svc.GetUser(r.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		log.Printf("me error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("encode me response error: %v", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}
