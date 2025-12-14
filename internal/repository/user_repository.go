package repository

import (
	"context"
	"database/sql"
	"errors"
	"uttc-hackathon-backend/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetUser(ctx context.Context, id string) (*models.User, error) {
	// Simulated SQL query
	if id == "abc" {
		return &models.User{ID: "abc", Name: "Alice", Email: "alice@example.com"}, nil
	}
	return nil, errors.New("user not found")
}
