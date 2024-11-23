package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewContainer(t *testing.T) {
	c := NewContainer()
	assert.NotNil(t, c.Web)
	assert.NotNil(t, c.Config)
	assert.NotNil(t, c.Cache)
	assert.NotNil(t, c.Database)
	assert.NotNil(t, c.ORM)
	assert.NotNil(t, c.Mail)
	assert.NotNil(t, c.Auth)
	assert.NotNil(t, c.TemplateRenderer)
}
