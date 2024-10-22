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
		Static("/public", "public")
	c.Web.Group("", middleware.CacheControl(15552000)).
		Static("/images", "static/images")

	ctr := controllers.NewController(c)

	errorHandler := controllers.Error{
		Controller: ctr,
	}

	c.Web.HTTPErrorHandler = errorHandler.Handler

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
	login := controllers.Login{Controller: ctr}
	register := controllers.Register{Controller: ctr}
	userRoute := e.Group("/user")
	{
		userRoute.GET("/login", login.Get).Name = "login"
		userRoute.POST("/login", login.Post).Name = "login.post"

		userRoute.GET("/register", register.Get).Name = "register"
		userRoute.POST("/register", register.Post).Name = "register.post"
	}

}
