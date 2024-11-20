package controller

import (
	"github.com/hippo-an/goranchise/context"
	"github.com/hippo-an/goranchise/msg"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestNewPage(t *testing.T) {
	homeUrl := "/"
	ctx, _ := newContext(homeUrl)
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
	ctx, _ = newContext(url)
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
	ctx, _ := newContext("/")

	p := NewPage(ctx)
	initSession(t, ctx)

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
