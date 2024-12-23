package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/hippo-an/goranchise/handlers"
	"github.com/hippo-an/goranchise/services"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	c := services.NewContainer()

	defer func() {
		if err := c.Shutdown(); err != nil {
			c.Web.Logger.Fatal(err)
		}
	}()

	handlers.BuildRouter(c)

	go func() {
		srv := http.Server{
			Addr:         fmt.Sprintf("%s:%d", c.Config.Http.Hostname, c.Config.Http.Port),
			Handler:      c.Web,
			ReadTimeout:  c.Config.Http.ReadTimeout,
			WriteTimeout: c.Config.Http.WriteTimeout,
			IdleTimeout:  c.Config.Http.IdleTimeout,
		}

		if c.Config.Http.TLS.Enabled {
			certs, err := tls.LoadX509KeyPair(c.Config.Http.TLS.Certificate, c.Config.Http.TLS.Key)
			if err != nil {
				c.Web.Logger.Fatalf("cannot load TLS certificate: %v", err)
			}

			srv.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{certs},
			}
		}

		if err := c.Web.StartServer(&srv); !errors.Is(err, http.ErrServerClosed) {
			c.Web.Logger.Fatalf("shutting down the server: %v", err)

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
