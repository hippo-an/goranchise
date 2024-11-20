package controller

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/go-playground/validator/v10"
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/funcmap"
	"github.com/hippo-an/goranchise/middleware"
	"github.com/hippo-an/goranchise/msg"
	"github.com/hippo-an/goranchise/services"
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
)

var (
	templates    = sync.Map{}
	funcMap      = funcmap.GetFuncMap()
	templatePath = getTemplatesDirectoryPath()
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

	if err := c.parsePageTemplate(p); err != nil {
		ctx.Logger().Errorf("failed to parse templates: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	buf, err := c.executeTemplate(ctx, p)

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

func (c *Controller) executeTemplate(ctx echo.Context, p Page) (*bytes.Buffer, error) {
	tmpl, ok := templates.Load(p.PageName)

	if !ok {
		return nil, errors.New("uncached page template requested")
	}

	buf := new(bytes.Buffer)
	err := tmpl.(*template.Template).ExecuteTemplate(buf, p.Layout+config.TemplateExt, p)

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (c *Controller) cachePage(ctx echo.Context, p Page, html *bytes.Buffer) {
	if !p.Cache.Enabled {
		return
	}
	if p.Cache.Expiration == 0 {
		p.Cache.Expiration = c.Container.Config.Cache.Expiration.Page
	}

	key := ctx.Request().URL.String()

	cp := middleware.CachedPage{
		URL:        key,
		HTML:       html.Bytes(),
		Headers:    p.Headers,
		StatusCode: p.StatusCode,
	}

	err := marshaler.New(c.Container.Cache).Set(
		ctx.Request().Context(),
		key,
		cp,
		store.WithExpiration(p.Cache.Expiration),
		store.WithTags(p.Cache.Tags),
	)

	if err != nil {
		ctx.Logger().Errorf("failed to cache page: %s", key)
		ctx.Logger().Error(err)
		return
	}
	ctx.Logger().Infof("cached page for: %s", key)

}

func (c *Controller) parsePageTemplate(p Page) error {
	if _, ok := templates.Load(p.PageName); !ok || c.Container.Config.App.Environment == config.EnvironmentLocal {
		parsed, err := template.New(p.Layout+config.TemplateExt).
			Funcs(funcMap).
			ParseFiles(
				fmt.Sprintf("%s/layouts/%s%s", templatePath, p.Layout, config.TemplateExt),
				fmt.Sprintf("%s/pages/%s%s", templatePath, p.PageName, config.TemplateExt),
			)

		if err != nil {
			return err
		}

		parsed, err = parsed.ParseGlob(fmt.Sprintf("%s/components/*%s", templatePath, config.TemplateExt))

		if err != nil {
			return err
		}

		templates.Store(p.PageName, parsed)

	}
	return nil
}

func (c *Controller) Redirect(ctx echo.Context, route string, routeParams ...interface{}) error {
	return ctx.Redirect(http.StatusFound, ctx.Echo().Reverse(route, routeParams))
}

// getTemplatesDirectoryPath gets the templates directory path
// This is needed in case this is called from a package outside of main, such as testing
func getTemplatesDirectoryPath() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Join(filepath.Dir(d), config.TemplateDir)
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
