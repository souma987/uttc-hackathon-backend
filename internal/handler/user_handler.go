package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"uttc-hackathon-backend/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// We will return the created user directly (id, name, email)

func (h *UserHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	user, err := h.svc.FetchUser(r.Context(), id)
	if err != nil {
		// Do not expose internal error details
		log.Printf("fetch user error: %v", err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		// Do not expose internal error details
		log.Printf("encode user response error: %v", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}

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
//     "id": string,   // Firebase UID
//     "name": string,
//     "email": string
//     }
//
// Errors
//   - 401 Unauthorized: missing or invalid Authorization header, or invalid/expired token
//   - 404 Not Found: token is valid but no corresponding user exists in DB
//   - 500 Internal Server Error: failed to encode response
func (h *UserHandler) HandleMe(w http.ResponseWriter, r *http.Request) {
	authz := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if len(authz) < len(prefix) || authz[:len(prefix)] != prefix {
		http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
		return
	}
	idToken := authz[len(prefix):]
	if idToken == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	user, err := h.svc.GetCurrentUser(r.Context(), idToken)
	if err != nil {
		// Distinguish between not found (user absent in DB) vs invalid token/other
		// Our repo returns "user not found" for missing users
		if err.Error() == "user not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		log.Printf("me error: %v", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("encode me response error: %v", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}

// HandleCreate registers a new user via POST /users.
//
// Request
//   - Method: POST
//   - Headers: Content-Type: application/json
//   - Body (JSON object):
//     {
//     "name": string,      // optional display name
//     "email": string,     // required
//     "password": string   // required
//     }
//
// Response
//   - On success: 201 Created with application/json body representing the created user:
//     {
//     "id": string,   // Firebase UID
//     "name": string, // echoes the provided name (may be empty)
//     "email": string
//     }
//
// Errors
// - 400 Bad Request: invalid JSON, missing required fields (email/password), or sign-up failure
// - 405 Method Not Allowed: if the request method is not POST
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

	user, err := h.svc.SignUp(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
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
