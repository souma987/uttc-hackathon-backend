package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"uttc-hackathon-backend/internal/service"
)

type ListingHandler struct {
	svc *service.ListingService
}

func NewListingHandler(svc *service.ListingService) *ListingHandler {
	return &ListingHandler{svc: svc}
}

// HandleFeed returns a list of active listings.
//
// Route
//   - GET /listings/feed
//
// Query Parameters
//   - limit: int (optional, default 20, max 100)
//   - offset: int (optional, default 0)
//
// Success Response
//   - 200 OK
//   - Content-Type: application/json
//   - Body: []Listing
func (h *ListingHandler) HandleFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			limit = v
		}
	}

	offset := 0
	if offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil {
			offset = v
		}
	}

	listings, err := h.svc.GetFeed(r.Context(), limit, offset)
	if err != nil {
		log.Printf("get listings feed error: %v", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(listings); err != nil {
		log.Printf("encode listings response error: %v", err)
	}
}
