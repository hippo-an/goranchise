package controllers

import "github.com/labstack/echo/v4"

type Login struct {
	Controller
}

func (h *Login) Get(c echo.Context) error {
	p := NewPage(c)

	p.Layout = "auth"
	p.Name = "login"
	p.Title = "Login"
	p.Data = "Login Page"

	return h.RenderPage(c, p)
}
