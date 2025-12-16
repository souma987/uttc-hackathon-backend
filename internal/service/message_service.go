package service

import (
	"context"
	"errors"
	"time"
	"uttc-hackathon-backend/internal/models"
	"uttc-hackathon-backend/internal/repository"

	"github.com/oklog/ulid/v2"
)

var (
	ErrContentRequired = errors.New("content is required")
)

type MessageService struct {
	repo *repository.MessageRepository
}

func NewMessageService(repo *repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) CreateMessage(ctx context.Context, senderID, receiverID, content string) (*models.Message, error) {
	if content == "" {
		return nil, ErrContentRequired
	}

	id := "msg_" + ulid.Make().String()
	now := time.Now()

	msg := &models.Message{
		ID:         id,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  now,
	}

	if err := s.repo.CreateMessage(ctx, msg); err != nil {
		return nil, err
	}

	return msg, nil
}
