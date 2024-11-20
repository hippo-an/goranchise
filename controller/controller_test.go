package controller

import (
	"context"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/middleware"
	"github.com/hippo-an/goranchise/msg"
	"github.com/hippo-an/goranchise/services"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	c *services.Container
)

func TestMain(m *testing.M) {
	config.SwitchEnvironment(config.EnvironmentTest)

	c = services.NewContainer()

	defer func() {
		if err := c.Shutdown(); err != nil {
			panic(err)
		}
	}()
	exitCode := m.Run()

	os.Exit(exitCode)
}

func newContext(url string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, url, strings.NewReader(""))
	rec := httptest.NewRecorder()
	return c.Web.NewContext(req, rec), rec
}

func initSession(t *testing.T, ctx echo.Context) {
	mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
	handler := mw(echo.NotFoundHandler)
	err := handler(ctx)
	assert.Error(t, err)
	assert.ErrorIs(t, err, echo.ErrNotFound)
}

func TestController_Redirect(t *testing.T) {
	ctx, _ := newContext("/abc")
	ctr := NewController(c)
	err := ctr.Redirect(ctx, "home")
	require.NoError(t, err)
	assert.Equal(t, "", ctx.Response().Header().Get(echo.HeaderLocation))
	assert.Equal(t, http.StatusFound, ctx.Response().Status)
}

func TestController_SetValidationErrorMessages(t *testing.T) {
	type example struct {
		Name string `validate:"required" label:"Label test"`
	}

	e := example{}
	v := validator.New()
	err := v.Struct(e)
	require.Error(t, err)

	ctx, _ := newContext("/")
	initSession(t, ctx)
	ctr := NewController(c)
	ctr.SetValidationErrorMessages(ctx, err, e)
	msgs := msg.Get(ctx, msg.TypeDanger)
	require.Len(t, msgs, 1)
	assert.Equal(t, "<strong>Label test</strong> is required.", msgs[0])
}

func TestController_RenderPage(t *testing.T) {
	setup := func() (echo.Context, *httptest.ResponseRecorder, Controller, Page) {
		ctx, rec := newContext("/test/TestController_RenderPage")
		initSession(t, ctx)
		ctr := NewController(c)
		p := NewPage(ctx)
		p.PageName = "home"
		p.Layout = "main"
		p.Cache.Enabled = false
		p.Headers["a"] = "b"
		p.Headers["c"] = "d"
		p.StatusCode = http.StatusCreated
		return ctx, rec, ctr, p
	}

	t.Run("missing name", func(t *testing.T) {
		ctx, _, ctr, p := setup()
		p.PageName = ""
		err := ctr.RenderPage(ctx, p)
		assert.Error(t, err)
	})

	t.Run("no page cache", func(t *testing.T) {
		ctx, _, ctr, p := setup()
		err := ctr.RenderPage(ctx, p)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, ctx.Response().Status)
		for k, v := range p.Headers {
			assert.Equal(t, v, ctx.Response().Header().Get(k))
		}

		parsed, err := c.Templates.Load("controller", p.PageName)
		assert.NoError(t, err)

		expectedTemplates := make(map[string]bool)
		expectedTemplates[p.PageName+config.TemplateExt] = true
		expectedTemplates[p.Layout+config.TemplateExt] = true
		components, err := os.ReadDir(c.Templates.GetTemplatesPath() + "/components")
		require.NoError(t, err)
		for _, f := range components {
			expectedTemplates[f.Name()] = true
		}
		for _, v := range parsed.Templates() {
			delete(expectedTemplates, v.Name())
		}
		assert.Empty(t, expectedTemplates)
	})

	t.Run("page cache", func(t *testing.T) {
		ctx, rec, ctr, p := setup()
		p.Cache.Enabled = true
		p.Cache.Tags = []string{"tag1"}
		err := ctr.RenderPage(ctx, p)
		require.NoError(t, err)

		res, err := marshaler.New(c.Cache).
			Get(context.Background(), p.URL, new(middleware.CachedPage))
		require.NoError(t, err)

		cp, ok := res.(*middleware.CachedPage)
		require.True(t, ok)
		assert.Equal(t, p.URL, cp.URL)
		assert.Equal(t, p.Headers, cp.Headers)
		assert.Equal(t, p.StatusCode, cp.StatusCode)
		assert.Equal(t, rec.Body.Bytes(), cp.HTML)

		err = c.Cache.Invalidate(context.Background(), store.WithInvalidateTags([]string{p.Cache.Tags[0]}))
		require.NoError(t, err)

		_, err = marshaler.New(c.Cache).
			Get(context.Background(), p.URL, new(middleware.CachedPage))
		assert.Error(t, err)
	})
}
