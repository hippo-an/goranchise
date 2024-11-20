package controller

import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/go-playground/validator/v10"
	"github.com/hippo-an/goranchise/middleware"
	"github.com/hippo-an/goranchise/msg"
	"github.com/hippo-an/goranchise/services"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
)

type Controller struct {
	Container *services.Container
}

func NewController(c *services.Container) Controller {
	return Controller{
		Container: c,
	}
}

func (c *Controller) RenderPage(ctx echo.Context, p Page) error {
	if p.PageName == "" {
		ctx.Logger().Error("page render failed due to missing name")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	if p.AppName == "" {
		p.AppName = c.Container.Config.App.Name
	}

	if err := c.parsePageTemplates(p); err != nil {
		ctx.Logger().Errorf("failed to parse templates: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	buf, err := c.executeTemplate(p)

	if err != nil {
		ctx.Logger().Errorf("failed to execute templates: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	c.cachePage(ctx, p, buf)

	for k, v := range p.Headers {
		ctx.Response().Header().Set(k, v)
	}

	return ctx.HTMLBlob(p.StatusCode, buf.Bytes())
}

func (c *Controller) executeTemplate(page Page) (*bytes.Buffer, error) {
	return c.Container.Templates.Execute("controller", page.PageName, page.Layout, page)
}

func (c *Controller) cachePage(ctx echo.Context, page Page, html *bytes.Buffer) {
	if !page.Cache.Enabled {
		return
	}
	if page.Cache.Expiration == 0 {
		page.Cache.Expiration = c.Container.Config.Cache.Expiration.Page
	}

	key := ctx.Request().URL.String()

	cp := middleware.CachedPage{
		URL:        key,
		HTML:       html.Bytes(),
		Headers:    page.Headers,
		StatusCode: page.StatusCode,
	}

	err := marshaler.New(c.Container.Cache).Set(
		ctx.Request().Context(),
		key,
		cp,
		store.WithExpiration(page.Cache.Expiration),
		store.WithTags(page.Cache.Tags),
	)

	if err != nil {
		ctx.Logger().Errorf("failed to cache page: %s", key)
		ctx.Logger().Error(err)
		return
	}

	ctx.Logger().Infof("cached page for: %s", key)
}

func (c *Controller) parsePageTemplates(page Page) error {
	return c.Container.Templates.Parse(
		"controller",
		page.PageName,
		page.Layout,
		[]string{
			fmt.Sprintf("layouts/%s", page.Layout),
			fmt.Sprintf("pages/%s", page.PageName),
		},
		[]string{"components"},
	)
}

func (c *Controller) Redirect(ctx echo.Context, route string, routeParams ...interface{}) error {
	return ctx.Redirect(http.StatusFound, ctx.Echo().Reverse(route, routeParams))
}

func (c *Controller) SetValidationErrorMessages(ctx echo.Context, err error, data interface{}) {

	var ves validator.ValidationErrors
	ok := errors.As(err, &ves)
	if !ok {
		return
	}

	for _, ve := range ves {
		var message string
		label := ve.StructField()
		ty := reflect.TypeOf(data)

		if ty.Kind() == reflect.Ptr {
			ty = ty.Elem()
		}

		if ty.Kind() == reflect.Struct {
			if field, ok := ty.FieldByName(ve.Field()); ok {
				if labelTag := field.Tag.Get("label"); labelTag != "" {
					label = labelTag
				}
			}
		}

		switch ve.Tag() {
		case "required":
			message = "%s is required."
		default:
			message = "%s is not a valid value."
		}

		msg.Danger(ctx, fmt.Sprintf(message, "<strong>"+label+"</strong>"))
	}
}
