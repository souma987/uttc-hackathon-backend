package service

import (
	"context"
	"errors"
	"strings"
	"uttc-hackathon-backend/internal/models"

	"github.com/oklog/ulid/v2"
)

const (
	MinListingPrice       = 100
	FirebaseStoragePrefix = "https://firebasestorage.googleapis.com"
)

type ListingService struct {
	repo ListingRepository
}

type ListingRepository interface {
	GetListingsFeed(ctx context.Context, limit, offset int) ([]*models.Listing, error)
	CreateListing(ctx context.Context, l *models.Listing) error
	GetListing(ctx context.Context, id string) (*models.Listing, error)
}

func NewListingService(repo ListingRepository) *ListingService {
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
	ErrTitleRequired   = errors.New("title is required")
	ErrPriceInvalid    = errors.New("price must be at least 100")
	ErrNoImages        = errors.New("at least one image is required")
	ErrListingNotFound = errors.New("listing not found")
	ErrInvalidImageURL = errors.New("image url must start with " + FirebaseStoragePrefix)
)

func (s *ListingService) GetListing(ctx context.Context, id string) (*models.Listing, error) {
	l, err := s.repo.GetListing(ctx, id)
	if err != nil {
		return nil, err
	}
	if l == nil {
		return nil, ErrListingNotFound
	}
	return l, nil
}

func (s *ListingService) CreateListing(ctx context.Context, req *models.Listing) (*models.Listing, error) {
	if req.Title == "" {
		return nil, ErrTitleRequired
	}
	if req.Price < MinListingPrice {
		return nil, ErrPriceInvalid
	}
	if len(req.Images) == 0 {
		return nil, ErrNoImages
	}
	for _, img := range req.Images {
		if !strings.HasPrefix(img.URL, FirebaseStoragePrefix) {
			return nil, ErrInvalidImageURL
		}
	}

	req.ID = "lst_" + ulid.Make().String()

	if err := s.repo.CreateListing(ctx, req); err != nil {
		return nil, err
	}

	return req, nil
}
