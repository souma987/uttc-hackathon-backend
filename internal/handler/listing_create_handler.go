package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"uttc-hackathon-backend/internal/middleware"
	"uttc-hackathon-backend/internal/models"
	"uttc-hackathon-backend/internal/service"
)

// HandleCreate creates a new listing.
//
// Route
//   - POST /listings
//
// Required Headers
//   - Authorization: Bearer <Firebase ID token>
//   - Content-Type: application/json
//
// Request Body
//   - title: string (required)
//   - description: string
//   - images: []{url: string} (required)
//   - price: int (required)
//   - quantity: int
//   - item_condition: string (new, excellent, good, not_good, bad)
//
// Success Response
//   - 201 Created
//   - Content-Type: application/json
//   - Body: Listing
func (h *ListingHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())

	var req models.Listing
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.SellerID = userID

	createdListing, err := h.svc.CreateListing(r.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrTitleRequired) ||
			errors.Is(err, service.ErrPriceInvalid) ||
			errors.Is(err, service.ErrNoImages) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("create listing error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(createdListing); err != nil {
		log.Printf("encode create listing response error: %v", err)
	}
}
