package controller

import (
	"github.com/gorilla/sessions"
	"github.com/hippo-an/goranchise/context"
	"github.com/hippo-an/goranchise/msg"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
)

func TestNewPage(t *testing.T) {
	homeUrl := "/"
	ctx := newContext(homeUrl)
	p := NewPage(ctx)

	assert.Same(t, ctx, p.Context)
	assert.NotNil(t, p.ToURL)
	assert.Equal(t, homeUrl, p.Path)
	assert.Equal(t, homeUrl, p.URL)
	assert.Equal(t, http.StatusOK, p.StatusCode)
	assert.Equal(t, NewPager(ctx, DefaultItemsPerPage), p.Pager)
	assert.Empty(t, p.Headers)
	assert.True(t, p.IsHome)
	assert.False(t, p.IsAuth)
	assert.Empty(t, p.CSRF)
	assert.Empty(t, p.RequestId)
	assert.False(t, p.Cache.Enabled)

	url := "/abc?def=123"
	csrfToken := "csrf"
	ctx = newContext(url)
	ctx.Set(context.AuthenticatedUserKey, 1)
	ctx.Set(echomw.DefaultCSRFConfig.ContextKey, csrfToken)
	p = NewPage(ctx)

	assert.Equal(t, strings.Split(url, "?")[0], p.Path)
	assert.Equal(t, url, p.URL)
	assert.False(t, p.IsHome)
	assert.True(t, p.IsAuth)
	assert.Equal(t, csrfToken, p.CSRF)
}

func TestPage_GetMessages(t *testing.T) {
	ctx := newContext("/")

	p := NewPage(ctx)
	mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
	handler := mw(echo.NotFoundHandler)
	err := handler(ctx)
	assert.Error(t, err)
	require.ErrorIs(t, err, echo.ErrNotFound)

	msgTests := make(map[msg.Type][]string)
	msgTests[msg.TypeWarning] = []string{
		"abc",
		"def",
	}
	msgTests[msg.TypeInfo] = []string{
		"123",
		"456",
	}

	for typ, values := range msgTests {
		for _, value := range values {
			msg.Set(ctx, typ, value)
		}
	}
	for typ, values := range msgTests {
		msgs := p.GetMessages(typ)
		for i, message := range msgs {
			assert.Equal(t, values[i], string(message))
		}
	}
}
