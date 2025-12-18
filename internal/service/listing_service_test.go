package service

import (
	"context"
	"testing"

	"uttc-hackathon-backend/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockListingRepository struct {
	mock.Mock
}

func (m *MockListingRepository) GetListingsFeed(ctx context.Context, limit, offset int) ([]*models.Listing, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Listing), args.Error(1)
}

func (m *MockListingRepository) CreateListing(ctx context.Context, l *models.Listing) error {
	args := m.Called(ctx, l)
	return args.Error(0)
}

func (m *MockListingRepository) GetListing(ctx context.Context, id string) (*models.Listing, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Listing), args.Error(1)
}

func TestListingService_GetFeed(t *testing.T) {
	// We want to test the parameter normalization logic in the service
	tests := []struct {
		name       string
		limit      int
		offset     int
		wantLimit  int
		wantOffset int
		mockReturn []*models.Listing
	}{
		{
			name:       "Default limit",
			limit:      0,
			offset:     0,
			wantLimit:  20,
			wantOffset: 0,
			mockReturn: []*models.Listing{},
		},
		{
			name:       "Max limit",
			limit:      1000,
			offset:     0,
			wantLimit:  100,
			wantOffset: 0,
			mockReturn: []*models.Listing{},
		},
		{
			name:       "Negative offset",
			limit:      10,
			offset:     -5,
			wantLimit:  10,
			wantOffset: 0,
			mockReturn: []*models.Listing{},
		},
		{
			name:       "Normal Case",
			limit:      15,
			offset:     5,
			wantLimit:  15,
			wantOffset: 5,
			mockReturn: []*models.Listing{{Title: "Item"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockListingRepository)
			// Expect repository to be called with normalized values
			repo.On("GetListingsFeed", mock.Anything, tt.wantLimit, tt.wantOffset).Return(tt.mockReturn, nil)

			s := NewListingService(repo)
			got, err := s.GetFeed(context.Background(), tt.limit, tt.offset)

			assert.NoError(t, err)
			assert.Equal(t, tt.mockReturn, got)
			repo.AssertExpectations(t)
		})
	}
}

func TestListingService_CreateListing(t *testing.T) {
	validImage := models.ListingImage{URL: "https://firebasestorage.googleapis.com/v0/b/bucket/o/image.jpg"}
	invalidImage := models.ListingImage{URL: "http://malicious.com/image.jpg"}

	tests := []struct {
		name      string
		req       *models.Listing
		mockSetup func(*MockListingRepository)
		wantErr   bool
		errType   error
	}{
		{
			name: "Success",
			req: &models.Listing{
				Title: "Valid Item", Price: 500, Images: []models.ListingImage{validImage},
			},
			mockSetup: func(m *MockListingRepository) {
				m.On("CreateListing", mock.Anything, mock.MatchedBy(func(l *models.Listing) bool {
					return l.ID != "" && l.Title == "Valid Item"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Title Missing",
			req: &models.Listing{
				Title: "", Price: 500, Images: []models.ListingImage{validImage},
			},
			mockSetup: func(m *MockListingRepository) {},
			wantErr:   true,
			errType:   ErrTitleRequired,
		},
		{
			name: "Price too low",
			req: &models.Listing{
				Title: "Cheap Item", Price: 50, Images: []models.ListingImage{validImage},
			},
			mockSetup: func(m *MockListingRepository) {},
			wantErr:   true,
			errType:   ErrPriceInvalid,
		},
		{
			name: "No Images",
			req: &models.Listing{
				Title: "No Image Item", Price: 500, Images: []models.ListingImage{},
			},
			mockSetup: func(m *MockListingRepository) {},
			wantErr:   true,
			errType:   ErrNoImages,
		},
		{
			name: "Invalid Image URL",
			req: &models.Listing{
				Title: "Bad Image Item", Price: 500, Images: []models.ListingImage{invalidImage},
			},
			mockSetup: func(m *MockListingRepository) {},
			wantErr:   true,
			errType:   ErrInvalidImageURL,
		},
		{
			name: "Repo Error",
			req: &models.Listing{
				Title: "Error Item", Price: 500, Images: []models.ListingImage{validImage},
			},
			mockSetup: func(m *MockListingRepository) {
				m.On("CreateListing", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			wantErr: true,
			errType: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockListingRepository)
			tt.mockSetup(repo)

			s := NewListingService(repo)
			got, err := s.CreateListing(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Contains(t, got.ID, "lst_")
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestListingService_GetListing(t *testing.T) {
	listing := &models.Listing{ID: "lst1", Title: "My Item"}

	tests := []struct {
		name      string
		id        string
		mockSetup func(*MockListingRepository)
		want      *models.Listing
		wantErr   bool
		errType   error
	}{
		{
			name: "Success",
			id:   "lst1",
			mockSetup: func(m *MockListingRepository) {
				m.On("GetListing", mock.Anything, "lst1").Return(listing, nil)
			},
			want:    listing,
			wantErr: false,
		},
		{
			name: "Not Found",
			id:   "lst1",
			mockSetup: func(m *MockListingRepository) {
				m.On("GetListing", mock.Anything, "lst1").Return(nil, nil)
			},
			want:    nil,
			wantErr: true,
			errType: ErrListingNotFound,
		},
		{
			name: "DB Error",
			id:   "lst1",
			mockSetup: func(m *MockListingRepository) {
				m.On("GetListing", mock.Anything, "lst1").Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: true,
			errType: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockListingRepository)
			tt.mockSetup(repo)

			s := NewListingService(repo)
			got, err := s.GetListing(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			repo.AssertExpectations(t)
		})
	}
}
