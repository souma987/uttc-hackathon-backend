package service

import (
	"context"
	"testing"

	"uttc-hackathon-backend/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVertexGenerativeClient struct {
	mock.Mock
}

func (m *MockVertexGenerativeClient) GenerateContent(ctx context.Context, modelName string, prompt string, config repository.GenerationConfig) (string, error) {
	args := m.Called(ctx, modelName, prompt, config)
	return args.String(0), args.Error(1)
}

func TestSuggestionService_GetListingSuggestion(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		condition   string
		language    string
		mockSetup   func(*MockVertexGenerativeClient)
		want        []string
	}{
		{
			name:        "Success - English",
			title:       "iPhone 13",
			description: "Good condition",
			condition:   "Used",
			language:    "en",
			mockSetup: func(m *MockVertexGenerativeClient) {
				m.On("GenerateContent", mock.Anything, "gemini-2.5-flash", mock.Anything, mock.MatchedBy(func(c repository.GenerationConfig) bool {
					return c.JsonResponse == true
				})).Return(`["Storage", "Color", "Battery Health"]`, nil)
			},
			want: []string{"Storage", "Color", "Battery Health"},
		},
		{
			name:        "Success - Japanese",
			title:       "iPad",
			description: "美品",
			condition:   "Used",
			language:    "ja",
			mockSetup: func(m *MockVertexGenerativeClient) {
				m.On("GenerateContent", mock.Anything, "gemini-2.5-flash", mock.Anything, mock.MatchedBy(func(c repository.GenerationConfig) bool {
					// We could inspect system instructions too but simplified here
					return c.JsonResponse == true
				})).Return(`["容量", "色"]`, nil)
			},
			want: []string{"容量", "色"},
		},
		{
			name:        "Vertex Error",
			title:       "iPhone",
			description: "desc",
			condition:   "New",
			language:    "en",
			mockSetup: func(m *MockVertexGenerativeClient) {
				m.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", assert.AnError)
			},
			want: []string{},
		},
		{
			name:        "Invalid JSON Response",
			title:       "iPhone",
			description: "desc",
			condition:   "New",
			language:    "en",
			mockSetup: func(m *MockVertexGenerativeClient) {
				m.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(`invalid json`, nil)
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := new(MockVertexGenerativeClient)
			tt.mockSetup(client)

			s := NewSuggestionService(client)
			got := s.GetListingSuggestion(context.Background(), tt.title, tt.description, tt.condition, tt.language)

			assert.Equal(t, tt.want, got)
			client.AssertExpectations(t)
		})
	}
}
