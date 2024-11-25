package middleware

import (
	"github.com/hippo-an/goranchise/context"
	"github.com/hippo-an/goranchise/ent"
	"github.com/hippo-an/goranchise/msg"
	"github.com/hippo-an/goranchise/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

func LoadAuthenticatedUser(authClient *services.AuthClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if user, err := authClient.GetAuthenticatedUser(c); err == nil {
				switch err.(type) {
				case *ent.NotFoundError:
					c.Logger().Debug("auth user not found")
				case services.NotAuthenticatedError:
				case nil:
					c.Set(context.AuthenticatedUserKey, user)
					c.Logger().Infof("auth user loaded in to context: %d", user.ID)
				default:
					c.Logger().Errorf("error querying for authenticated user: %v", err)
				}
			}

			return next(c)
		}
	}
}

func LoadValidPasswordToken(a *services.AuthClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var u *ent.User
			cu := c.Get(context.UserKey)
			if cu == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}
			u = cu.(*ent.User)

			tokenPathParam := c.Param("password_token")
			if tokenPathParam == "" {
				c.Logger().Warn("missing password token path parameter")
				return echo.NewHTTPError(http.StatusNotFound, "not fount password token")
			}

			_, err := a.GetValidPasswordToken(c, tokenPathParam, u.ID)
			switch err.(type) {
			case nil:
			case services.InvalidPasswordTokenError:
				msg.Warning(c, "The link is either invalid or has expired. Please request a new one.")
				return c.Redirect(http.StatusFound, c.Echo().Reverse("forgot_password"))
			default:
				c.Logger().Error(err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}

			return next(c)
		}
	}
}

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
