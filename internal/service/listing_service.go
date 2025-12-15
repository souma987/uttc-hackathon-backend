package service

import (
	"context"
	"uttc-hackathon-backend/internal/models"
	"uttc-hackathon-backend/internal/repository"
)

type ListingService struct {
	repo *repository.ListingRepo
}

func NewListingService(repo *repository.ListingRepo) *ListingService {
	return &ListingService{repo: repo}
}

func (s *ListingService) GetFeed(ctx context.Context, limit, offset int) ([]*models.Listing, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.GetListingsFeed(ctx, limit, offset)
}
