package service

import (
	"context"
	"errors"
	"time"
	"uttc-hackathon-backend/internal/models"

	"github.com/oklog/ulid/v2"
)

var (
	ErrQuantityInvalid = errors.New("quantity must be greater than 0")
	ErrUnauthorized    = errors.New("unauthorized to access this order")
)

type OrderService struct {
	repo OrderRepository
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, listingID string, quantity int, fn func(*models.Listing) (*models.Order, error)) error
	GetOrder(ctx context.Context, orderID string) (*models.Order, error)
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, buyerID string, req *models.Order) (*models.Order, error) {
	if req.Quantity <= 0 {
		return nil, ErrQuantityInvalid
	}

	// Use repository transactional callback to ensure consistency
	err := s.repo.CreateOrder(ctx, req.ListingID, req.Quantity, func(l *models.Listing) (*models.Order, error) {
		req.ID = "ord_" + ulid.Make().String()
		req.BuyerID = buyerID
		req.SellerID = l.SellerID
		req.ListingTitle = l.Title
		req.ListingPrice = l.Price
		if len(l.Images) > 0 {
			req.ListingMainImage = l.Images[0].URL
		}
		req.Status = models.OrderStatusAwaitingPayment
		req.CreatedAt = time.Now()
		req.UpdatedAt = time.Now()

		// Calculate totals
		req.TotalPrice = l.Price * req.Quantity
		// Fee is 10% of total price, rounded up (ceil)
		req.PlatformFee = (req.TotalPrice + 9) / 10
		req.NetPayout = req.TotalPrice - req.PlatformFee

		return req, nil
	})

	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *OrderService) GetOrder(ctx context.Context, userID, orderID string) (*models.Order, error) {
	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Authorization check: User must be buyer or seller
	if order.BuyerID != userID && order.SellerID != userID {
		return nil, ErrUnauthorized
	}

	return order, nil
}
