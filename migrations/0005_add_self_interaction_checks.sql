-- Add validtion to ensure sender/receiver and buyer/seller are different
-- Dialect: MySQL (InnoDB, utf8mb4)

-- Update messages table
ALTER TABLE messages
    DROP FOREIGN KEY fk_messages_sender,
    DROP FOREIGN KEY fk_messages_receiver;

ALTER TABLE messages
    ADD CONSTRAINT fk_messages_sender FOREIGN KEY (sender_id) REFERENCES users (id)
        ON UPDATE RESTRICT ON DELETE RESTRICT,
    ADD CONSTRAINT fk_messages_receiver FOREIGN KEY (receiver_id) REFERENCES users (id)
        ON UPDATE RESTRICT ON DELETE RESTRICT,
    ADD CONSTRAINT chk_messages_self CHECK (sender_id <> receiver_id);

-- Update orders table
ALTER TABLE orders
    DROP FOREIGN KEY fk_orders_buyer,
    DROP FOREIGN KEY fk_orders_seller;

ALTER TABLE orders
    ADD CONSTRAINT fk_orders_buyer FOREIGN KEY (buyer_id) REFERENCES users (id)
        ON UPDATE RESTRICT ON DELETE RESTRICT,
    ADD CONSTRAINT fk_orders_seller FOREIGN KEY (seller_id) REFERENCES users (id)
        ON UPDATE RESTRICT ON DELETE RESTRICT,
    ADD CONSTRAINT chk_orders_self CHECK (buyer_id <> seller_id);
