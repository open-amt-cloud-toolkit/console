package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

var ConsoleConfig *Config

type (
	// Config -.
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		DB   `yaml:"postgres"`
		EA   `yaml:"ea"`
	}

	// App -.
	App struct {
		Name          string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Repo          string `env-required:"true" yaml:"repo" env:"APP_REPO"`
		Version       string `env-required:"true"`
		EncryptionKey string `yaml:"encryption_key" env:"APP_ENCRYPTION_KEY"`
		JWTKey        string `env-required:"true" yaml:"jwtKey" env:"APP_JWT_KEY"`
		AuthDisabled  bool   `yaml:"authDisabled" env:"APP_AUTH_DISABLED"`
		AdminUsername string `yaml:"adminUsername" env:"APP_ADMIN_USERNAME"`
		AdminPassword string `yaml:"adminPassword" env:"APP_ADMIN_PASSWORD"`
	}

	// HTTP -.
	HTTP struct {
		Host           string   `env-required:"true" yaml:"host" env:"HTTP_HOST"`
		Port           string   `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		AllowedOrigins []string `env-required:"true" yaml:"allowed_origins" env:"HTTP_ALLOWED_ORIGINS"`
		AllowedHeaders []string `env-required:"true" yaml:"allowed_headers" env:"HTTP_ALLOWED_HEADERS"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// DB -.
	DB struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"DB_POOL_MAX"`
		URL     string `env:"DB_URL"`
	}

	// EA -.
	EA struct {
		URL      string `yaml:"url" env:"EA_URL"`
		Username string `yaml:"username" env:"EA_USERNAME"`
		Password string `yaml:"password" env:"EA_PASSWORD"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	// set defaults
	ConsoleConfig = &Config{
		App: App{
			Name:          "console",
			Repo:          "open-amt-cloud-toolkit/console",
			Version:       "DEVELOPMENT",
			EncryptionKey: "",
			JWTKey:        "your_secret_jwt_key",
			AdminUsername: "standalone",
			AdminPassword: "G@ppm0ym",
		},
		HTTP: HTTP{
			Host:           "localhost",
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
		EA: EA{
			URL:      "http://localhost:8000",
			Username: "",
			Password: "",
		},
	}

	_ = cleanenv.ReadConfig("./config/config.yml", ConsoleConfig)
	// its ok to ignore the error here, as we have default values set if the config file is not found

	err := cleanenv.ReadEnv(ConsoleConfig)
	if err != nil {
		return nil, err
	}

	return ConsoleConfig, nil
}
