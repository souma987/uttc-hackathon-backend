package service

import (
	"context"
	"errors"

	"uttc-hackathon-backend/internal/models"
)

var ErrInvalidPasswordLength = errors.New("password must be between 8 and 4096 characters")
var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	repo         UserRepository
	firebaseAuth FirebaseRepository
}

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*models.User, error)
	GetUserProfile(ctx context.Context, id string) (*models.UserProfile, error)
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

// VerifyToken validates a Firebase ID token and returns the UID.
func (s *UserService) VerifyToken(ctx context.Context, idToken string) (string, error) {
	return s.firebaseAuth.VerifyIDToken(ctx, idToken)
}

// GetUser returns the user from DB by ID.
func (s *UserService) GetUser(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetUser(ctx, id)
}

func (s *UserService) GetUserProfile(ctx context.Context, id string) (*models.UserProfile, error) {
	p, err := s.repo.GetUserProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrUserNotFound
	}
	return p, nil
}
