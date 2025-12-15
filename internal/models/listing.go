package models

import "time"

type ListingStatus string
type ItemCondition string

const (
	ListingStatusDraft  ListingStatus = "draft"
	ListingStatusActive ListingStatus = "active"
	ListingStatusSold   ListingStatus = "sold"

	ItemConditionNew       ItemCondition = "new"
	ItemConditionExcellent ItemCondition = "excellent"
	ItemConditionGood      ItemCondition = "good"
	ItemConditionNotGood   ItemCondition = "not_good"
	ItemConditionBad       ItemCondition = "bad"
)

type ListingImage struct {
	URL string `json:"url"`
}

type Listing struct {
	ID            string         `json:"id"`
	SellerID      string         `json:"seller_id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Images        []ListingImage `json:"images"`
	Price         int            `json:"price"`
	Quantity      int            `json:"quantity"`
	Status        ListingStatus  `json:"status"`
	ItemCondition ItemCondition  `json:"item_condition"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}
