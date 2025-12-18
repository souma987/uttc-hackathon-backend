package client

import (
	"context"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/transport/http"
	"google.golang.org/genai"
)

func InitVertexAI(projectID, location, googleCredentialsJson string) *genai.Client {
	ctx := context.Background()

	jsonCreds := []byte(googleCredentialsJson)

	// Create an HTTP client with the credentials
	httpClient, _, err := http.NewClient(ctx, option.WithCredentialsJSON(jsonCreds), option.WithScopes("https://www.googleapis.com/auth/cloud-platform"))
	if err != nil {
		log.Fatalf("error creating http client for vertex ai: %v\n", err)
	}

	cfg := &genai.ClientConfig{
		Project:    projectID,
		Location:   location,
		Backend:    genai.BackendVertexAI,
		HTTPClient: httpClient,
	}

	client, err := genai.NewClient(ctx, cfg)
	if err != nil {
		log.Fatalf("error initializing vertex ai client: %v\n", err)
	}

	log.Println("Initialized Vertex AI Client")
	return client
}
