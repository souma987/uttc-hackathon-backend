package service

import (
	"context"
	"errors"

	"uttc-hackathon-backend/internal/models"
)

var ErrInvalidPasswordLength = errors.New("password must be between 8 and 4096 characters")

type UserService struct {
	repo         UserRepository
	firebaseAuth FirebaseRepository
}

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

type FirebaseRepository interface {
	CreateUser(ctx context.Context, email, password string) (string, error)
	DeleteUser(ctx context.Context, uid string) error
	VerifyIDToken(ctx context.Context, idToken string) (string, error)
}

func NewUserService(repo UserRepository, fb FirebaseRepository) *UserService {
	return &UserService{repo: repo, firebaseAuth: fb}
}

// SignUp creates a user in Firebase and then inserts a matching user in DB.
func (s *UserService) SignUp(ctx context.Context, name, email, password, avatarURL string) (*models.User, error) {
	if len(password) < 8 || len(password) > 4096 {
		return nil, ErrInvalidPasswordLength
	}

	uid, err := s.firebaseAuth.CreateUser(ctx, email, password)
	if err != nil {
		return nil, err
	}

	u := &models.User{ID: uid, Name: name, Email: email, AvatarURL: avatarURL}
	if err := s.repo.CreateUser(ctx, u); err != nil {
		// Rollback Firebase user on DB failure
		_ = s.firebaseAuth.DeleteUser(ctx, uid)
		return nil, err
	}
	return u, nil
}

// GetCurrentUser validates a Firebase ID token and returns the corresponding user from DB.
func (s *UserService) GetCurrentUser(ctx context.Context, idToken string) (*models.User, error) {
	uid, err := s.firebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUser(ctx, uid)
}
