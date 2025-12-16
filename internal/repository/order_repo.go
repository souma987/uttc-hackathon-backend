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

func (r *OrderRepo) CreateOrder(ctx context.Context, o *models.Order) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// 1. Lock and Get Listing Details
	// We select FOR UPDATE to prevent race conditions on quantity
	queryGet := `
		SELECT seller_id, title, images, price, quantity, status
		FROM listings
		WHERE id = ?
		FOR UPDATE
	`
	var sellerID, title string
	var imagesJSON []byte
	var price, currentQty int
	var status string

	err = tx.QueryRowContext(ctx, queryGet, o.ListingID).Scan(
		&sellerID, &title, &imagesJSON, &price, &currentQty, &status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrListingNotFound
		}
		return fmt.Errorf("get listing for update: %w", err)
	}

	if status != string(models.ListingStatusActive) {
		return ErrListingNotActive
	}

	if currentQty < o.Quantity {
		return ErrInsufficientStock
	}

	// 2. Populate Order Details from Listing
	o.SellerID = sellerID
	o.ListingTitle = title
	o.ListingPrice = price

	// Parse images to get the main image (first one)
	var images []struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(imagesJSON, &images); err != nil {
		return fmt.Errorf("unmarshal images: %w", err)
	}
	if len(images) > 0 {
		o.ListingMainImage = images[0].URL
	}

	// Calculate totals
	o.PriceTotal = o.ListingPrice * o.Quantity
	// Platform fee logic can be placed here or service. For now, assuming simple calculation or 0 if not specified.
	// Let's assume 10% for now or keep what the service passed if any.
	// Actually, the request might not set it, so let's set it here if not set?
	// Or better, let the Service handle the business logic of Fee calculation.
	// But the Repo is responsible for filling the snapshot fields from the DB.
	// I'll trust the Service has set Quantity. I will calculate totals based on DB price.

	// Recalculate totals to ensure they match DB price
	o.PriceTotal = o.ListingPrice * o.Quantity
	// We'll leave PlatformFee and TotalCharged to be updated, or calculate them here.
	// Let's stick to the plan: Service might handle ID generation, but Repo does the transaction work.
	// I'll update the PriceTotal here to be safe.
	// I'll assume PlatformFee is calculated by Service or is 0.
	// Actually, safer to just calculate everything based on DB data.
	// Let's assume 0 fee for simplicty unless requirements say otherwise.
	// The problem says "listing_price, quantity... price_total, platform_fee, total_charged".
	// I will just compute PriceTotal.
	// I will update TotalCharged = PriceTotal + PlatformFee.
	o.TotalCharged = o.PriceTotal + o.PlatformFee

	// 3. Update Listing Quantity
	queryUpdate := `
		UPDATE listings
		SET quantity = quantity - ?
		WHERE id = ?
	`
	_, err = tx.ExecContext(ctx, queryUpdate, o.Quantity, o.ListingID)
	if err != nil {
		return fmt.Errorf("update listing quantity: %w", err)
	}

	// 4. Insert Order
	queryInsert := `
		INSERT INTO orders (
			id, buyer_id, seller_id, listing_id, listing_title, listing_main_image,
			listing_price, quantity, price_total, platform_fee, total_charged, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, queryInsert,
		o.ID, o.BuyerID, o.SellerID, o.ListingID, o.ListingTitle, o.ListingMainImage,
		o.ListingPrice, o.Quantity, o.PriceTotal, o.PlatformFee, o.TotalCharged, o.Status,
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
