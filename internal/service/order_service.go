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
)

type OrderService struct {
	repo OrderRepository
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, o *models.Order) error
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, buyerID string, req *models.Order) (*models.Order, error) {
	if req.Quantity <= 0 {
		return nil, ErrQuantityInvalid
	}

	// Generate Order ID
	req.ID = "ord_" + ulid.Make().String()

	req.BuyerID = buyerID
	req.Status = models.OrderStatusAwaitingPayment
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	// Add platform fee calculation logic if needed.
	// For now, let's keep it simple or 0.
	req.PlatformFee = 0

	if err := s.repo.CreateOrder(ctx, req); err != nil {
		return nil, err
	}

	return req, nil
}
