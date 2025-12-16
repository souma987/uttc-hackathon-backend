package app

import (
	"database/sql"
	"net/http"
	"uttc-hackathon-backend/internal/handler"
	"uttc-hackathon-backend/internal/middleware"
	"uttc-hackathon-backend/internal/repository"
	"uttc-hackathon-backend/internal/service"

	"firebase.google.com/go/v4/auth"
)

type App struct {
	UserHandler    *handler.UserHandler
	listingHandler *handler.ListingHandler
	orderHandler   *handler.OrderHandler
	MessageHandler *handler.MessageHandler
	authMiddleware func(http.Handler) http.Handler
}

func NewApp(db *sql.DB, fbAuth *auth.Client) *App {
	userRepo := repository.NewUserRepo(db)
	listingRepo := repository.NewListingRepo(db)
	orderRepo := repository.NewOrderRepo(db)
	messageRepo := repository.NewMessageRepository(db)
	fbRepo := repository.NewFirebaseAuthRepo(fbAuth)

	userSvc := service.NewUserService(userRepo, fbRepo)
	listingSvc := service.NewListingService(listingRepo)
	orderSvc := service.NewOrderService(orderRepo)
	messageSvc := service.NewMessageService(messageRepo)

	userHandler := handler.NewUserHandler(userSvc)
	listingHandler := handler.NewListingHandler(listingSvc, userSvc)
	orderHandler := handler.NewOrderHandler(orderSvc, userSvc)
	messageHandler := handler.NewMessageHandler(messageSvc, userSvc)

	authMW := middleware.AuthMiddleware(userSvc)

	return &App{
		UserHandler:    userHandler,
		listingHandler: listingHandler,
		orderHandler:   orderHandler,
		MessageHandler: messageHandler,
		authMiddleware: authMW,
	}
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()

	// Users
	mux.HandleFunc("POST /users", a.UserHandler.HandleCreate)
	mux.Handle("GET /me", a.authMiddleware(http.HandlerFunc(a.UserHandler.HandleMe)))

	// Listings
	mux.HandleFunc("GET /listings/feed", a.listingHandler.HandleFeed)
	mux.Handle("POST /listings", a.authMiddleware(http.HandlerFunc(a.listingHandler.HandleCreate)))
	mux.HandleFunc("GET /listings/{id}", a.listingHandler.HandleGetListing)

	// Orders
	mux.Handle("POST /orders", a.authMiddleware(http.HandlerFunc(a.orderHandler.HandleCreate)))
	mux.Handle("GET /orders/{orderId}", a.authMiddleware(http.HandlerFunc(a.orderHandler.HandleGet)))

	// Messages
	mux.Handle("POST /messages", a.authMiddleware(http.HandlerFunc(a.MessageHandler.HandleCreate)))
	mux.Handle("GET /messages/with/{userid}", a.authMiddleware(http.HandlerFunc(a.MessageHandler.HandleGetMessages)))

	return mux
}
