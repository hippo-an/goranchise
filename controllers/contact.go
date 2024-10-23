package controllers

import (
	"github.com/hippo-an/goranchise/msg"
	"github.com/labstack/echo/v4"
)

type Contact struct {
	Controller
}

func (c *Contact) Get(ctx echo.Context) error {
	p := NewPage(ctx)
	p.Layout = "main"
	p.PageName = "contact"
	p.Title = "Contact us"
	p.Data = "This is contact page"

	return c.RenderPage(ctx, p)
}

func (c *Contact) Post(ctx echo.Context) error {
	msg.Success(ctx, "Thank you for contacting us!")
	msg.Info(ctx, "We will respond to you shortly")
	return c.Redirect(ctx, "home")
}
