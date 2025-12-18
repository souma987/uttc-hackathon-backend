package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"uttc-hackathon-backend/internal/middleware"
)

// HandleGetConversations returns a list of conversations (latest message per user).
//
// Route:
//   - GET /messages/conversations
//
// Required Headers:
//   - Authorization: Bearer <Firebase ID token>
//
// Success Response:
//   - 200 OK
//   - Body: []Conversation
//
// Error Responses:
//   - 401 Unauthorized
//   - 500 Internal Server Error
func (h *MessageHandler) HandleGetConversations(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())

	conversations, err := h.svc.GetConversations(r.Context(), userID)
	if err != nil {
		log.Printf("failed to get conversations: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(conversations); err != nil {
		log.Printf("encode conversations error: %v", err)
	}
}
