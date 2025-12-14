package app

import (
	"database/sql"
	"net/http"
	"uttc-hackathon-backend/internal/handler"
	"uttc-hackathon-backend/internal/repository"
	"uttc-hackathon-backend/internal/service"

	firebase "firebase.google.com/go/v4"
)

type App struct {
	UserHandler *handler.UserHandler
}

func NewApp(db *sql.DB, fb *firebase.App) *App {
	userRepo := repository.NewUserRepo(db)
	fbRepo := repository.NewFirebaseAuthRepo(fb)

	userSvc := service.NewUserService(userRepo, fbRepo)

	return &App{
		UserHandler: handler.NewUserHandler(userSvc),
	}
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", a.UserHandler.HandleGet)
	mux.HandleFunc("POST /users", a.UserHandler.HandleCreate)

	return mux
}
