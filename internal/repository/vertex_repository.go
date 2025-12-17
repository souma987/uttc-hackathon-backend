package repository

import (
	"context"

	"cloud.google.com/go/vertexai/genai"
)

type VertexRepository struct {
	client *genai.Client
}

func NewVertexRepository(client *genai.Client) *VertexRepository {
	return &VertexRepository{client: client}
}

func (r *VertexRepository) GenerateContent(ctx context.Context, modelName string, prompt string, systemInstruction string) (string, error) {
	model := r.client.GenerativeModel(modelName)
	if systemInstruction != "" {
		model.SystemInstruction = &genai.Content{
			Role:  "",
			Parts: []genai.Part{genai.Text(systemInstruction)},
		}
	}
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", nil // Or an error indicating no response
	}

	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			result += string(txt)
		}
	}

	return result, nil
}
