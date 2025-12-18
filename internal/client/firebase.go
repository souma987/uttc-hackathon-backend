package client

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

func InitFirebaseAuth(googleCredentialsJson string) *auth.Client {
	ctx := context.Background()

	jsonCreds := []byte(googleCredentialsJson)
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(jsonCreds))
	if err != nil {
		log.Fatalf("error initializing firebase: %v\n", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing firebase client: %v\n", err)
	}
	log.Println("Initialized Firebase Auth")
	return client
}
