package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uttc-hackathon-backend/internal/app"
	"uttc-hackathon-backend/internal/database"
	"uttc-hackathon-backend/internal/middleware"
)

func main() {
	log.Println("Running main")

	http.ListenAndServe(":8080", nil)
	return

	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlUserPwd := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlConnectionParms := os.Getenv("MYSQL_CONNECTION_PARAMS")
	corsAllowOrigin := os.Getenv("CORS_ALLOW_ORIGIN")
	googleCredentials := os.Getenv("GOOGLE_CREDENTIALS_JSON")

	db := database.InitDB(mysqlUser, mysqlUserPwd, mysqlDatabase, mysqlHost, mysqlConnectionParms)
	defer func() {
		log.Println("Closing DB connection...")
		database.CloseDB(db)
	}()

	fb := database.InitFirebase(googleCredentials)

	routes := app.NewApp(db, fb).Routes()
	handlerWithCors := middleware.CorsMiddleware(routes, corsAllowOrigin)

	srv := &http.Server{Addr: ":8080", Handler: handlerWithCors}
	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // BLOCKING WAIT
	log.Println("Shutting down server...")

	// Graceful Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
