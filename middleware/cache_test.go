package middleware

import (
	"context"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/hippo-an/goranchise/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestServeCachedPage(t *testing.T) {
	cp := CachedPage{
		URL:        "/cache",
		HTML:       []byte("html"),
		StatusCode: http.StatusCreated,
		Headers:    make(map[string]string),
	}

	contentTypeHeaderKey := "Content-Type"
	cp.Headers[contentTypeHeaderKey] = "text/html"
	cacheControlHeaderKey := "Cache-Control"
	cp.Headers[cacheControlHeaderKey] = "no-cache"

	err := marshaler.New(c.Cache).Set(context.Background(), cp.URL, cp)
	require.NoError(t, err)

	ctx, rec := tests.NewContext(c.Web, cp.URL)
	err = tests.ExecuteMiddleware(ctx, ServeCachedPage(c.Cache))
	assert.NoError(t, err)
	assert.Equal(t, cp.StatusCode, ctx.Response().Status)
	assert.Equal(t, cp.Headers[contentTypeHeaderKey], ctx.Response().Header().Get(contentTypeHeaderKey))
	assert.Equal(t, cp.Headers[cacheControlHeaderKey], ctx.Response().Header().Get(cacheControlHeaderKey))
	assert.Equal(t, cp.HTML, rec.Body.Bytes())

	tests.InitSession(ctx)
	err = c.Auth.Login(ctx, usr.ID)
	require.NoError(t, err)

	_ = tests.ExecuteMiddleware(ctx, LoadAuthenticatedUser(c.Auth))
	err = tests.ExecuteMiddleware(ctx, ServeCachedPage(c.Cache))
	tests.AssertHTTPErrorCode(t, err, http.StatusNotFound)
}

func TestCacheControl(t *testing.T) {
	ctx, _ := tests.NewContext(c.Web, "/")
	_ = tests.ExecuteMiddleware(ctx, CacheControl(time.Second*5))
	assert.Equal(t, "public, max-age=5", ctx.Response().Header().Get("Cache-Control"))
	_ = tests.ExecuteMiddleware(ctx, CacheControl(0))
	assert.Equal(t, "no-cache, no-store", ctx.Response().Header().Get("Cache-Control"))
}
