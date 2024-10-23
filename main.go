package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/container"
	"github.com/hippo-an/goranchise/controllers"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	c := container.NewContainer()

	switch c.Config.App.Environment {
	case config.EnvProd:
		c.Web.Logger.SetLevel(log.WARN)
	default:
		c.Web.Logger.SetLevel(log.DEBUG)
	}

	controllers.BuildRouter(c)

	go func() {
		srv := http.Server{
			Addr:         fmt.Sprintf("%s:%d", c.Config.Http.Hostname, c.Config.Http.Port),
			Handler:      c.Web,
			ReadTimeout:  c.Config.Http.ReadTimeout,
			WriteTimeout: c.Config.Http.WriteTimeout,
			IdleTimeout:  c.Config.Http.IdleTimeout,
		}

		if err := c.Web.StartServer(&srv); !errors.Is(err, http.ErrServerClosed) {
			c.Web.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	if err := c.Web.Shutdown(ctx); err != nil {
		c.Web.Logger.Fatal(err)
	}
}
