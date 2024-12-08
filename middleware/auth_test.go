package middleware

import (
	"errors"
	"fmt"
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
	tests.AssertHTTPErrorCodeNot(t, err, http.StatusForbidden)

	err = c.Auth.Login(ctx, usr.ID)
	require.NoError(t, err)
	_ = tests.ExecuteMiddleware(ctx, LoadAuthenticatedUser(c.Auth))

	err = tests.ExecuteMiddleware(ctx, RequireNoAuthentication())
	tests.AssertHTTPErrorCode(t, err, http.StatusForbidden)
}

func TestLoadValidPasswordToken(t *testing.T) {
	ctx, _ := tests.NewContext(c.Web, "/")
	tests.InitSession(ctx)

	err := tests.ExecuteMiddleware(ctx, LoadValidPasswordToken(c.Auth))
	tests.AssertHTTPErrorCode(t, err, http.StatusInternalServerError)

	userId, passwordToken := "userId", "password_token"

	ctx.SetParamNames(userId)
	ctx.SetParamValues(fmt.Sprintf("%d", usr.ID))
	_ = tests.ExecuteMiddleware(ctx, LoadUser(c.ORM))
	err = tests.ExecuteMiddleware(ctx, LoadValidPasswordToken(c.Auth))
	assert.Error(t, err)
	tests.AssertHTTPErrorCode(t, err, http.StatusNotFound)

	ctx.SetParamNames(userId, passwordToken)
	ctx.SetParamValues(fmt.Sprintf("%d", usr.ID), "faketoken")
	_ = tests.ExecuteMiddleware(ctx, LoadUser(c.ORM))
	err = tests.ExecuteMiddleware(ctx, LoadValidPasswordToken(c.Auth))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, ctx.Response().Status)

	token, pt, err := c.Auth.GeneratePasswordResetToken(ctx, usr.ID)
	require.NoError(t, err)

	ctx.SetParamNames(userId, passwordToken)
	ctx.SetParamValues(fmt.Sprintf("%d", usr.ID), token)
	_ = tests.ExecuteMiddleware(ctx, LoadUser(c.ORM))
	err = tests.ExecuteMiddleware(ctx, LoadValidPasswordToken(c.Auth))
	tests.AssertHTTPErrorCode(t, err, http.StatusNotFound)
	ctxPt, ok := ctx.Get(context.PasswordTokenKey).(*ent.PasswordToken)
	require.True(t, ok)
	assert.Equal(t, pt.ID, ctxPt.ID)
	assert.Equal(t, pt.Hash, ctxPt.Hash)
}
