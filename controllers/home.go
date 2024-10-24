package controllers

import "github.com/labstack/echo/v4"

type Home struct {
	Controller
}

func (h *Home) Get(c echo.Context) error {
	p := NewPage(c)

	p.Layout = "main"
	p.PageName = "home"
	p.Data = "Hello world"
	p.IsHome = true
	return h.RenderPage(c, p)
}
