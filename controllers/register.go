package controllers

import (
	"github.com/hippo-an/goranchise/msg"
	"github.com/labstack/echo/v4"
)

type Register struct {
	Controller
}

func (r *Register) Get(c echo.Context) error {
	p := NewPage(c)

	p.Layout = "auth"
	p.PageName = "register"
	p.Title = "Register"
	p.Data = "Register page"

	return r.RenderPage(c, p)
}

func (r *Register) Post(c echo.Context) error {
	u, err := r.Container.Ent.User.
		Create().
		SetUsername(c.FormValue("username")).
		SetPassword(c.FormValue("password")).
		Save(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
	} else {
		c.Logger().Infof("user created: %s", u.Username)
	}
	msg.Danger(c, "Registration is currently disabled.")
	return r.Get(c)
}
