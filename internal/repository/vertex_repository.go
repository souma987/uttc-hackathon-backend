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

func (r *VertexRepository) GenerateContent(ctx context.Context, modelName string, prompt string, systemInstruction string) (string, error) {
	var config *genai.GenerateContentConfig
	if systemInstruction != "" {
		config = &genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{{Text: systemInstruction}},
			},
		}
	}

	contents := genai.Text(prompt)

	resp, err := r.client.Models.GenerateContent(ctx, modelName, contents, config)
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
