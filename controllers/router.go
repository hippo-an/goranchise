package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/hippo-an/goranchise/container"
	"github.com/hippo-an/goranchise/middleware"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"net/http"
)

const (
	StaticDir = "static"
	PublicDir = "public"
)

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func BuildRouter(c *container.Container) {

	c.Web.Group("", middleware.CacheControl(c.Config.Cache.Expiration.StaticFile)).
		Static("/public", PublicDir)
	c.Web.Group("", middleware.CacheControl(c.Config.Cache.Expiration.StaticFile)).
		Static("/static", StaticDir)

	c.Web.Use(
		echomw.RemoveTrailingSlashWithConfig(
			echomw.TrailingSlashConfig{
				RedirectCode: http.StatusMovedPermanently,
			}),
		echomw.Secure(),
		echomw.RequestID(),
		echomw.Recover(),
		echomw.Gzip(),
		echomw.Logger(),
		middleware.LogRequestId(),
		//middleware.Static(StaticDir),
		echomw.TimeoutWithConfig(echomw.TimeoutConfig{
			Timeout: c.Config.App.Timeout,
		}),
		middleware.PageCache(c.Cache),
		session.Middleware(sessions.NewCookieStore([]byte(c.Config.App.EncryptionKey))),
		echomw.CSRFWithConfig(echomw.CSRFConfig{
			TokenLookup: "form:csrf",
		}),
		middleware.LoadAuthenticatedUser(c.ORM),
	)

	c.Web.Validator = &Validator{validator: validator.New()}

	ctr := NewController(c)

	errorHandler := Error{
		Controller: ctr,
	}

	c.Web.HTTPErrorHandler = errorHandler.Handler

	navRoutes(c.Web, ctr)
	userRoutes(c.Web, ctr)
}

func navRoutes(e *echo.Echo, ctr Controller) {
	home := Home{Controller: ctr}
	e.GET("/", home.Get).Name = "home"

	about := About{Controller: ctr}
	e.GET("/about", about.Get).Name = "about"

	contact := Contact{Controller: ctr}
	e.GET("/contact", contact.Get).Name = "contact"
	e.POST("/contact", contact.Post).Name = "contact.post"
}

func userRoutes(e *echo.Echo, ctr Controller) {
	login := Login{Controller: ctr}
	register := Register{Controller: ctr}
	{
		logout := Logout{Controller: ctr}
		e.GET("/logout", logout.Get, middleware.RequireAuthentication()).
			Name = "logout"
	}

	noAuth := e.Group("/user", middleware.RequireNoAuthentication())

	{
		noAuth.GET("/login", login.Get).Name = "login"
		noAuth.POST("/login", login.Post).Name = "login.post"

		noAuth.GET("/register", register.Get).Name = "register"
		noAuth.POST("/register", register.Post).Name = "register.post"
	}

}
