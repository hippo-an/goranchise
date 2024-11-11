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

	exitCode := m.Run()
	os.Exit(exitCode)
}

func newContext(url string) echo.Context {
	req := httptest.NewRequest(http.MethodGet, url, strings.NewReader(""))
	return c.Web.NewContext(req, httptest.NewRecorder())
}
