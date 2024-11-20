package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetConfig(t *testing.T) {
	_, err := GetConfig()
	require.NoError(t, err)

	var env Environment
	env = "abc"
	SwitchEnvironment(env)
	cfg, err := GetConfig()
	require.NoError(t, err)
	assert.Equal(t, env, cfg.App.Environment)
}
