package repository

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type FirebaseAuthRepo struct {
	app *firebase.App
}

func NewFirebaseAuthRepo(app *firebase.App) *FirebaseAuthRepo {
	return &FirebaseAuthRepo{app: app}
}

func (r *FirebaseAuthRepo) CreateUser(ctx context.Context, email, password string) (string, error) {
	client, err := r.app.Auth(ctx)
	if err != nil {
		return "", err
	}
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)
	u, err := client.CreateUser(ctx, params)
	if err != nil {
		return "", err
	}
	return u.UID, nil
}

func (r *FirebaseAuthRepo) DeleteUser(ctx context.Context, uid string) error {
	client, err := r.app.Auth(ctx)
	if err != nil {
		return err
	}
	return client.DeleteUser(ctx, uid)
}
