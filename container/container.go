package container

import (
	"fmt"
	"github.com/eko/gocache/lib/v4/cache"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/hippo-an/goranchise/config"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Container struct {
	Web    *echo.Echo
	Cache  *cache.Cache[any]
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
	cacheClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Config.Cache.Hostname, c.Config.Cache.Port),
		Password: c.Config.Cache.Password,
	})

	//if _, err = cacheClient.Ping(context.Background()).Result(); err != nil {
	//	c.Web.Logger.Error(err)
	//	c.Web.Logger.Fatal("Failed to connect to cache server")
	//}

	cacheStore := redis_store.NewRedis(cacheClient)
	c.Cache = cache.New[any](cacheStore)

	return &c
}
