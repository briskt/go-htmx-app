package data

import (
	"time"

	"github.com/briskt/go-htmx-app/app"
)

func (s *Suite) TestCreateAccessToken() {
	user := insertUser(s.db)
	token, err := CreateAccessToken(s.ctx, s.db, int(user.ID), "fakehash")
	s.NoError(err)
	s.Equal(user.ID, token.UserID)
	s.Equal("fakehash", token.Hash)
	s.WithinDuration(time.Now().Add(app.AccessTokenLifetime), token.ExpiresAt, time.Second)
	s.WithinDuration(time.Now(), token.CreatedUTC, time.Second)
	s.WithinDuration(time.Now(), token.UpdatedUTC, time.Second)
}

func (s *Suite) TestFindAccessTokenByHash() {
	user := insertUser(s.db)
	token, err := CreateAccessToken(s.ctx, s.db, int(user.ID), "fakehash")
	s.NoError(err)

	got, err := FindAccessTokenByHash(s.ctx, s.db, token.Hash)
	s.NoError(err)
	s.Equal(token, got)

	_, err = FindAccessTokenByHash(s.ctx, s.db, "")
	s.Error(err)
}
