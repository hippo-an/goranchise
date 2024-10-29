package middleware

import (
	"github.com/hippo-an/goranchise/auth"
	"github.com/hippo-an/goranchise/context"
	"github.com/hippo-an/goranchise/ent"
	"github.com/hippo-an/goranchise/ent/user"
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

func LoadAuthenticatedUser(orm *ent.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if userId, err := auth.GetUserID(c); err == nil {
				u, err := orm.User.Query().
					Where(user.ID(userId)).
					First(c.Request().Context())

				switch err.(type) {
				case *ent.NotFoundError:
					c.Logger().Debug("auth user not found: %d", userId)
				case nil:
					c.Set(context.AuthenticatedUserKey, u)
					c.Logger().Info("auth user loaded in to context: %d", userId)
				default:
					c.Logger().Errorf("error querying for authenticated user: %v", err)
				}
			}

			return next(c)
		}
	}
}
