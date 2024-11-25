package services

import (
	"context"
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/ent"
	"github.com/hippo-an/goranchise/tests"
	"github.com/labstack/echo/v4"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	c   *Container
	ctx echo.Context
	usr *ent.User
	rec *httptest.ResponseRecorder
)

func TestMain(m *testing.M) {
	config.SwitchEnvironment(config.EnvironmentTest)

	con, err := tests.RunTestDB()
	if err != nil {
		panic(err)
	}
	defer con.Terminate(context.Background())

	c = NewContainer()

	defer func() {
		if err := c.Shutdown(); err != nil {
			c.Web.Logger.Fatal(err)
		}
	}()

	ctx, _ = tests.NewContext(c.Web, "/")
	tests.InitSession(ctx)

	usr, err = tests.CreateUser(c.ORM)
	if err != nil {
		panic(err)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}
