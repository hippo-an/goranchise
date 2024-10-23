package msg

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Type string

const (
	TypeSuccess Type = "success"
	TypeInfo    Type = "info"
	TypeWarning Type = "warning"
	TypeDanger  Type = "danger"
)

func Success(c echo.Context, message string) {
	Set(c, TypeSuccess, message)
}

func Info(c echo.Context, message string) {
	Set(c, TypeInfo, message)
}

func Warning(c echo.Context, message string) {
	Set(c, TypeWarning, message)
}

func Danger(c echo.Context, message string) {
	Set(c, TypeDanger, message)
}

func getSession(ctx echo.Context) *sessions.Session {
	sess, _ := session.Get("msg", ctx)
	return sess
}

func Set(ctx echo.Context, typ Type, message string) {
	sess := getSession(ctx)
	sess.AddFlash(message, string(typ))
	_ = sess.Save(ctx.Request(), ctx.Response())
}

func Get(c echo.Context, typ Type) []string {
	sess := getSession(c)
	fm := sess.Flashes(string(typ))
	if len(fm) > 0 {
		_ = sess.Save(c.Request(), c.Response())

		var flashes []string
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
		return flashes
	}
	return nil
}
