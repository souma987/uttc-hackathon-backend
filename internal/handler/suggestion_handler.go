package handler

import (
	"encoding/json"
	"net/http"
	"uttc-hackathon-backend/internal/service"
)

type SuggestionHandler struct {
	suggestionService *service.SuggestionService
}

func NewSuggestionHandler(suggestionService *service.SuggestionService) *SuggestionHandler {
	return &SuggestionHandler{
		suggestionService: suggestionService,
	}
}

type NewListingSuggestionRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Condition   string `json:"condition"`
	Language    string `json:"language"`
}

type NewListingSuggestionResponse struct {
	Suggestions []string `json:"suggestions"`
}

// HandleGetSuggestion handles the request to generate a listing suggestion.
// Route: POST /suggestions/newListing
// Required Headers:
//   - Content-Type: application/json
//   - Authorization: Bearer <token>
//
// Request Body:
//   - title (string, required): The title of the listing.
//   - description (string, required): The description of the listing.
//   - condition (string, required): The condition of the listing.
//   - language (string, required): The language for the suggestion ("ja" or "en").
//
// Response Body (Success - 200 OK):
//   - suggestions ([]string): The generated suggestions for the listing.
//
// Response Body (Error):
//   - 400 Bad Request: Invalid request body or missing required fields.
//   - 500 Internal Server Error: Failed to generate suggestion.
func (h *SuggestionHandler) HandleGetSuggestion(w http.ResponseWriter, r *http.Request) {
	var req NewListingSuggestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Description == "" || req.Condition == "" || req.Language == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	suggestions := h.suggestionService.GetListingSuggestion(r.Context(), req.Title, req.Description, req.Condition, req.Language)

	resp := NewListingSuggestionResponse{
		Suggestions: suggestions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
