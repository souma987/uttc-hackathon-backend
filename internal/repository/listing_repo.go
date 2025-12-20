package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"uttc-hackathon-backend/internal/models"
)

type ListingRepo struct {
	db *sql.DB
}

func NewListingRepo(db *sql.DB) *ListingRepo {
	return &ListingRepo{db: db}
}

func (r *ListingRepo) GetListingsFeed(ctx context.Context, limit, offset int) ([]*models.Listing, error) {
	query := `
		SELECT id, seller_id, title, description, images, price, quantity, status, item_condition, created_at, updated_at
		FROM listings
		WHERE status = 'active'
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query listings feed: %w", err)
	}
	defer rows.Close()

	var listings []*models.Listing
	for rows.Next() {
		var l models.Listing
		var imagesJSON []byte

		err := rows.Scan(
			&l.ID,
			&l.SellerID,
			&l.Title,
			&l.Description,
			&imagesJSON,
			&l.Price,
			&l.Quantity,
			&l.Status,
			&l.ItemCondition,
			&l.CreatedAt,
			&l.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan listing: %w", err)
		}

		if err := json.Unmarshal(imagesJSON, &l.Images); err != nil {
			return nil, fmt.Errorf("unmarshal listing images: %w", err)
		}

		listings = append(listings, &l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listings rows: %w", err)
	}

	return listings, nil
}

func (r *ListingRepo) CreateListing(ctx context.Context, l *models.Listing) error {
	imagesJSON, err := json.Marshal(l.Images)
	if err != nil {
		return fmt.Errorf("marshal listing images: %w", err)
	}

	query := `
		INSERT INTO listings (id, seller_id, title, description, images, price, quantity, status, item_condition)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, query,
		l.ID,
		l.SellerID,
		l.Title,
		l.Description,
		imagesJSON,
		l.Price,
		l.Quantity,
		l.Status,
		l.ItemCondition,
		l.CreatedAt,
		l.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert listing: %w", err)
	}

	return nil
}

func (r *ListingRepo) GetListing(ctx context.Context, id string) (*models.Listing, error) {
	query := `
		SELECT id, seller_id, title, description, images, price, quantity, status, item_condition, created_at, updated_at
		FROM listings
		WHERE id = ?
	`

	var l models.Listing
	var imagesJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&l.ID,
		&l.SellerID,
		&l.Title,
		&l.Description,
		&imagesJSON,
		&l.Price,
		&l.Quantity,
		&l.Status,
		&l.ItemCondition,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Or nil, sql.ErrNoRows - service will handle
		}
		return nil, fmt.Errorf("get listing: %w", err)
	}

	if err := json.Unmarshal(imagesJSON, &l.Images); err != nil {
		return nil, fmt.Errorf("unmarshal listing images: %w", err)
	}

	return &l, nil
}
