package msg

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Type string

const (
	Success Type = "success"
	Info    Type = "info"
	Warning Type = "warning"
	Danger  Type = "danger"
)

func getSession(ctx echo.Context) *sessions.Session {
	sess, _ := session.Get("msg", ctx)
	return sess
}

func Set(ctx echo.Context, typ Type, value string) {
	sess := getSession(ctx)
	sess.AddFlash(value, string(typ))
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
