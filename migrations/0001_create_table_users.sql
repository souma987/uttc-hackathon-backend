CREATE TABLE users (
    id VARCHAR(128) PRIMARY KEY, -- Matches Firebase UID
    email VARCHAR(255) NOT NULL,
    username VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);