package service

import (
	"context"
	"testing"
	"time"

	"uttc-hackathon-backend/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) CreateMessage(ctx context.Context, msg *models.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MockMessageRepository) GetMessages(ctx context.Context, userID, otherUserID string) ([]*models.Message, error) {
	args := m.Called(ctx, userID, otherUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Message), args.Error(1)
}

func (m *MockMessageRepository) GetLatestIncomingMessages(ctx context.Context, userID string) ([]*models.Message, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Message), args.Error(1)
}

func (m *MockMessageRepository) GetLatestOutgoingMessages(ctx context.Context, userID string) ([]*models.Message, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Message), args.Error(1)
}

type MockMessageUserRepo struct {
	mock.Mock
}

func (m *MockMessageUserRepo) GetUserProfile(ctx context.Context, id string) (*models.UserProfile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserProfile), args.Error(1)
}

func TestMessageService_CreateMessage(t *testing.T) {
	tests := []struct {
		name       string
		senderID   string
		receiverID string
		content    string
		mockSetup  func(*MockMessageRepository)
		wantErr    bool
		errType    error
	}{
		{
			name:       "Success",
			senderID:   "u1",
			receiverID: "u2",
			content:    "Hello",
			mockSetup: func(m *MockMessageRepository) {
				m.On("CreateMessage", mock.Anything, mock.MatchedBy(func(msg *models.Message) bool {
					return msg.SenderID == "u1" && msg.ReceiverID == "u2" && msg.Content == "Hello" && msg.ID != ""
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "Self Message",
			senderID:   "u1",
			receiverID: "u1",
			content:    "Hello",
			mockSetup:  func(m *MockMessageRepository) {},
			wantErr:    true,
			errType:    ErrSelfMessage,
		},
		{
			name:       "Empty Content",
			senderID:   "u1",
			receiverID: "u2",
			content:    "",
			mockSetup:  func(m *MockMessageRepository) {},
			wantErr:    true,
			errType:    ErrContentRequired,
		},
		{
			name:       "DB Error",
			senderID:   "u1",
			receiverID: "u2",
			content:    "Hello",
			mockSetup: func(m *MockMessageRepository) {
				m.On("CreateMessage", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			wantErr: true,
			errType: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockMessageRepository)
			userRepo := new(MockMessageUserRepo)
			tt.mockSetup(repo)

			s := NewMessageService(repo, userRepo)
			got, err := s.CreateMessage(context.Background(), tt.senderID, tt.receiverID, tt.content)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestMessageService_GetConversations(t *testing.T) {
	// Setup timestamps
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	msgIncoming := &models.Message{ID: "m1", SenderID: "u2", ReceiverID: "u1", CreatedAt: yesterday, Content: "Hi"}
	// Outgoing message to u2 is newer
	msgOutgoing := &models.Message{ID: "m2", SenderID: "u1", ReceiverID: "u2", CreatedAt: now, Content: "Hello back"}

	// Another conversation with u3 (only incoming)
	msgIncoming3 := &models.Message{ID: "m3", SenderID: "u3", ReceiverID: "u1", CreatedAt: yesterday, Content: "Yo"}

	user2 := &models.UserProfile{ID: "u2", Name: "User 2"}
	user3 := &models.UserProfile{ID: "u3", Name: "User 3"}

	tests := []struct {
		name      string
		userID    string
		mockSetup func(*MockMessageRepository, *MockMessageUserRepo)
		wantLen   int
		wantErr   bool
	}{
		{
			name:   "Success - Merging and Sorting",
			userID: "u1",
			mockSetup: func(m *MockMessageRepository, u *MockMessageUserRepo) {
				m.On("GetLatestIncomingMessages", mock.Anything, "u1").Return([]*models.Message{msgIncoming, msgIncoming3}, nil)
				m.On("GetLatestOutgoingMessages", mock.Anything, "u1").Return([]*models.Message{msgOutgoing}, nil)

				// Expect profile fetch for u2 (partner in first conv)
				u.On("GetUserProfile", mock.Anything, "u2").Return(user2, nil)
				// Expect profile fetch for u3 (partner in second conv)
				u.On("GetUserProfile", mock.Anything, "u3").Return(user3, nil)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name:   "Repo Error",
			userID: "u1",
			mockSetup: func(m *MockMessageRepository, u *MockMessageUserRepo) {
				m.On("GetLatestIncomingMessages", mock.Anything, "u1").Return(nil, assert.AnError)
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockMessageRepository)
			userRepo := new(MockMessageUserRepo)
			tt.mockSetup(repo, userRepo)

			s := NewMessageService(repo, userRepo)
			got, err := s.GetConversations(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, got, tt.wantLen)
				// Verify sorting: msgOutgoing (now) should be before msgIncoming3 (yesterday)
				// Wait, msgOutgoing and msgIncoming are same conversation (u2).
				// We merged them, keeping newer (msgOutgoing).
				// The other conversation is u3 (msgIncoming3 at yesterday).
				// So u2 conversation is newer than u3 conversation.
				// got[0] should be u2

				assert.Equal(t, "u2", got[0].User.ID)
				assert.Equal(t, "u3", got[1].User.ID)

				// Check message content
				assert.Equal(t, "Hello back", got[0].Message.Content) // merged result
			}
			repo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}
