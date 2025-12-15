package service

import (
	"context"
	"errors"
	"time"
	"uttc-hackathon-backend/internal/models"
	"uttc-hackathon-backend/internal/repository"

	"github.com/oklog/ulid/v2"
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

var (
	ErrTitleRequired = errors.New("title is required")
	ErrPriceInvalid  = errors.New("price must be greater than 0")
	ErrNoImages      = errors.New("at least one image is required")
)

func (s *ListingService) CreateListing(ctx context.Context, req *models.Listing) (*models.Listing, error) {
	// Basic Validation
	if req.Title == "" {
		return nil, ErrTitleRequired
	}
	if req.Price <= 0 {
		return nil, ErrPriceInvalid
	}
	if len(req.Images) == 0 {
		return nil, ErrNoImages
	}

	// Set defaults and system fields
	req.ID = "lst_" + ulid.Make().String()
	req.Status = models.ListingStatusActive
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	// Ensure caller (handler) sets SellerID

	if err := s.repo.CreateListing(ctx, req); err != nil {
		return nil, err
	}

	return req, nil
}
