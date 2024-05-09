package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		DB   `yaml:"postgres"`
	}

	// App -.
	App struct {
		Name string `env-required:"false" yaml:"name" env:"APP_NAME"`
	}

	// HTTP -.
	HTTP struct {
		Port           string   `env-required:"false" yaml:"port" env:"HTTP_PORT"`
		AllowedOrigins []string `env-required:"false" yaml:"allowed_origins" env:"HTTP_ALLOWED_ORIGINS"`
		AllowedHeaders []string `env-required:"false" yaml:"allowed_headers" env:"HTTP_ALLOWED_HEADERS"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"false" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// DB -.
	DB struct {
		PoolMax int    `env-required:"false" yaml:"pool_max" env:"DB_POOL_MAX"`
		URL     string `env:"DB_URL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	// set defaults
	cfg := &Config{
		App: App{
			Name: "console",
		},
		HTTP: HTTP{
			Port:           "8181",
			AllowedOrigins: []string{"*"},
			AllowedHeaders: []string{"*"},
		},
		Log: Log{
			Level: "info",
		},
		DB: DB{
			PoolMax: 2,
		},
	}

	_ = cleanenv.ReadConfig("./config/config.yml", cfg)
	// its ok to ignore the error here, as we have default values set if the config file is not found

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
