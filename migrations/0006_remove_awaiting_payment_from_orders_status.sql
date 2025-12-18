-- Remove awaiting_payment from orders status
-- Dialect: MySQL (InnoDB, utf8mb4)

ALTER TABLE orders
    MODIFY COLUMN status ENUM (
        'paid',
        'shipped',
        'delivered',
        'completed',
        'cancelled',
        'disputed'
        ) NOT NULL;
