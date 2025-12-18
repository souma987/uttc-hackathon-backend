package repository

import (
	"context"

	"google.golang.org/genai"
)

type VertexRepository struct {
	client *genai.Client
}

func NewVertexRepository(client *genai.Client) *VertexRepository {
	return &VertexRepository{client: client}
}

type GenerationConfig struct {
	Temperature       *float32
	JsonResponse      bool
	SystemInstruction string
}

func (r *VertexRepository) GenerateContent(ctx context.Context, modelName string, prompt string, config GenerationConfig) (string, error) {
	genaiConfig := &genai.GenerateContentConfig{}

	if config.SystemInstruction != "" {
		genaiConfig.SystemInstruction = &genai.Content{
			Parts: []*genai.Part{{Text: config.SystemInstruction}},
		}
	}
	if config.Temperature != nil {
		genaiConfig.Temperature = config.Temperature
	}
	if config.JsonResponse {
		genaiConfig.ResponseMIMEType = "application/json"
	}

	contents := genai.Text(prompt)

	resp, err := r.client.Models.GenerateContent(ctx, modelName, contents, genaiConfig)
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", nil
	}

	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			result += part.Text
		}
	}

	return result, nil
}
