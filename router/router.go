package router

import (
	"github.com/gorilla/sessions"
	"github.com/hippo-an/goranchise/container"
	"github.com/hippo-an/goranchise/controllers"
	"github.com/hippo-an/goranchise/middleware"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echomv "github.com/labstack/echo/v4/middleware"
	"net/http"
)

const (
	StaticDir = "static"
	PublicDir = "public"
)

func BuildRouter(c *container.Container) {
	c.Web.Use(
		echomv.RemoveTrailingSlashWithConfig(
			echomv.TrailingSlashConfig{
				RedirectCode: http.StatusMovedPermanently,
			}),
		echomv.RequestID(),
		echomv.Recover(),
		echomv.Gzip(),
		echomv.Logger(),
		//middleware.Static(StaticDir),
		session.Middleware(sessions.NewCookieStore([]byte(c.Config.App.EncryptionKey))),
		echomv.CSRFWithConfig(echomv.CSRFConfig{
			TokenLookup: "form:csrf",
		}),
	)

	c.Web.Group("", middleware.CacheControl(15552000)).
		Static("/", StaticDir)
	c.Web.Group("", middleware.CacheControl(15552000)).
		Static("/", PublicDir)

	ctr := controllers.NewController(c)

	err := controllers.Error{
		Controller: ctr,
	}

	c.Web.HTTPErrorHandler = err.Handler

	navRoutes(c.Web, ctr)
	userRoutes(c.Web, ctr)
}

func navRoutes(e *echo.Echo, ctr controllers.Controller) {
	home := controllers.Home{Controller: ctr}
	e.GET("/", home.Get).Name = "home"

	about := controllers.About{Controller: ctr}
	e.GET("/about", about.Get).Name = "about"

	contact := controllers.Contact{Controller: ctr}
	e.GET("/contact", contact.Get).Name = "contact"
	e.POST("/contact", contact.Post).Name = "contact.post"
}
func userRoutes(e *echo.Echo, ctr controllers.Controller) {
}
