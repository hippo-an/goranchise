package container

import (
	"context"
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/ent"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	c   *Container
	ctx echo.Context
	usr *ent.User
	rec *httptest.ResponseRecorder
)

func TestMain(m *testing.M) {
	// Set the environment to test
	config.SwitchEnvironment(config.EnvironmentTest)

	// Create a new container
	c = NewContainer()

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	rec = httptest.NewRecorder()
	ctx = c.Web.NewContext(req, rec)

	var err error
	usr, err = c.ORM.User.
		Create().
		SetEmail("test@test.dev").
		SetPassword("abc").
		SetName("Test User").
		Save(context.Background())
	if err != nil {
		panic(err)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestNewContainer(t *testing.T) {
	c := NewContainer()
	assert.NotNil(t, c.Web)
	assert.NotNil(t, c.Config)
	assert.NotNil(t, c.Cache)
	assert.NotNil(t, c.Database)
	assert.NotNil(t, c.ORM)
	assert.NotNil(t, c.Mail)
	assert.NotNil(t, c.Auth)
}
