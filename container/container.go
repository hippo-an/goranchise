package container

import (
	"github.com/hippo-an/goranchise/config"
	"github.com/labstack/echo/v4"
)

type Container struct {
	Web    *echo.Echo
	Config *config.Config
}

func NewContainer() *Container {
	var c Container

	c.Web = echo.New()

	cfg, err := config.GetConfig()

	if err != nil {
		c.Web.Logger.Fatal("Failed to load configuration")
		panic(err)
	}

	c.Config = &cfg

	return &c
}
