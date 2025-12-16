package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"uttc-hackathon-backend/internal/models"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

var (
	ErrListingNotFound   = errors.New("listing not found or not active")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrListingNotActive  = errors.New("listing is not active")
)

func (r *OrderRepo) CreateOrder(ctx context.Context, listingID string, quantity int, fn func(*models.Listing) (*models.Order, error)) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// 1. Lock and Get Listing Details
	// We select FOR UPDATE to prevent race conditions on quantity
	queryGet := `
		SELECT id, seller_id, title, images, price, quantity, status
		FROM listings
		WHERE id = ?
		FOR UPDATE
	`
	var l models.Listing
	var imagesJSON []byte

	err = tx.QueryRowContext(ctx, queryGet, listingID).Scan(
		&l.ID, &l.SellerID, &l.Title, &imagesJSON, &l.Price, &l.Quantity, &l.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrListingNotFound
		}
		return fmt.Errorf("get listing for update: %w", err)
	}

	if l.Status != models.ListingStatusActive {
		return ErrListingNotActive
	}

	if l.Quantity < quantity {
		return ErrInsufficientStock
	}

	// Parse images
	if err := json.Unmarshal(imagesJSON, &l.Images); err != nil {
		return fmt.Errorf("unmarshal images: %w", err)
	}

	// 2. Call Service Callback
	o, err := fn(&l)
	if err != nil {
		return err // Service error, rollback happens via defer
	}

	// 3. Update Listing Quantity
	queryUpdate := `
		UPDATE listings
		SET quantity = quantity - ?
		WHERE id = ?
	`
	_, err = tx.ExecContext(ctx, queryUpdate, quantity, listingID)
	if err != nil {
		return fmt.Errorf("update listing quantity: %w", err)
	}

	// 4. Insert Order
	queryInsert := `
		INSERT INTO orders (
			id, buyer_id, seller_id, listing_id, listing_title, listing_main_image,
			listing_price, quantity, total_price, platform_fee, net_payout, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, queryInsert,
		o.ID, o.BuyerID, o.SellerID, o.ListingID, o.ListingTitle, o.ListingMainImage,
		o.ListingPrice, o.Quantity, o.TotalPrice, o.PlatformFee, o.NetPayout, o.Status,
		o.CreatedAt, o.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
