package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// HandleGetMessages fetches all messages between the current user and the specified user.
//
// Route:
//   - GET /messages/with/{userid}
//
// Required Headers:
//   - Authorization: Bearer <Firebase ID token>
//
// Success Response:
//   - 200 OK
//   - Body: []Message
//
// Error Responses:
//   - 401 Unauthorized: Missing or invalid token
//   - 500 Internal Server Error: Database error
func (h *MessageHandler) HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authz := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if len(authz) < len(prefix) || authz[:len(prefix)] != prefix {
		http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
		return
	}
	idToken := authz[len(prefix):]

	user, err := h.userSvc.GetCurrentUser(r.Context(), idToken)
	if err != nil {
		log.Printf("get messages auth error: %v", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	otherUserID := r.PathValue("userid")
	if otherUserID == "" {
		http.Error(w, "userid is required", http.StatusBadRequest)
		return
	}

	messages, err := h.svc.GetMessages(r.Context(), user.ID, otherUserID)
	if err != nil {
		log.Printf("get messages error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Printf("encode get messages response error: %v", err)
	}
}
