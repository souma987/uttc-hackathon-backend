package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"uttc-hackathon-backend/internal/middleware"
	"uttc-hackathon-backend/internal/service"
)

type createMessageRequest struct {
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

// HandleCreate creates a new message.
//
// Route:
//   - POST /messages
//
// Required Headers:
//   - Authorization: Bearer <Firebase ID token>
//
// Request Body:
//   - receiver_id: string (required)
//   - content: string (required)
//
// Success Response:
//   - 201 Created
//   - Body: Message
//
// Error Responses:
//   - 400 Bad Request: Invalid body or missing content
//   - 401 Unauthorized: Missing or invalid token
//   - 500 Internal Server Error: Database error
func (h *MessageHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())

	var req createMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	msg, err := h.svc.CreateMessage(r.Context(), userID, req.ReceiverID, req.Content)
	if err != nil {
		if errors.Is(err, service.ErrContentRequired) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// FIXME: handle FK violation
		log.Printf("create message error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		log.Printf("encode create message response error: %v", err)
	}
}
