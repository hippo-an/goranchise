package controllers

import (
	"github.com/labstack/echo/v4"
)

type About struct {
	Controller
}

func (a *About) Get(ctx echo.Context) error {
	p := NewPage(ctx)
	p.Layout = "main"
	p.PageName = "about"
	p.Title = "About"
	p.Data = "This is the about page"
	p.Cache.Enabled = true
	p.Cache.Tags = []string{"page_about", "page:list"}
	return a.RenderPage(ctx, p)

}
