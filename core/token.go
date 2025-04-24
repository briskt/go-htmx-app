package core

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/briskt/go-htmx-app/api"
	"github.com/briskt/go-htmx-app/data"
)

func FindUserByToken(ctx context.Context, tx *sql.Tx, token string) (data.User, error) {
	accessToken, err := data.FindAccessTokenByHash(ctx, tx, HashAccessToken(token))
	if err != nil {
		return data.User{}, errors.New("invalid access token")
	}

	if accessToken.ExpiresAt.Before(time.Now()) {
		return data.User{}, errors.New("expired access token")
	}

	user, err := data.GetUser(ctx, tx, int(accessToken.UserID))
	if err != nil {
		return data.User{}, fmt.Errorf("error getting authenticated user: %w", err)
	}

	return user, nil
}

// NewToken creates a new user authentication token.
func NewToken(ctx context.Context, tx *sql.Tx, staffID string) (string, error) {
	user, err := data.FindUserByEmployeeID(ctx, tx, staffID)
	if err != nil {
		err = fmt.Errorf("no user found with staff ID %q: %w", staffID, err)
		return "", api.NewAppError(err, api.ErrorNotAuthenticated, http.StatusUnauthorized)
	}

	rawToken, err := getRandomToken()
	if err != nil {
		err = fmt.Errorf("error generating random token: %w", err)
		return "", api.NewAppError(err, api.ErrorGeneratingRandomToken, http.StatusInternalServerError)
	}

	_, err = data.CreateAccessToken(ctx, tx, int(user.ID), HashAccessToken(rawToken))
	if err != nil {
		err = fmt.Errorf("error creating access token: %w", err)
		return "", api.NewAppError(err, api.ErrorCreatingAccessToken, http.StatusInternalServerError)
	}

	return rawToken, nil
}

func getRandomToken() (string, error) {
	rb := make([]byte, 32)

	_, err := rand.Read(rb)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rb), nil
}

func HashAccessToken(accessToken string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(accessToken)))
}
