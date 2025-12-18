package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"uttc-hackathon-backend/internal/middleware"
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
	userID := middleware.GetUserIDFromContext(r.Context())

	otherUserID := r.PathValue("userid")
	if otherUserID == "" {
		http.Error(w, "userid is required", http.StatusBadRequest)
		return
	}

	messages, err := h.svc.GetMessages(r.Context(), userID, otherUserID)
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
