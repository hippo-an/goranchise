package container

import (
	"github.com/hippo-an/goranchise/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	config.SwitchEnvironment(config.EnvironmentTest)

	exitCode := m.Run()
	os.Exit(exitCode)
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
