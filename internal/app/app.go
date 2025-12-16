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
	listingHandler *handler.ListingHandler
	orderHandler   *handler.OrderHandler
}

func NewApp(db *sql.DB, fbAuth *auth.Client) *App {
	userRepo := repository.NewUserRepo(db)
	listingRepo := repository.NewListingRepo(db)
	orderRepo := repository.NewOrderRepo(db)
	fbRepo := repository.NewFirebaseAuthRepo(fbAuth)

	userSvc := service.NewUserService(userRepo, fbRepo)
	listingSvc := service.NewListingService(listingRepo)
	orderSvc := service.NewOrderService(orderRepo)

	userHandler := handler.NewUserHandler(userSvc)
	listingHandler := handler.NewListingHandler(listingSvc, userSvc)
	orderHandler := handler.NewOrderHandler(orderSvc, userSvc)

	return &App{
		UserHandler:    userHandler,
		listingHandler: listingHandler,
		orderHandler:   orderHandler,
	}
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", a.UserHandler.HandleCreate)
	mux.HandleFunc("GET /me", a.UserHandler.HandleMe)

	// Listings
	mux.HandleFunc("GET /listings/feed", a.listingHandler.HandleFeed)
	mux.HandleFunc("POST /listings", a.listingHandler.HandleCreate)
	mux.HandleFunc("GET /listings/{id}", a.listingHandler.HandleGetListing)

	// Orders
	mux.HandleFunc("POST /orders", a.orderHandler.HandleCreate)

	return mux
}
