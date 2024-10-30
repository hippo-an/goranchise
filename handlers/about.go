package handlers

import (
	"github.com/hippo-an/goranchise/controller"
	"github.com/labstack/echo/v4"
)

type About struct {
	controller.Controller
}

func (a *About) Get(ctx echo.Context) error {
	p := controller.NewPage(ctx)
	p.Layout = "main"
	p.PageName = "about"
	p.Title = "About"
	p.Data = "This is the about page"
	p.Cache.Enabled = true
	p.Cache.Tags = []string{"page_about", "page:list"}
	return a.RenderPage(ctx, p)

}
