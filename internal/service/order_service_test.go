package service

import (
	"context"
	"testing"

	"uttc-hackathon-backend/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, listingID string, fn func(*models.Listing) (*models.Order, error)) error {
	args := m.Called(ctx, listingID, fn)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	return args.Error(0)
}

func (m *MockOrderRepository) GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetOrdersByUserID(ctx context.Context, userID string) ([]*models.Order, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Order), args.Error(1)
}

func TestOrderService_CreateOrder(t *testing.T) {
	tests := []struct {
		name      string
		buyerID   string
		req       *models.Order
		listing   *models.Listing // The listing currently in DB
		mockSetup func(*MockOrderRepository, *models.Order, *models.Listing)
		wantErr   bool
		errType   error
	}{
		{
			name:    "Success",
			buyerID: "buyer1",
			req:     &models.Order{ListingID: "lst1", Quantity: 1},
			listing: &models.Listing{
				ID: "lst1", SellerID: "seller1", Status: models.ListingStatusActive, Quantity: 10, Price: 1000,
				Images: []models.ListingImage{{URL: "img.jpg"}},
			},
			mockSetup: func(m *MockOrderRepository, req *models.Order, l *models.Listing) {
				m.On("CreateOrder", mock.Anything, "lst1", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					fn := args.Get(2).(func(*models.Listing) (*models.Order, error))
					_, err := fn(l) // Execute the callback with our test listing
					assert.NoError(t, err)
				})
			},
			wantErr: false,
		},
		{
			name:    "Buy Own Listing",
			buyerID: "seller1",
			req:     &models.Order{ListingID: "lst1", Quantity: 1},
			listing: &models.Listing{ID: "lst1", SellerID: "seller1", Status: models.ListingStatusActive, Quantity: 10},
			mockSetup: func(m *MockOrderRepository, req *models.Order, l *models.Listing) {
				m.On("CreateOrder", mock.Anything, "lst1", mock.Anything).Return(ErrBuyOwnListing).Run(func(args mock.Arguments) {
					fn := args.Get(2).(func(*models.Listing) (*models.Order, error))
					_, err := fn(l)
					assert.Equal(t, ErrBuyOwnListing, err)
				})
			},
			wantErr: true,
			errType: ErrBuyOwnListing,
		},
		{
			name:    "Listing Not Active",
			buyerID: "buyer1",
			req:     &models.Order{ListingID: "lst1", Quantity: 1},
			listing: &models.Listing{ID: "lst1", SellerID: "seller1", Status: models.ListingStatusSold, Quantity: 10},
			mockSetup: func(m *MockOrderRepository, req *models.Order, l *models.Listing) {
				m.On("CreateOrder", mock.Anything, "lst1", mock.Anything).Return(ErrListingNotActive).Run(func(args mock.Arguments) {
					fn := args.Get(2).(func(*models.Listing) (*models.Order, error))
					_, err := fn(l)
					assert.Equal(t, ErrListingNotActive, err)
				})
			},
			wantErr: true,
			errType: ErrListingNotActive,
		},
		{
			name:    "Insufficient Stock",
			buyerID: "buyer1",
			req:     &models.Order{ListingID: "lst1", Quantity: 11},
			listing: &models.Listing{ID: "lst1", SellerID: "seller1", Status: models.ListingStatusActive, Quantity: 10},
			mockSetup: func(m *MockOrderRepository, req *models.Order, l *models.Listing) {
				m.On("CreateOrder", mock.Anything, "lst1", mock.Anything).Return(ErrInsufficientStock).Run(func(args mock.Arguments) {
					fn := args.Get(2).(func(*models.Listing) (*models.Order, error))
					_, err := fn(l)
					assert.Equal(t, ErrInsufficientStock, err)
				})
			},
			wantErr: true,
			errType: ErrInsufficientStock,
		},
		{
			name:    "Quantity Zero",
			buyerID: "buyer1",
			req:     &models.Order{ListingID: "lst1", Quantity: 0},
			listing: nil,
			mockSetup: func(m *MockOrderRepository, req *models.Order, l *models.Listing) {
				// Repo not called
			},
			wantErr: true,
			errType: ErrQuantityInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockOrderRepository)
			if tt.mockSetup != nil {
				tt.mockSetup(repo, tt.req, tt.listing)
			}

			s := NewOrderService(repo)
			got, err := s.CreateOrder(context.Background(), tt.buyerID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				// Check calculated fields
				assert.Equal(t, tt.listing.Price*tt.req.Quantity, got.TotalPrice)
				assert.NotEmpty(t, got.ID)
				assert.Equal(t, models.OrderStatusPaid, got.Status)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestOrderService_GetOrder(t *testing.T) {
	order := &models.Order{ID: "ord1", BuyerID: "buyer1", SellerID: "seller1"}

	tests := []struct {
		name      string
		userID    string
		orderID   string
		mockSetup func(*MockOrderRepository)
		want      *models.Order
		wantErr   bool
		errType   error
	}{
		{
			name:    "Success - Buyer",
			userID:  "buyer1",
			orderID: "ord1",
			mockSetup: func(m *MockOrderRepository) {
				m.On("GetOrder", mock.Anything, "ord1").Return(order, nil)
			},
			want:    order,
			wantErr: false,
		},
		{
			name:    "Success - Seller",
			userID:  "seller1",
			orderID: "ord1",
			mockSetup: func(m *MockOrderRepository) {
				m.On("GetOrder", mock.Anything, "ord1").Return(order, nil)
			},
			want:    order,
			wantErr: false,
		},
		{
			name:    "Unauthorized",
			userID:  "other",
			orderID: "ord1",
			mockSetup: func(m *MockOrderRepository) {
				m.On("GetOrder", mock.Anything, "ord1").Return(order, nil)
			},
			want:    nil,
			wantErr: true,
			errType: ErrUnauthorized,
		},
		{
			name:    "Not Found",
			userID:  "buyer1",
			orderID: "ord1",
			mockSetup: func(m *MockOrderRepository) {
				m.On("GetOrder", mock.Anything, "ord1").Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: true,
			errType: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockOrderRepository)
			tt.mockSetup(repo)

			s := NewOrderService(repo)
			got, err := s.GetOrder(context.Background(), tt.userID, tt.orderID)

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

func TestOrderService_GetOrdersByUser(t *testing.T) {
	orders := []*models.Order{{ID: "ord1"}, {ID: "ord2"}}
	repo := new(MockOrderRepository)
	repo.On("GetOrdersByUserID", mock.Anything, "user1").Return(orders, nil)

	s := NewOrderService(repo)
	got, err := s.GetOrdersByUser(context.Background(), "user1")

	assert.NoError(t, err)
	assert.Equal(t, orders, got)
	repo.AssertExpectations(t)
}
