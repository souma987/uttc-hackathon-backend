package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"uttc-hackathon-backend/internal/middleware"
	"uttc-hackathon-backend/internal/repository"
	"uttc-hackathon-backend/internal/service"
)

// HandleGet returns an order by ID.
//
// Route
//   - GET /orders/{orderId}
//
// Required Headers
//   - Authorization: Bearer <Firebase ID token>
//
// Success Response
//   - 200 OK
//   - Body: Order
//
// Error Responses
//   - 401 Unauthorized
//   - 404 Not Found
func (h *OrderHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())

	// 2. Parse Order ID
	orderID := r.PathValue("orderId")
	if orderID == "" {
		http.Error(w, "missing order id", http.StatusBadRequest)
		return
	}

	// 3. Call Service
	order, err := h.svc.GetOrder(r.Context(), userID, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}
		log.Printf("get order error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// 4. Return Response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("encode get order response error: %v", err)
	}
}
