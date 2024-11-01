package middleware

import (
	"github.com/hippo-an/goranchise/auth"
	"github.com/hippo-an/goranchise/context"
	"github.com/hippo-an/goranchise/ent"
	"github.com/hippo-an/goranchise/msg"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func LoadAuthenticatedUser(authClient *auth.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if user, err := authClient.GetAuthenticatedUser(c); err == nil {
				switch err.(type) {
				case *ent.NotFoundError:
					c.Logger().Debug("auth user not found")
				case auth.NotAuthenticatedError:
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

func LoadValidPasswordToken(a *auth.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userId, err := strconv.Atoi(c.Param("userId"))
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound, "Not found")
			}

			tokenPathParam := c.Param("password_token")
			if tokenPathParam == "" {
				c.Logger().Warn("missing password token path parameter")
				return echo.NewHTTPError(http.StatusNotFound, "not fount password token")
			}

			token, err := a.GetValidPasswordToken(c, tokenPathParam, userId)
			switch err.(type) {
			case nil:
			case auth.InvalidTokenError:
				msg.Warning(c, "The link is either invalid or has expired. Please request a new one.")
				return c.Redirect(http.StatusFound, c.Echo().Reverse("forgot_password"))
			default:
				c.Logger().Error(err)
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}

			c.Set(context.PasswordTokenKey, token)
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
