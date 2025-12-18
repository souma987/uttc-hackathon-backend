package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"uttc-hackathon-backend/internal/middleware"
)

// HandleGetMyOrders returns all orders for the current user (buyer or seller).
//
// Route
//   - GET /orders/my
//
// Required Headers
//   - Authorization: Bearer <Firebase ID token>
//
// Success Response
//   - 200 OK
//   - Body: []Order
//
// Error Responses
//   - 401 Unauthorized
//   - 500 Internal Server Error
func (h *OrderHandler) HandleGetMyOrders(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())

	// Call Service
	orders, err := h.svc.GetOrdersByUser(r.Context(), userID)
	if err != nil {
		log.Printf("get my orders error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Return Response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		log.Printf("encode get my orders response error: %v", err)
	}
}
