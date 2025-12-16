package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
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

	// 1. Auth Check
	authz := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if len(authz) < len(prefix) || authz[:len(prefix)] != prefix {
		http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
		return
	}
	idToken := authz[len(prefix):]

	user, err := h.userSvc.GetCurrentUser(r.Context(), idToken)
	if err != nil {
		log.Printf("get order auth error: %v", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Parse Order ID
	orderID := r.PathValue("orderId")
	if orderID == "" {
		http.Error(w, "missing order id", http.StatusBadRequest)
		return
	}

	// 3. Call Service
	order, err := h.svc.GetOrder(r.Context(), user.ID, orderID)
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
