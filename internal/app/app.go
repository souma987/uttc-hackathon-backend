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
	UserHandler *handler.UserHandler
}

func NewApp(db *sql.DB, fbAuth *auth.Client) *App {
	userRepo := repository.NewUserRepo(db)
	fbRepo := repository.NewFirebaseAuthRepo(fbAuth)

	userSvc := service.NewUserService(userRepo, fbRepo)

	return &App{
		UserHandler: handler.NewUserHandler(userSvc),
	}
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", a.UserHandler.HandleCreate)
	mux.HandleFunc("GET /me", a.UserHandler.HandleMe)

	return mux
}
