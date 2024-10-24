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
	Ent      *ent.Client
}

func NewContainer() *Container {
	var c Container

	// web configuration
	c.Web = echo.New()

	// config configuration
	cfg, err := config.GetConfig()

	if err != nil {
		c.Web.Logger.Fatalf("failed to load configuration: %v", err)
	}

	c.Config = &cfg

	// cache configuration
	cacheClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Config.Cache.Hostname, c.Config.Cache.Port),
		Password: c.Config.Cache.Password,
	})

	//if _, err = cacheClient.Ping(context.Background()).Result(); err != nil {
	//	c.Web.Logger.Fatalf("failed to connect to cache server: %v", err)
	//}

	cacheStore := redis_store.NewRedis(cacheClient)
	c.Cache = cache.New[any](cacheStore)

	// database
	addr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		c.Config.Database.User,
		c.Config.Database.Password,
		c.Config.Database.Hostname,
		c.Config.Database.Database,
	)

	// ent
	st, err := entsql.Open("pgx", addr)
	if err != nil {
		c.Web.Logger.Fatalf("failed to connect to database: %v", err)
	}

	c.Database = st.DB()

	drv := entsql.OpenDB(dialect.Postgres, c.Database)
	c.Ent = ent.NewClient(ent.Driver(drv))
	if err := c.Ent.Schema.Create(context.Background()); err != nil {
		c.Web.Logger.Fatalf("failed to create database schema: %v", err)
	}
	return &c
}
