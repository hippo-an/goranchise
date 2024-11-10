package services

import (
	"errors"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuth(t *testing.T) {
	mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
	handler := mw(echo.NotFoundHandler)
	assert.Error(t, handler(ctx))

	assertNoAuth := func() {
		_, err := c.Auth.GetAuthenticatedUserId(ctx)
		assert.True(t, errors.Is(err, NotAuthenticatedError{}))
		_, err = c.Auth.GetAuthenticatedUser(ctx)
		assert.True(t, errors.Is(err, NotAuthenticatedError{}))
	}

	assertNoAuth()

	err := c.Auth.Login(ctx, usr.ID)
	require.NoError(t, err)

	id, err := c.Auth.GetAuthenticatedUserId(ctx)
	require.NoError(t, err)
	require.Equal(t, usr.ID, id)

	user, err := c.Auth.GetAuthenticatedUser(ctx)
	require.NoError(t, err)
	require.Equal(t, usr.ID, user.ID)

	err = c.Auth.Logout(ctx)
	require.NoError(t, err)

	assertNoAuth()
}

func TestHashPassword(t *testing.T) {
	pw := "abcdef"
	hash, err := c.Auth.HashPassword(pw)
	assert.NoError(t, err)
	assert.NotEqual(t, hash, pw)
	err = c.Auth.CheckPassword(pw, hash)
	assert.NoError(t, err)
}

func TestCheckPassword(t *testing.T) {
	pw := "testcheckpassword"
	hash, err := c.Auth.HashPassword(pw)
	assert.NoError(t, err)
	err = c.Auth.CheckPassword(pw, hash)
	assert.NoError(t, err)
}

func TestGetValidPasswordToken(t *testing.T) {
	_, err := c.Auth.GetValidPasswordToken(ctx, "notValidToken:)", usr.ID)
	assert.True(t, errors.Is(err, InvalidPasswordTokenError{}))
}

func TestGeneratePasswordResetToken(t *testing.T) {
	token, pt, err := c.Auth.GeneratePasswordResetToken(ctx, usr.ID)
	require.NoError(t, err)
	assert.Len(t, token, c.Config.App.PasswordToken.Length)
	assert.NoError(t, c.Auth.CheckPassword(token, pt.Hash))

	pt2, err := c.Auth.GetValidPasswordToken(ctx, token, usr.ID)
	assert.NoError(t, err)
	assert.Equal(t, pt.ID, pt2.ID)
}

func TestRandomToken(t *testing.T) {
	length := 64
	a, err := c.Auth.RandomToken(length)
	require.NoError(t, err)
	b, err := c.Auth.RandomToken(length)
	require.NoError(t, err)
	assert.Len(t, a, length)
	assert.Len(t, b, length)
	assert.NotEqual(t, a, b)
}
