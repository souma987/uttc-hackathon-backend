package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"uttc-hackathon-backend/internal/middleware"
	"uttc-hackathon-backend/internal/models"
	"uttc-hackathon-backend/internal/repository"
	"uttc-hackathon-backend/internal/service"
)

type OrderHandler struct {
	svc     *service.OrderService
	userSvc *service.UserService
}

func NewOrderHandler(svc *service.OrderService, userSvc *service.UserService) *OrderHandler {
	return &OrderHandler{
		svc:     svc,
		userSvc: userSvc,
	}
}

// HandleCreate creates a new order.
//
// Route
//   - POST /orders
//
// Required Headers
//   - Authorization: Bearer <Firebase ID token>
//
// Request Body
//   - listing_id: string (required)
//   - quantity: int (required)
//
// Success Response
//   - 201 Created
//   - Body: Order
func (h *OrderHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())

	var req models.Order
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	createdOrder, err := h.svc.CreateOrder(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrQuantityInvalid) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, repository.ErrListingNotFound) || errors.Is(err, repository.ErrListingNotActive) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, repository.ErrInsufficientStock) {
			http.Error(w, err.Error(), http.StatusConflict) // Or 422
			return
		}

		log.Printf("create order error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(createdOrder); err != nil {
		log.Printf("encode create order response error: %v", err)
	}
}
