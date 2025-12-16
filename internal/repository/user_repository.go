package repository

import (
	"context"
	"database/sql"
	"errors"
	"uttc-hackathon-backend/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetUser(ctx context.Context, id string) (*models.User, error) {
	const q = "SELECT id, username, email, avatarUrl FROM users WHERE id = ?"
	row := r.db.QueryRowContext(ctx, q, id)
	var u models.User
	var avatarURL sql.NullString
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &avatarURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	u.AvatarURL = avatarURL.String
	return &u, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	const q = "INSERT INTO users (id, email, username, avatarUrl) VALUES (?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, q, user.ID, user.Email, user.Name, user.AvatarURL)
	return err
}
