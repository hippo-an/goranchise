package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/controller"
	"github.com/hippo-an/goranchise/middleware"
	"github.com/hippo-an/goranchise/services"
	"github.com/labstack/echo-contrib/session"
	echomw "github.com/labstack/echo/v4/middleware"
	"net/http"
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

func BuildRouter(c *services.Container) {

	c.Web.Group("", middleware.CacheControl(c.Config.Cache.Expiration.StaticFile)).
		Static("/public", config.PublicDir)
	c.Web.Group("", middleware.CacheControl(c.Config.Cache.Expiration.StaticFile)).
		Static("/static", config.StaticDir)

	if c.Config.Http.TLS.Enabled {
		c.Web.Use(echomw.HTTPSRedirect())
	}

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
		middleware.LoadAuthenticatedUser(c.Auth),
		middleware.ServeCachedPage(c.Cache),
		session.Middleware(sessions.NewCookieStore([]byte(c.Config.App.EncryptionKey))),
		echomw.CSRFWithConfig(echomw.CSRFConfig{
			TokenLookup: "form:csrf",
		}),
	)

	c.Web.Validator = &Validator{validator: validator.New()}

	ctr := controller.NewController(c)

	errorHandler := Error{
		Controller: ctr,
	}

	c.Web.HTTPErrorHandler = errorHandler.Handler

	navRoutes(c, ctr)
	userRoutes(c, ctr)
}

func navRoutes(c *services.Container, ctr controller.Controller) {
	home := Home{Controller: ctr}
	c.Web.GET("/", home.Get).Name = "home"

	about := About{Controller: ctr}
	c.Web.GET("/about", about.Get).Name = "about"

	contact := Contact{Controller: ctr}
	c.Web.GET("/contact", contact.Get).Name = "contact"
	c.Web.POST("/contact", contact.Post).Name = "contact.post"
}

func userRoutes(c *services.Container, ctr controller.Controller) {
	{
		logout := Logout{Controller: ctr}
		c.Web.GET("/logout", logout.Get, middleware.RequireAuthentication()).
			Name = "logout"
	}

	noAuth := c.Web.Group("/user", middleware.RequireNoAuthentication())
	{
		login := Login{Controller: ctr}
		noAuth.GET("/login", login.Get).Name = "login"
		noAuth.POST("/login", login.Post).Name = "login.post"

		register := Register{Controller: ctr}
		noAuth.GET("/register", register.Get).Name = "register"
		noAuth.POST("/register", register.Post).Name = "register.post"

		forgot := ForgotPassword{Controller: ctr}
		noAuth.GET("/password", forgot.Get).Name = "forgot_password"
		noAuth.POST("/password", forgot.Post).Name = "forgot_password.post"

		resetGroup := noAuth.Group(
			"/password/reset",
			middleware.LoadUser(c.ORM),
			middleware.LoadValidPasswordToken(c.Auth),
		)
		reset := ResetPassword{Controller: ctr}
		resetGroup.GET("/token/:userId/:password_token", reset.Get).Name = "reset_password"
		resetGroup.POST("/token/:userId/:password_token", reset.Post).Name = "reset_password.post"
	}

}
