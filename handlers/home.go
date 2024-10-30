package handlers

import (
	"github.com/hippo-an/goranchise/controller"
	"github.com/labstack/echo/v4"
)

type Home struct {
	controller.Controller
}

func (h *Home) Get(c echo.Context) error {
	p := controller.NewPage(c)

	p.Layout = "main"
	p.PageName = "home"
	p.Data = "Hello world"
	p.IsHome = true

	return h.RenderPage(c, p)
}
