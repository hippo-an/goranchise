package middleware

import (
	"context"
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/ent"
	"github.com/hippo-an/goranchise/services"
	"github.com/hippo-an/goranchise/tests"
	"os"
	"testing"
)

var (
	c   *services.Container
	usr *ent.User
)

func TestMain(m *testing.M) {
	config.SwitchEnvironment(config.EnvironmentTest)
	con, err := tests.RunTestDB()
	if err != nil {
		panic(err)
	}
	defer con.Terminate(context.Background())

	c = services.NewContainer()

	defer func() {
		if err := c.Shutdown(); err != nil {
			c.Web.Logger.Fatal(err)
		}
	}()

	if usr, err = tests.CreateUser(c.ORM); err != nil {
		panic(err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}
