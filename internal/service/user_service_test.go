package service

import (
	"context"
	"testing"

	"uttc-hackathon-backend/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock mocks
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUser(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserProfile(ctx context.Context, id string) (*models.UserProfile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserProfile), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockFirebaseRepository struct {
	mock.Mock
}

func (m *MockFirebaseRepository) CreateUser(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func (m *MockFirebaseRepository) DeleteUser(ctx context.Context, uid string) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *MockFirebaseRepository) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	args := m.Called(ctx, idToken)
	return args.String(0), args.Error(1)
}

func TestUserService_SignUp(t *testing.T) {
	tests := []struct {
		name          string
		inputName     string
		inputEmail    string
		inputPassword string
		inputAvatar   string
		mockSetup     func(*MockUserRepository, *MockFirebaseRepository)
		wantErr       bool
		errType       error
	}{
		{
			name:          "Success",
			inputName:     "Test User",
			inputEmail:    "test@example.com",
			inputPassword: "password123",
			inputAvatar:   "http://example.com/avatar.jpg",
			mockSetup: func(userRepo *MockUserRepository, fbRepo *MockFirebaseRepository) {
				fbRepo.On("CreateUser", mock.Anything, "test@example.com", "password123").Return("uid123", nil)
				userRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
					return u.ID == "uid123" && u.Name == "Test User" && u.Email == "test@example.com"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:          "Password too short",
			inputName:     "Test User",
			inputEmail:    "test@example.com",
			inputPassword: "short",
			inputAvatar:   "",
			mockSetup:     func(userRepo *MockUserRepository, fbRepo *MockFirebaseRepository) {},
			wantErr:       true,
			errType:       ErrInvalidPasswordLength,
		},
		{
			name:          "Firebase error",
			inputName:     "Test User",
			inputEmail:    "test@example.com",
			inputPassword: "password123",
			inputAvatar:   "",
			mockSetup: func(userRepo *MockUserRepository, fbRepo *MockFirebaseRepository) {
				fbRepo.On("CreateUser", mock.Anything, "test@example.com", "password123").Return("", assert.AnError)
			},
			wantErr: true,
			errType: assert.AnError,
		},
		{
			name:          "DB error - cleanup firebase user",
			inputName:     "Test User",
			inputEmail:    "test@example.com",
			inputPassword: "password123",
			inputAvatar:   "",
			mockSetup: func(userRepo *MockUserRepository, fbRepo *MockFirebaseRepository) {
				fbRepo.On("CreateUser", mock.Anything, "test@example.com", "password123").Return("uid123", nil)
				userRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
					return u.ID == "uid123"
				})).Return(assert.AnError)
				fbRepo.On("DeleteUser", mock.Anything, "uid123").Return(nil)
			},
			wantErr: true,
			errType: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			fbRepo := new(MockFirebaseRepository)
			tt.mockSetup(userRepo, fbRepo)

			s := NewUserService(userRepo, fbRepo)
			got, err := s.SignUp(context.Background(), tt.inputName, tt.inputEmail, tt.inputPassword, tt.inputAvatar)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.inputName, got.Name)
			}
			userRepo.AssertExpectations(t)
			fbRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserProfile(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		mockSetup func(*MockUserRepository)
		want      *models.UserProfile
		wantErr   bool
		errType   error
	}{
		{
			name:   "Success",
			userID: "uid123",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserProfile", mock.Anything, "uid123").Return(&models.UserProfile{ID: "uid123", Name: "Test"}, nil)
			},
			want:    &models.UserProfile{ID: "uid123", Name: "Test"},
			wantErr: false,
		},
		{
			name:   "Not found",
			userID: "uid123",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserProfile", mock.Anything, "uid123").Return(nil, nil)
			},
			want:    nil,
			wantErr: true,
			errType: ErrUserNotFound,
		},
		{
			name:   "DB Error",
			userID: "uid123",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserProfile", mock.Anything, "uid123").Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: true,
			errType: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			fbRepo := new(MockFirebaseRepository) // Not used in this test
			tt.mockSetup(userRepo)

			s := NewUserService(userRepo, fbRepo)
			got, err := s.GetUserProfile(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			userRepo.AssertExpectations(t)
		})
	}
}
