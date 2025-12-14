package app

import (
	"database/sql"
	"net/http"
	"uttc-hackathon-backend/internal/handler"
	"uttc-hackathon-backend/internal/repository"
	"uttc-hackathon-backend/internal/service"
)

type App struct {
	UserHandler *handler.UserHandler
}

func NewApp(db *sql.DB) *App {
	userRepo := repository.NewUserRepo(db)

	userSvc := service.NewUserService(userRepo)

	return &App{
		UserHandler: handler.NewUserHandler(userSvc),
	}
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", a.UserHandler.HandleGet)

	return mux
}
