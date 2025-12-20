package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"uttc-hackathon-backend/internal/repository"
)

type TranslationService struct {
	vertexRepo VertexGenerativeClient
}

func NewTranslationService(vertexRepo VertexGenerativeClient) *TranslationService {
	return &TranslationService{
		vertexRepo: vertexRepo,
	}
}

type TranslationResponse struct {
	TranslatedTitle        string `json:"translated_title"`
	TranslatedDescription  string `json:"translated_description"`
	DetectedSourceLanguage string `json:"detected_source_language"`
}

func (s *TranslationService) TranslateContent(ctx context.Context, title, description, targetLanguage string) (*TranslationResponse, error) {
	// Standardize target language to uppercase for the prompt logic
	targetLanguage = strings.ToUpper(targetLanguage)

	instruction := `You are an expert AI Localization and Translation engine for a C2C marketplace application. Your goal is to localize product listings based on a specified target language.

### INPUT DATA
You will receive:
1. "title": Product title.
2. "description": Product description.
3. "target_language": The 2-letter ISO code for the output language (e.g., "en", "ja").

### OPERATIONAL RULES
1. **Language Detection & Direction:**
   - Analyze the language of the 'title' and 'description'.
   - If the content is **already** in the <target_language>, return it exactly as is (do not translate).
   - If the content is in a different language, translate it into the <target_language>.

2. **Localization Style:**
   - **If Target is English (EN):** Use a persuasive, concise e-commerce tone.
   - **If Target is Japanese (JA):** Use a polite and natural tone (Desu/Masu style - です/ます調).
   - **General:** Preserve model numbers, sizes, and technical specs exactly.

3. **Output Format:**
   - You must output valid, parseable JSON only.
   - Structure:
     {
       "translated_title": "String",
       "translated_description": "String",
       "detected_source_language": "String (2-letter ISO code)"
     }

### CRITICAL CONSTRAINTS
- Do not output any conversational text. Only the JSON object.
- Strictly adhere to the requested <target_language>.`

	instruction = strings.ReplaceAll(instruction, "<target_language>", targetLanguage)

	prompt := fmt.Sprintf(`### INPUT
{
  "title": "%s",
  "description": "%s",
  "target_language": "%s"
}`, title, description, targetLanguage)

	temperature := float32(0.1)
	config := repository.GenerationConfig{
		SystemInstruction: instruction,
		Temperature:       &temperature,
		JsonResponse:      true,
	}

	// Using gemini-2.0-flash-lite as seen in other services
	respStr, err := s.vertexRepo.GenerateContent(ctx, "gemini-2.0-flash-lite", prompt, config)
	if err != nil {
		log.Printf("failed to generate translation: %v", err)
		return nil, fmt.Errorf("failed to generate translation")
	}

	var response TranslationResponse
	if err := json.Unmarshal([]byte(respStr), &response); err != nil {
		log.Printf("failed to unmarshal translation response: %v, response: %s", err, respStr)
		return nil, fmt.Errorf("failed to parse translation response")
	}
	response.DetectedSourceLanguage = strings.ToLower(response.DetectedSourceLanguage)

	return &response, nil
}
