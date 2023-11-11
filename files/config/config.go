package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App  `yaml:"app"`
	HTTP `yaml:"http"`
	Log  `yaml:"logger"`
	DB   `yaml:"db"`
}

type App struct {
	Name       string `env-required:"true" yaml:"name" env:"APP_NAME"`
	Version    string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	NumWorkers int    `env-required:"true" yaml:"numworkers" env:"NUM_WORKERS"`
}

type HTTP struct {
	Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
}

type Log struct {
	Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
}

type DB struct {
	Host     string `env-required:"true" yaml:"db_host"   env:"DB_HOST"`
	Port     string `env-required:"true" yaml:"db_port"   env:"DB_PORT"`
	Username string `env-required:"true" yaml:"db_username"   env:"DB_USERNAME"`
	Password string `env-required:"true" yaml:"db_password"   env:"DB_PASSWORD"`
	Name     string `env-required:"true" yaml:"db_name"   env:"DB_NAME"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./files/config/config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
