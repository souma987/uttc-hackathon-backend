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
	ErrListingNotFound = errors.New("listing not found or not active")
	ErrOrderNotFound   = errors.New("order not found")
)

// CreateOrder updates listing and creates order atomically preventing race conditions
func (r *OrderRepo) CreateOrder(ctx context.Context, listingID string, fn func(*models.Listing) (*models.Order, error)) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// LOCK Listing
	queryGet := `
		SELECT id, seller_id, title, description, images, price, quantity, status, item_condition, created_at, updated_at
		FROM listings
		WHERE id = ?
		FOR UPDATE
	`
	var l models.Listing
	var imagesJSON []byte

	err = tx.QueryRowContext(ctx, queryGet, listingID).Scan(
		&l.ID, &l.SellerID, &l.Title, &l.Description, &imagesJSON, &l.Price, &l.Quantity, &l.Status, &l.ItemCondition, &l.CreatedAt, &l.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrListingNotFound
		}
		return fmt.Errorf("get listing for update: %w", err)
	}

	// Parse images
	if err := json.Unmarshal(imagesJSON, &l.Images); err != nil {
		return fmt.Errorf("unmarshal images: %w", err)
	}

	o, err := fn(&l)
	if err != nil {
		return err
	}

	newImagesJSON, err := json.Marshal(l.Images)
	if err != nil {
		return fmt.Errorf("marshal images: %w", err)
	}

	queryUpdate := `
		UPDATE listings
		SET title = ?, description = ?, images = ?, price = ?, quantity = ?, status = ?, item_condition = ?
		WHERE id = ?
	`
	_, err = tx.ExecContext(ctx, queryUpdate,
		l.Title, l.Description, newImagesJSON, l.Price,
		l.Quantity, l.Status, l.ItemCondition,
		listingID,
	)
	if err != nil {
		return fmt.Errorf("update listing: %w", err)
	}

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

func (r *OrderRepo) GetOrder(ctx context.Context, orderID string) (*models.Order, error) {
	query := `
		SELECT id, buyer_id, seller_id, listing_id, listing_title, listing_main_image,
			listing_price, quantity, total_price, platform_fee, net_payout, status, created_at, updated_at
		FROM orders
		WHERE id = ?
	`
	var o models.Order
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&o.ID, &o.BuyerID, &o.SellerID, &o.ListingID, &o.ListingTitle, &o.ListingMainImage,
		&o.ListingPrice, &o.Quantity, &o.TotalPrice, &o.PlatformFee, &o.NetPayout, &o.Status,
		&o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("get order: %w", err)
	}
	return &o, nil
}

func (r *OrderRepo) GetOrdersByUserID(ctx context.Context, userID string) ([]*models.Order, error) {
	query := `
		SELECT id, buyer_id, seller_id, listing_id, listing_title, listing_main_image,
			listing_price, quantity, total_price, platform_fee, net_payout, status, created_at, updated_at
		FROM orders
		WHERE buyer_id = ? OR seller_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID, &o.BuyerID, &o.SellerID, &o.ListingID, &o.ListingTitle, &o.ListingMainImage,
			&o.ListingPrice, &o.Quantity, &o.TotalPrice, &o.PlatformFee, &o.NetPayout, &o.Status,
			&o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}
		orders = append(orders, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}
	return orders, nil
}
