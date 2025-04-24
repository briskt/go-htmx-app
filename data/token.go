package data

import (
	"context"
	"fmt"
	"time"

	"github.com/briskt/go-htmx-app/app"
	"github.com/briskt/go-htmx-app/data/sqlc"
)

type AccessToken struct {
	sqlc.Token
}

func CreateAccessToken(ctx context.Context, tx sqlc.DBTX, userID int, tokenHash string) (AccessToken, error) {
	params := sqlc.CreateAccessTokenParams{
		UserID:     int32(userID),
		Hash:       tokenHash,
		ExpiresAt:  time.Now().Add(app.AccessTokenLifetime),
		CreatedUTC: time.Now(),
		UpdatedUTC: time.Now(),
	}
	if app.Env.AppEnv == app.EnvDevelopment {
		params.ExpiresAt = time.Now().Add(app.AccessTokenLifetime * 100)
	}
	token, err := q(tx).CreateAccessToken(ctx, params)
	if err != nil {
		err = fmt.Errorf("error creating access token: %w", err)
		return AccessToken{}, err
	}

	return AccessToken{token}, nil
}

func FindAccessTokenByHash(ctx context.Context, tx sqlc.DBTX, hash string) (AccessToken, error) {
	token, err := q(tx).FindAccessTokenByHash(ctx, hash)
	if err != nil {
		err = fmt.Errorf("error finding access token: %w", err)
		return AccessToken{}, err
	}
	return AccessToken{token}, nil
}
