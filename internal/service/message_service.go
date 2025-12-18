package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"log" // Moved log to its own group as per instruction
	"uttc-hackathon-backend/internal/models"

	"sort"

	"github.com/oklog/ulid/v2"
)

var (
	ErrContentRequired = errors.New("content is required")
	ErrSelfMessage     = errors.New("cannot send message to yourself")
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, m *models.Message) error
	GetMessages(ctx context.Context, userID, otherUserID string) ([]*models.Message, error)
	GetLatestIncomingMessages(ctx context.Context, userID string) ([]*models.Message, error)
	GetLatestOutgoingMessages(ctx context.Context, userID string) ([]*models.Message, error)
}

type MessageUserRepo interface {
	GetUserProfile(ctx context.Context, id string) (*models.UserProfile, error)
}

type MessageService struct {
	repo     MessageRepository
	userRepo MessageUserRepo
}

func NewMessageService(repo MessageRepository, userRepo MessageUserRepo) *MessageService {
	return &MessageService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, senderID, receiverID, content string) (*models.Message, error) {
	if content == "" {
		return nil, ErrContentRequired
	}

	if senderID == receiverID {
		return nil, ErrSelfMessage
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

func (s *MessageService) GetMessages(ctx context.Context, userID, otherUserID string) ([]*models.Message, error) {
	return s.repo.GetMessages(ctx, userID, otherUserID)
}

func (s *MessageService) GetConversations(ctx context.Context, userID string) ([]models.Conversation, error) {
	incoming, err := s.repo.GetLatestIncomingMessages(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get incoming messages: %w", err)
	}
	outgoing, err := s.repo.GetLatestOutgoingMessages(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get outgoing messages: %w", err)
	}

	conversationMap := make(map[string]*models.Message)
	for _, m := range incoming {
		conversationMap[m.SenderID] = m
	}
	for _, m := range outgoing {
		partnerID := m.ReceiverID
		// compare timestamp and keep newer message
		if existing, ok := conversationMap[partnerID]; !ok || m.CreatedAt.After(existing.CreatedAt) {
			conversationMap[partnerID] = m
		}
	}

	// Sort by CreatedAt DESC
	var messages []*models.Message
	for _, m := range conversationMap {
		messages = append(messages, m)
	}
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.After(messages[j].CreatedAt)
	})

	// Enrich with User Profile
	var conversations []models.Conversation
	for _, m := range messages {
		partnerID := m.SenderID
		if partnerID == userID {
			partnerID = m.ReceiverID
		}

		userProfile, err := s.userRepo.GetUserProfile(ctx, partnerID)
		if err != nil {
			log.Printf("failed to get user profile for conversation: %v", err)
			return nil, fmt.Errorf("get user profile for conversation: %w", err)
		}
		if userProfile == nil {
			log.Printf("conversations: user profile not found for partnerID: %s", partnerID)
			continue
		}

		conversations = append(conversations, models.Conversation{
			Message: m,
			User:    userProfile,
		})
	}

	return conversations, nil
}
