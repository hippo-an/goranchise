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
	p.Name = "contact"
	p.Data = "This is contact page"

	return c.RenderPage(ctx, p)
}

func (c *Contact) Post(ctx echo.Context) error {
	msg.Set(ctx, msg.Success, "Thank you for contacting us!")
	msg.Set(ctx, msg.Info, "We will respond to you shortly")
	return c.Redirect(ctx, "home")
}
