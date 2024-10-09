package controllers

import (
	"bytes"
	"fmt"
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/container"
	"github.com/hippo-an/goranchise/funcmap"
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
	"sync"
)

const (
	TemplateDir = "views"
	TemplateExt = ".tmpl"
)

var (
	templates = sync.Map{}
	funcMap   = funcmap.GetFuncMap()
)

type Controller struct {
	Container *container.Container
}

func NewController(c *container.Container) Controller {
	return Controller{
		Container: c,
	}
}

func (c *Controller) RenderPage(ctx echo.Context, p Page) error {
	if p.Name == "" {
		ctx.Logger().Error("Page render failed due to missing name")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	if p.AppName == "" {
		p.AppName = c.Container.Config.App.Name
	}

	if err := c.parsePageTemplate(p); err != nil {
		return err
	}

	tmpl, ok := templates.Load(p.Name)

	if !ok {
		ctx.Logger().Error("Uncached page template requested")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	buf := new(bytes.Buffer)
	err := tmpl.(*template.Template).ExecuteTemplate(buf, p.Layout+TemplateExt, p)

	if err != nil {
		return err
	}

	return ctx.HTMLBlob(p.StatusCode, buf.Bytes())

}

func (c *Controller) parsePageTemplate(p Page) error {
	if _, ok := templates.Load(p.Name); !ok || c.Container.Config.App.Environment == config.EnvLocal {
		parsed, err := template.New(p.Layout+TemplateExt).
			Funcs(funcMap).
			ParseFiles(
				fmt.Sprintf("%s/layouts/%s%s", TemplateDir, p.Layout, TemplateExt),
				fmt.Sprintf("%s/pages/%s%s", TemplateDir, p.Name, TemplateExt),
			)

		if err != nil {
			return err
		}

		parsed, err = parsed.ParseGlob(fmt.Sprintf("%s/components/*%s", TemplateDir, TemplateExt))

		if err != nil {
			return err
		}

		templates.Store(p.Name, parsed)

	}
	return nil
}

func (c *Controller) Redirect(ctx echo.Context, route string, routeParams ...interface{}) error {
	return ctx.Redirect(http.StatusFound, ctx.Echo().Reverse(route, routeParams))
}
