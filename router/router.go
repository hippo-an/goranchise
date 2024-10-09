package router

import (
	"github.com/gorilla/sessions"
	"github.com/hippo-an/goranchise/container"
	"github.com/hippo-an/goranchise/controllers"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

const StaticDir = "static"

func BuildRouter(c *container.Container) {
	c.Web.Use(
		middleware.RemoveTrailingSlashWithConfig(
			middleware.TrailingSlashConfig{
				RedirectCode: http.StatusMovedPermanently,
			}),
		middleware.RequestID(),
		middleware.Recover(),
		middleware.Gzip(),
		middleware.Logger(),
		middleware.Static(StaticDir),
		session.Middleware(sessions.NewCookieStore([]byte(c.Config.App.EncryptionKey))),
	)

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
