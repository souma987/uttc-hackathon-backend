package database

import (
	"context"
	"log"
	"strings"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func InitFirebase(googleCredentialsJson string) *firebase.App {
	ctx := context.Background()

	cleanJson := strings.ReplaceAll(googleCredentialsJson, "\\n", "\n")
	jsonCreds := []byte(cleanJson)
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(jsonCreds))
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	log.Println("Initialized Firebase")
	return app
}
