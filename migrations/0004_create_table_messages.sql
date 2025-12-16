ALTER TABLE users
    ADD avatarUrl TEXT NULL;

CREATE TABLE messages
(
    id          CHAR(30)     NOT NULL PRIMARY KEY,
    sender_id   VARCHAR(128) NULL,
    receiver_id VARCHAR(128) NULL,
    content     TEXT         NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_messages_id CHECK (id LIKE 'msg_%'),
    CONSTRAINT fk_messages_sender FOREIGN KEY (sender_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_messages_receiver FOREIGN KEY (receiver_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE SET NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
