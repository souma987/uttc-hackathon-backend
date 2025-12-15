-- Marketplace core schema
-- Dialect: MySQL (InnoDB, utf8mb4)

ALTER TABLE users
    CONVERT TO CHARACTER SET utf8mb4
        COLLATE utf8mb4_unicode_ci;

-- Core listings table
CREATE TABLE listings
(
    id             CHAR(30)                                             NOT NULL PRIMARY KEY,
    seller_id      VARCHAR(128)                                         NOT NULL, -- FK to users.id (Firebase UID)
    title          VARCHAR(200)                                         NOT NULL,
    description    TEXT                                                 NOT NULL,
    images         JSON                                                 NOT NULL,
    price          INT UNSIGNED                                         NOT NULL,
    quantity       INT UNSIGNED                                         NOT NULL,
    status         ENUM ('draft','active','sold')                       NOT NULL,
    item_condition ENUM ('new', 'excellent', 'good', 'not_good', 'bad') NOT NULL,
    created_at     TIMESTAMP                                            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP                                            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT chk_images_schema CHECK (
        JSON_SCHEMA_VALID(
                '{
                  "type": "array",
                  "maxItems": 10,
                  "minItems": 1,
                  "items": {
                    "type": "object",
                    "required": [
                      "url"
                    ],
                    "properties": {
                      "url": {
                        "type": "string"
                      }
                    },
                    "additionalProperties": false
                  }
                }',
                images
        )
        ),
    CONSTRAINT chk_listings_id CHECK (id LIKE 'lst_%'),
    CONSTRAINT fk_listings_seller FOREIGN KEY (seller_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    INDEX idx_listings_seller (seller_id),
    INDEX idx_listings_status (status),
    FULLTEXT INDEX idx_listings_search (title, description) WITH PARSER ngram
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE categories
(
    id        CHAR(30)     NOT NULL PRIMARY KEY,
    parent_id CHAR(30)     NULL,
    path      VARCHAR(512) NOT NULL, -- slash delimited
    CONSTRAINT chk_categories_id CHECK (id LIKE 'cat_%'),
    CONSTRAINT fk_categories_parent FOREIGN KEY (parent_id) REFERENCES categories (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    INDEX idx_categories_path (path)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- Many-to-many between listings and categories
CREATE TABLE listing_categories
(
    listing_id  CHAR(30)    NOT NULL,
    category_id VARCHAR(64) NOT NULL,

    PRIMARY KEY (listing_id, category_id),
    CONSTRAINT fk_listcat_listing FOREIGN KEY (listing_id) REFERENCES listings (id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_listcat_category FOREIGN KEY (category_id) REFERENCES categories (id)
        ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE orders
(
    id                 CHAR(30)     NOT NULL PRIMARY KEY,
    buyer_id           VARCHAR(128) NOT NULL,
    seller_id          VARCHAR(128) NOT NULL,
    listing_id         CHAR(30)     NOT NULL,
    listing_title      VARCHAR(200) NOT NULL,
    listing_main_image TEXT         NOT NULL,
    listing_price      INT UNSIGNED NOT NULL,
    quantity           INT UNSIGNED NOT NULL,
    price_total        INT UNSIGNED NOT NULL,
    platform_fee       INT UNSIGNED NOT NULL,
    total_charged      INT UNSIGNED NOT NULL,
    status             ENUM (
        'awaiting_payment',
        'paid',
        'shipped',
        'delivered',
        'completed',
        'cancelled',
        'disputed'
        )                           NOT NULL,
    created_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT chk_orders_id CHECK (id LIKE 'ord_%'),
    CONSTRAINT fk_orders_buyer FOREIGN KEY (buyer_id) REFERENCES users (id),
    CONSTRAINT fk_orders_seller FOREIGN KEY (seller_id) REFERENCES users (id),
    CONSTRAINT fk_orders_listing FOREIGN KEY (listing_id) REFERENCES listings (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
