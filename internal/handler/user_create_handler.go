package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"uttc-hackathon-backend/internal/service"
)

type CreateUserRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarURL string `json:"avatar_url"`
}

// HandleCreate registers a new user.
//
// Route
//   - POST /users
//
// Required Headers
//   - Content-Type: application/json
//
// Request
//   - Body (application/json):
//     {
//     "name": string (optional),
//     "email": string (required),
//     "password": string (required, 8-4096 characters),
//     "avatar_url": string (optional)
//     }
//
// Success Response
//   - 201 Created
//   - Content-Type: application/json
//   - Body:
//     {
//     "id": string,
//     "name": string,
//     "email": string
//     }
//
// Errors
//   - 400 Bad Request: invalid JSON, missing required fields, or sign-up failure
//   - 405 Method Not Allowed: request method is not POST
//   - 500 Internal Server Error: failed to encode response
func (h *UserHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.svc.SignUp(r.Context(), req.Name, req.Email, req.Password, req.AvatarURL)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPasswordLength) {
			http.Error(w, "password must be between 8 and 4096 characters", http.StatusBadRequest)
			return
		}

		// Do not expose internal error details
		log.Printf("sign up error: %v", err)
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		// Do not expose internal error details
		log.Printf("encode create user response error: %v", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}
