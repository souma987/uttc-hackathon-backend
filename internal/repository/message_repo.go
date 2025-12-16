package repository

import (
	"context"
	"database/sql"
	"fmt"
	"uttc-hackathon-backend/internal/models"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) CreateMessage(ctx context.Context, m *models.Message) error {
	query := `
		INSERT INTO messages (id, sender_id, receiver_id, content, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		m.ID,
		m.SenderID,
		m.ReceiverID,
		m.Content,
		m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}
	return nil
}

func (r *MessageRepository) GetMessages(ctx context.Context, userID, otherUserID string) ([]*models.Message, error) {
	query := `
		SELECT id, sender_id, receiver_id, content, created_at
		FROM messages
		WHERE (sender_id = ? AND receiver_id = ?)
		   OR (sender_id = ? AND receiver_id = ?)
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, otherUserID, otherUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("query messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		messages = append(messages, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return messages, nil
}
