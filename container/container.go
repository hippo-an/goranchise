package container

import (
	"fmt"
	"github.com/hippo-an/goranchise/config"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Container struct {
	Web    *echo.Echo
	Cache  *redis.Client
	Config *config.Config
}

func NewContainer() *Container {
	var c Container

	// web configuration
	c.Web = echo.New()

	// config configuration
	cfg, err := config.GetConfig()

	if err != nil {
		c.Web.Logger.Error(err)
		c.Web.Logger.Fatal("Failed to load configuration")
		panic(err)
	}

	c.Config = &cfg

	// cache configuration
	c.Cache = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Config.Cache.Hostname, c.Config.Cache.Port),
		Password: c.Config.Cache.Password,
	})

	//if _, err = c.Cache.Ping(context.Background()).Result(); err != nil {
	//	c.Web.Logger.Error(err)
	//	c.Web.Logger.Fatal("Failed to connect to cache server")
	//}

	return &c
}
