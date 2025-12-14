package service

import (
	"context"
	"uttc-hackathon-backend/internal/models"
)

type UserService struct {
	repo UserRepository
}

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*models.User, error)
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) FetchUser(ctx context.Context, id string) (*models.User, error) {
	// Business logic...
	return s.repo.GetUser(ctx, id)
}
