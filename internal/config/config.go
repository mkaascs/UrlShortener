package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env                string `yaml:"env" env-default:"local"`
	DbConnectionString string `yaml:"db_connection_string"`
	DbHost             string `yaml:"db_host" env-required:"true"`
	DbPassword         string `yaml:"-" env-required:"true" env:"MYSQL_ROOT_PASSWORD"`
	HttpServer         `yaml:"http_server"`
	Service            `yaml:"service"`
}

type HttpServer struct {
	Address     string        `yaml:"address" address-default:"localhost:5050"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type Service struct {
	AliasLen int `yaml:"alias_len" env-default:"6"`
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		return nil, errors.New("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", cfgPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		return nil, fmt.Errorf("unable to read config: %s", err)
	}

	cfg.DbConnectionString = fmt.Sprintf(
		"root:%s@tcp(%s)/URLShortener?charset=utf8&parseTime=True",
		cfg.DbPassword,
		cfg.DbHost,
	)

	return &cfg, nil
}
