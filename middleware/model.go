package middleware

import (
	"github.com/hippo-an/goranchise/context"
	"github.com/hippo-an/goranchise/ent"
	"github.com/hippo-an/goranchise/ent/user"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func LoadUser(orm *ent.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			userId, err := strconv.Atoi(c.Param("userId"))
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			u, err := orm.User.Query().
				Where(user.ID(userId)).
				Only(c.Request().Context())

			switch err.(type) {
			case nil:
				c.Set(context.UserKey, u)
				return next(c)
			case *ent.NotFoundError:
				return echo.NewHTTPError(http.StatusNotFound)
			default:
				c.Logger().Error(err)
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
		}
	}
}
