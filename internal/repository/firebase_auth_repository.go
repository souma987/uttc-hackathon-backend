package repository

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

type FirebaseAuthRepo struct {
	client *auth.Client
}

func NewFirebaseAuthRepo(client *auth.Client) *FirebaseAuthRepo {
	return &FirebaseAuthRepo{client: client}
}

func (r *FirebaseAuthRepo) CreateUser(ctx context.Context, email, password string) (string, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)
	u, err := r.client.CreateUser(ctx, params)
	if err != nil {
		return "", err
	}
	return u.UID, nil
}

func (r *FirebaseAuthRepo) DeleteUser(ctx context.Context, uid string) error {
	return r.client.DeleteUser(ctx, uid)
}

// VerifyIDToken validates a Firebase ID token and returns the UID.
func (r *FirebaseAuthRepo) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	token, err := r.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", err
	}
	return token.UID, nil
}
