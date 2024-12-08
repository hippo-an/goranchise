package tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/gorilla/sessions"
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/ent"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func NewContext(e *echo.Echo, url string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, url, strings.NewReader(""))
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func InitSession(ctx echo.Context) {
	mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
	_ = ExecuteMiddleware(ctx, mw)
}

func ExecuteMiddleware(ctx echo.Context, mw echo.MiddlewareFunc) error {
	handler := mw(echo.NotFoundHandler)
	return handler(ctx)
}

func CreateUser(orm *ent.Client) (*ent.User, error) {
	seed := fmt.Sprintf("%d-%d", time.Now().UnixMilli(), rand.IntN(1000000))
	return orm.User.
		Create().
		SetEmail(fmt.Sprintf("testuser-%s@localhost.localhost", seed)).
		SetPassword("password").
		SetName(fmt.Sprintf("Test User %s", seed)).
		Save(context.Background())
}

func RunTestCache() (testcontainers.Container, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}

	cacheContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatal(err)
	}

	err = SwitchCacheToContainer(cacheContainer)
	if err != nil {
		return nil, err
	}

	return cacheContainer, nil
}

// RunTestDB
// https://golang.testcontainers.org/system_requirements/using_colima/
// using colima to docker container runtime in mac os
// use colima to default context for docker
func RunTestDB() (testcontainers.Container, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "admin",
			"POSTGRES_PASSWORD": "admin1234",
			"POSTGRES_DB":       "goranchise_test",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatal(err)
	}

	err = SwitchDatabaseToContainer(postgresContainer)
	if err != nil {
		return nil, err
	}

	return postgresContainer, nil
}

func SwitchDatabaseToContainer(con testcontainers.Container) error {
	host, err := con.Host(context.Background())

	if err != nil {
		return err
	}

	port, err := con.MappedPort(context.Background(), "5432")
	if err != nil {
		return err
	}

	config.SwitchDatabaseHostAndPort(host, port.Port())
	return nil
}

func SwitchCacheToContainer(con testcontainers.Container) error {
	host, err := con.Host(context.Background())

	if err != nil {
		return err
	}

	port, err := con.MappedPort(context.Background(), "6379")
	if err != nil {
		return err
	}

	config.SwitchCacheHostAndPort(host, port.Port())
	return nil
}

func AssertHTTPErrorCode(t *testing.T, err error, code int) {
	var httpError *echo.HTTPError
	ok := errors.As(err, &httpError)
	require.True(t, ok)
	assert.Equal(t, code, httpError.Code)
}

func AssertHTTPErrorCodeNot(t *testing.T, err error, code int) {
	var httpError *echo.HTTPError
	ok := errors.As(err, &httpError)
	require.True(t, ok)
	assert.NotEqual(t, code, httpError.Code)
}
