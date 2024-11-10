package controller

import (
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/services"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	c *services.Container
)

func TestMain(m *testing.M) {
	config.SwitchEnvironment(config.EnvironmentTest)

	c = services.NewContainer()

	exitCodt := m.Run()
	os.Exit(exitCodt)
}

func newContext(url string) echo.Context {
	req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(""))
	return c.Web.NewContext(req, httptest.NewRecorder())
}
