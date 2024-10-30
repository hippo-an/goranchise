package middleware

import (
	"github.com/hippo-an/goranchise/auth"
	"github.com/hippo-an/goranchise/context"
	"github.com/hippo-an/goranchise/ent"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RequireAuthentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if u := c.Get(context.AuthenticatedUserKey); u == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			return next(c)
		}
	}
}

func RequireNoAuthentication() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if u := c.Get(context.AuthenticatedUserKey); u != nil {
				return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
			}
			return next(c)
		}
	}
}

func LoadAuthenticatedUser(authClient *auth.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if user, err := authClient.GetAuthenticatedUser(c); err == nil {
				switch err.(type) {
				case *ent.NotFoundError:
					c.Logger().Debug("auth user not found")
				case nil:
					c.Set(context.AuthenticatedUserKey, user)
					c.Logger().Info("auth user loaded in to context: %d", user.ID)
				default:
					c.Logger().Errorf("error querying for authenticated user: %v", err)
				}
			}

			return next(c)
		}
	}
}
