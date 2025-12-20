package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"uttc-hackathon-backend/internal/service"
)

type TranslationHandler struct {
	translationService *service.TranslationService
}

func NewTranslationHandler(translationService *service.TranslationService) *TranslationHandler {
	return &TranslationHandler{
		translationService: translationService,
	}
}

type TranslateRequest struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	TargetLanguage string `json:"target_language"`
}

// HandleTranslate handles the request to translate product listing content.
// Route: POST /translate
// Required Headers:
//   - Content-Type: application/json
//
// Request Body:
//   - title (string, optional): The title of the product.
//   - description (string, optional): The description of the product.
//   - target_language (string, required): The ISO code of the target language (e.g., "en", "ja").
//
// Response Body (Success - 200 OK):
//   - translated_title (string): The translated title.
//   - translated_description (string): The translated description.
//   - detected_source_language (string): The detected source language code.
//
// Response Body (Error):
//   - 400 Bad Request: Invalid request body or missing required fields.
//   - 500 Internal Server Error: Failed to generate translation.
func (h *TranslationHandler) HandleTranslate(w http.ResponseWriter, r *http.Request) {
	var req TranslateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" && req.Description == "" {
		http.Error(w, "Title or description is required", http.StatusBadRequest)
		return
	}

	if req.TargetLanguage == "" {
		http.Error(w, "target_language is required", http.StatusBadRequest)
		return
	}

	if len(req.TargetLanguage) != 2 {
		http.Error(w, "target_language must be a 2-letter ISO code", http.StatusBadRequest)
		return
	}

	req.TargetLanguage = strings.ToLower(req.TargetLanguage)

	resp, err := h.translationService.TranslateContent(r.Context(), req.Title, req.Description, req.TargetLanguage)
	if err != nil {
		http.Error(w, "Failed to translate content", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
