package app

import (
	"database/sql"
	"net/http"
	"uttc-hackathon-backend/internal/handler"
	"uttc-hackathon-backend/internal/repository"
	"uttc-hackathon-backend/internal/service"

	"firebase.google.com/go/v4/auth"
)

type App struct {
	UserHandler    *handler.UserHandler
	ListingHandler *handler.ListingHandler
}

func NewApp(db *sql.DB, fbAuth *auth.Client) *App {
	userRepo := repository.NewUserRepo(db)
	listingRepo := repository.NewListingRepo(db)
	fbRepo := repository.NewFirebaseAuthRepo(fbAuth)

	userSvc := service.NewUserService(userRepo, fbRepo)
	listingSvc := service.NewListingService(listingRepo)

	return &App{
		UserHandler:    handler.NewUserHandler(userSvc),
		ListingHandler: handler.NewListingHandler(listingSvc),
	}
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", a.UserHandler.HandleCreate)
	mux.HandleFunc("GET /me", a.UserHandler.HandleMe)

	mux.HandleFunc("GET /listings/feed", a.ListingHandler.HandleFeed)

	return mux
}
