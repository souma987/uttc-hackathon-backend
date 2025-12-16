ALTER TABLE orders
    RENAME COLUMN price_total TO total_price,
    RENAME COLUMN total_charged TO net_payout;
