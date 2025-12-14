package database

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
)

func InitFirebase() *firebase.App {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return app
}
