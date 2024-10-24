package container

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/eko/gocache/lib/v4/cache"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/ent"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Container struct {
	Web      *echo.Echo
	Cache    *cache.Cache[any]
	Config   *config.Config
	Database *sql.DB
	ORM      *ent.Client
}

func (c *Container) initWeb() {
	c.Web = echo.New()
}

func (c *Container) initConfig() {
	cfg, err := config.GetConfig()

	if err != nil {
		c.Web.Logger.Fatalf("failed to load configuration: %v", err)
	}

	c.Config = &cfg
}

func (c *Container) initCache() {
	cacheClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Config.Cache.Hostname, c.Config.Cache.Port),
		Password: c.Config.Cache.Password,
	})

	//if _, err = cacheClient.Ping(context.Background()).Result(); err != nil {
	//	c.Web.Logger.Fatalf("failed to connect to cache server: %v", err)
	//}

	cacheStore := redis_store.NewRedis(cacheClient)
	c.Cache = cache.New[any](cacheStore)
}

func (c *Container) initDatabase() {
	addr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		c.Config.Database.User,
		c.Config.Database.Password,
		c.Config.Database.Hostname,
		c.Config.Database.Database,
	)

	driver, err := entsql.Open("pgx", addr)
	if err != nil {
		c.Web.Logger.Fatalf("failed to connect to database: %v", err)
	}

	c.Database = driver.DB()
}

func (c *Container) initORM() {
	drv := entsql.OpenDB(dialect.Postgres, c.Database)
	c.ORM = ent.NewClient(ent.Driver(drv))
	if err := c.ORM.Schema.Create(context.Background()); err != nil {
		c.Web.Logger.Fatalf("failed to create database schema: %v", err)
	}
}

func NewContainer() *Container {
	var c Container
	c.initWeb()
	c.initConfig()
	c.initCache()
	c.initDatabase()
	c.initORM()
	return &c
}
