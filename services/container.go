package services

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/eko/gocache/lib/v4/cache"
	redisstore "github.com/eko/gocache/store/redis/v4"
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/ent"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

type Container struct {
	Web              *echo.Echo
	Cache            *cache.Cache[any]
	Config           *config.Config
	cacheClient      *redis.Client
	Database         *sql.DB
	ORM              *ent.Client
	Auth             *AuthClient
	TemplateRenderer *TemplateRenderer
	Mail             *MailClient
}

func NewContainer() *Container {
	c := new(Container)
	c.initConfig()
	c.initWeb()
	c.initCache()
	c.initDatabase()
	c.initORM()
	c.initAuth()
	c.initTemplateRenderer()
	c.initMail()
	return c
}

func (c *Container) Shutdown() error {
	if err := c.cacheClient.Close(); err != nil {
		return err
	}

	if err := c.ORM.Close(); err != nil {
		return err
	}

	if err := c.Database.Close(); err != nil {
		return err
	}

	return nil
}

func (c *Container) initWeb() {
	c.Web = echo.New()

	switch c.Config.App.Environment {
	case config.EnvironmentProd:
		c.Web.Logger.SetLevel(log.WARN)
	default:
		c.Web.Logger.SetLevel(log.DEBUG)
	}
	c.Web.HideBanner = true
}

func (c *Container) initConfig() {
	cfg, err := config.GetConfig()

	if err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}

	c.Config = &cfg
}

func (c *Container) initCache() {
	c.cacheClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Config.Cache.Hostname, c.Config.Cache.Port),
		Password: c.Config.Cache.Password,
	})

	if c.Config.App.Environment == config.EnvironmentProd {
		if _, err := c.cacheClient.Ping(context.Background()).Result(); err != nil {
			panic(fmt.Sprintf("failed to connect to cache server: %v", err))
		}
	}

	cacheStore := redisstore.NewRedis(c.cacheClient)
	c.Cache = cache.New[any](cacheStore)
}

func (c *Container) initDatabase() {

	getAddr := func(dbName string) string {
		return fmt.Sprintf("postgresql://%s:%s@%s/%s",
			c.Config.Database.User,
			c.Config.Database.Password,
			c.Config.Database.Hostname,
			dbName,
		)
	}

	switch c.Config.App.Environment {
	case config.EnvironmentTest:
		driver, err := entsql.Open("pgx", getAddr(c.Config.Database.TestDatabase))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to database: %v", err))
		}
		c.Database = driver.DB()
	case config.EnvironmentLocal:
		driver, err := entsql.Open("pgx", getAddr(c.Config.Database.Database))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to database: %v", err))
		}
		c.Database = driver.DB()
	case config.EnvironmentProd:
	default:
	}
}

func (c *Container) initORM() {
	drv := entsql.OpenDB(dialect.Postgres, c.Database)
	c.ORM = ent.NewClient(ent.Driver(drv))
	if err := c.ORM.Schema.Create(context.Background()); err != nil {
		panic(fmt.Sprintf("failed to create database schema: %v", err))
	}
}

func (c *Container) initAuth() {
	c.Auth = NewAuthClient(c.Config, c.ORM)
}

func (c *Container) initTemplateRenderer() {
	c.TemplateRenderer = NewTemplateRenderer(c.Config)
}

func (c *Container) initMail() {
	c.Mail = NewMailClient(c.Config, c.TemplateRenderer)
}
