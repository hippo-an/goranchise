package middleware

import (
	"errors"
	"fmt"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/hippo-an/goranchise/context"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type CachedPage struct {
	URL        string
	HTML       []byte
	StatusCode int
	Headers    map[string]string
}

func ServeCachedPage(ch *cache.Cache[any]) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if c.Request().Method != http.MethodGet {
				return next(c)
			}

			if c.Get(context.AuthenticatedUserKey) != nil {
				return next(c)
			}

			res, err := marshaler.New(ch).Get(c.Request().Context(), c.Request().URL.String(), new(CachedPage))

			if err != nil {
				if errors.Is(err, redis.Nil) {
					c.Logger().Info("no cached page found")
				} else {
					c.Logger().Errorf("failed getting cached page: %v", err)
				}
				return next(c)
			}

			page, ok := res.(*CachedPage)

			if !ok {
				c.Logger().Info("failed casting cached page")
				return next(c)
			}

			if page.Headers != nil {
				for k, v := range page.Headers {
					c.Response().Header().Set(k, v)
				}
			}
			c.Logger().Info("serving cached page")
			return c.HTMLBlob(page.StatusCode, page.HTML)
		}
	}
}

func CacheControl(maxAge time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			v := "no-cache, no-store"

			if maxAge > 0 {
				v = fmt.Sprintf("public,max-age=%.0f", maxAge.Seconds())
			}

			c.Response().Header().Set("Cache-Control", v)
			return next(c)
		}
	}
}
