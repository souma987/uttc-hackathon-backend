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
//   - is_active: bool (optional, default false/draft)
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

	var req struct {
		Title         string                `json:"title"`
		Description   string                `json:"description"`
		Images        []models.ListingImage `json:"images"`
		Price         int                   `json:"price"`
		Quantity      int                   `json:"quantity"`
		ItemCondition models.ItemCondition  `json:"item_condition"`
		IsActive      bool                  `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	status := models.ListingStatusDraft
	if req.IsActive {
		status = models.ListingStatusActive
	}

	listing := &models.Listing{
		SellerID:      userID,
		Title:         req.Title,
		Description:   req.Description,
		Images:        req.Images,
		Price:         req.Price,
		Quantity:      req.Quantity,
		ItemCondition: req.ItemCondition,
		Status:        status,
	}

	createdListing, err := h.svc.CreateListing(r.Context(), listing)
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
