package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Env string

const (
	EnvTest  Env = "test"
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type Config struct {
	Http  HttpConfig
	App   AppConfig
	Cache CacheConfig
}

type HttpConfig struct {
	Hostname     string
	Port         uint16
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type AppConfig struct {
	Name          string
	Environment   Env
	EncryptionKey string
}

type CacheConfig struct {
	Hostname string
	Port     uint16
	Password string
	MaxAge   struct {
		StaticFile int
	}
}

func GetConfig() (Config, error) {
	var c Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	viper.SetEnvPrefix("goranchise")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return c, err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return c, err
	}

	return c, nil
}
