package controllers

import (
	"github.com/hippo-an/goranchise/auth"
	"github.com/labstack/echo/v4"
)

type Home struct {
	Controller
}

func (h *Home) Get(c echo.Context) error {
	p := NewPage(c)

	p.Layout = "main"
	p.PageName = "home"
	p.Data = "Hello world"
	p.IsHome = true

	uid, _ := auth.GetUserID(c)
	c.Logger().Infof("logged in user ID: %d", uid)

	return h.RenderPage(c, p)
}
