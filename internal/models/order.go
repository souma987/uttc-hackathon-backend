package models

import "time"

type OrderStatus string

const (
	OrderStatusAwaitingPayment OrderStatus = "awaiting_payment"
	OrderStatusPaid            OrderStatus = "paid"
	OrderStatusShipped         OrderStatus = "shipped"
	OrderStatusDelivered       OrderStatus = "delivered"
	OrderStatusCompleted       OrderStatus = "completed"
	OrderStatusCancelled       OrderStatus = "cancelled"
	OrderStatusDisputed        OrderStatus = "disputed"
)

type Order struct {
	ID               string      `json:"id"`
	BuyerID          string      `json:"buyer_id"`
	SellerID         string      `json:"seller_id"`
	ListingID        string      `json:"listing_id"`
	ListingTitle     string      `json:"listing_title"`
	ListingMainImage string      `json:"listing_main_image"`
	ListingPrice     int         `json:"listing_price"`
	Quantity         int         `json:"quantity"`
	PriceTotal       int         `json:"price_total"`
	PlatformFee      int         `json:"platform_fee"`
	TotalCharged     int         `json:"total_charged"`
	Status           OrderStatus `json:"status"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}
