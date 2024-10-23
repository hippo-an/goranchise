package controllers

import (
	"github.com/hippo-an/goranchise/msg"
	"github.com/labstack/echo/v4"
)

type Login struct {
	Controller
}

func (l *Login) Get(c echo.Context) error {
	p := NewPage(c)

	p.Layout = "auth"
	p.Name = "login"
	p.Title = "Login"
	p.Data = "Login Page"

	return l.RenderPage(c, p)
}

func (l *Login) Post(c echo.Context) error {
	msg.Danger(c, "Invalid credentials. Please try again.")
	return l.Get(c)
}
