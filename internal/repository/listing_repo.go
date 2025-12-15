package repository

import (
	"context"
	"database/sql"
	"encoding/json"
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
