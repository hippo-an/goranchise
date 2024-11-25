package middleware

import (
	"errors"
	"github.com/hippo-an/goranchise/context"
	"github.com/hippo-an/goranchise/ent"
	"github.com/hippo-an/goranchise/tests"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestLoadAuthenticatedUser(t *testing.T) {
	ctx, _ := tests.NewContext(c.Web, "/")
	tests.InitSession(ctx)
	mw := LoadAuthenticatedUser(c.Auth)

	_ = tests.ExecuteMiddleware(ctx, mw)
	assert.Nil(t, ctx.Get(context.AuthenticatedUserKey))

	err := c.Auth.Login(ctx, usr.ID)
	require.NoError(t, err)

	_ = tests.ExecuteMiddleware(ctx, mw)
	require.NotNil(t, ctx.Get(context.AuthenticatedUserKey))
	ctxUsr, ok := ctx.Get(context.AuthenticatedUserKey).(*ent.User)
	require.True(t, ok)
	require.Equal(t, usr.ID, ctxUsr.ID)
}

func TestRequireAuthentication(t *testing.T) {
	ctx, _ := tests.NewContext(c.Web, "/")
	tests.InitSession(ctx)

	err := tests.ExecuteMiddleware(ctx, RequireAuthentication())
	var httpError *echo.HTTPError
	ok := errors.As(err, &httpError)
	require.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpError.Code)

	err = c.Auth.Login(ctx, usr.ID)
	require.NoError(t, err)
	_ = tests.ExecuteMiddleware(ctx, LoadAuthenticatedUser(c.Auth))

	err = tests.ExecuteMiddleware(ctx, RequireAuthentication())
	ok = errors.As(err, &httpError)
	require.True(t, ok)
	assert.NotEqual(t, http.StatusUnauthorized, httpError.Code)
}

func TestRequireNoAuthentication(t *testing.T) {
	ctx, _ := tests.NewContext(c.Web, "/")
	tests.InitSession(ctx)

	err := tests.ExecuteMiddleware(ctx, RequireNoAuthentication())
	var httpError *echo.HTTPError
	ok := errors.As(err, &httpError)
	require.True(t, ok)
	assert.NotEqual(t, http.StatusForbidden, httpError.Code)

	err = c.Auth.Login(ctx, usr.ID)
	require.NoError(t, err)
	_ = tests.ExecuteMiddleware(ctx, LoadAuthenticatedUser(c.Auth))

	err = tests.ExecuteMiddleware(ctx, RequireNoAuthentication())
	ok = errors.As(err, &httpError)
	require.True(t, ok)
	assert.Equal(t, http.StatusForbidden, httpError.Code)
}
