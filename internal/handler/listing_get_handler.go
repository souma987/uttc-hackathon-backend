package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"uttc-hackathon-backend/internal/service"
)

// HandleGetListing returns a specific listing by ID.
//
// Route
//   - GET /listings/{id}
//
// Success Response
//   - 200 OK
//   - Content-Type: application/json
//   - Body: Listing
func (h *ListingHandler) HandleGetListing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing listing id", http.StatusBadRequest)
		return
	}

	listing, err := h.svc.GetListing(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrListingNotFound) {
			http.Error(w, "listing not found", http.StatusNotFound)
			return
		}
		log.Printf("get listing error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(listing); err != nil {
		log.Printf("encode listing response error: %v", err)
	}
}
