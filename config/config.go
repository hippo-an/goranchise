package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
)

type Environment string

const (
	EnvironmentLocal Environment = "local"
	EnvironmentTest  Environment = "test"
	EnvironmentProd  Environment = "prod"
)

const (
	StaticDir = "static"
	PublicDir = "public"
)

const (
	TemplateDir = "templates"
	TemplateExt = ".tmpl"
)

type (
	Config struct {
		Http     HttpConfig
		App      AppConfig
		Cache    CacheConfig
		Database DatabaseConfig
		Mail     MailConfig
	}

	HttpConfig struct {
		Hostname     string
		Port         uint16
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
		TLS          struct {
			Enabled     bool
			Certificate string
			Key         string
		}
	}

	AppConfig struct {
		Name          string
		Environment   Environment
		EncryptionKey string
		Timeout       time.Duration
		PasswordToken struct {
			Expiration time.Duration
			Length     int
		}
	}

	CacheConfig struct {
		Hostname   string
		Port       uint16
		Password   string
		Expiration struct {
			StaticFile time.Duration
			Page       time.Duration
		}
	}

	DatabaseConfig struct {
		Hostname     string
		Port         uint16
		User         string
		Password     string
		Database     string
		TestDatabase string
	}

	MailConfig struct {
		Hostname    string
		Port        uint16
		User        string
		Password    string
		FromAddress string
	}
)

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

	log.Printf("%+v\n", viper.AllSettings())

	return c, nil
}

func SwitchEnvironment(env Environment) {
	SetEnvironmentVariable("GORANCHISE_APP_ENVIRONMENT", string(env))
}

func SwitchDatabaseHostAndPort(host, port string) {
	SetEnvironmentVariable("GORANCHISE_DATABASE_HOSTNAME", host)
	SetEnvironmentVariable("GORANCHISE_DATABASE_PORT", port)
}

func SwitchCacheHostAndPort(host, port string) {
	SetEnvironmentVariable("GORANCHISE_CACHE_HOSTNAME", host)
	SetEnvironmentVariable("GORANCHISE_CACHE_PORT", port)
}

func SetEnvironmentVariable(key, value string) {
	if err := os.Setenv(key, value); err != nil {
		panic(err)
	}
}
