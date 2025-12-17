package database

import (
	"context"
	"log"

	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/option"
)

func InitVertexAI(projectID, location, googleCredentialsJson string) *genai.Client {
	ctx := context.Background()

	jsonCreds := []byte(googleCredentialsJson)
	client, err := genai.NewClient(ctx, projectID, location, option.WithCredentialsJSON(jsonCreds))
	if err != nil {
		log.Fatalf("error initializing vertex ai client: %v\n", err)
	}

	log.Println("Initialized Vertex AI Client")
	return client
}
