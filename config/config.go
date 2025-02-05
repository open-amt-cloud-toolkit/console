package config

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v2"
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
		Name                     string        `env-required:"true" yaml:"name" env:"APP_NAME"`
		Repo                     string        `env-required:"true" yaml:"repo" env:"APP_REPO"`
		Version                  string        `env-required:"true"`
		CommonName               string        `env-required:"true" yaml:"common_name" env:"APP_COMMON_NAME"`
		EncryptionKey            string        `yaml:"encryption_key" env:"APP_ENCRYPTION_KEY"`
		JWTKey                   string        `env-required:"true" yaml:"jwtKey" env:"APP_JWT_KEY"`
		AuthDisabled             bool          `yaml:"authDisabled" env:"APP_AUTH_DISABLED"`
		AdminUsername            string        `yaml:"adminUsername" env:"APP_ADMIN_USERNAME"`
		AdminPassword            string        `yaml:"adminPassword" env:"APP_ADMIN_PASSWORD"`
		JWTExpiration            time.Duration `yaml:"jwtExpiration" env:"APP_JWT_EXPIRATION"`
		RedirectionJWTExpiration time.Duration `yaml:"redirectionJWTExpiration" env:"APP_REDIRECTION_JWT_EXPIRATION"`
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
			Name:                     "console",
			Repo:                     "open-amt-cloud-toolkit/console",
			Version:                  "DEVELOPMENT",
			CommonName:               "localhost",
			EncryptionKey:            "",
			JWTKey:                   "your_secret_jwt_key",
			AdminUsername:            "standalone",
			AdminPassword:            "G@ppm0ym",
			JWTExpiration:            24 * time.Hour,
			RedirectionJWTExpiration: 5 * time.Minute,
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
			URL:     "",
		},
		EA: EA{
			URL:      "http://localhost:8000",
			Username: "",
			Password: "",
		},
	}

	// Define a command line flag for the config path
	var configPathFlag string
	if flag.Lookup("config") == nil {
		flag.StringVar(&configPathFlag, "config", "", "path to config file")
	}

	flag.Parse()

	// Determine the config path
	var configPath string
	if configPathFlag != "" {
		configPath = configPathFlag
	} else {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}

		exPath := filepath.Dir(ex)

		configPath = filepath.Join(exPath, "config", "config.yml")
	}

	err := cleanenv.ReadConfig(configPath, ConsoleConfig)

	var pathErr *os.PathError

	if errors.As(err, &pathErr) {
		// Write config file out to disk
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
			return nil, err
		}

		file, err := os.Create(configPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		encoder := yaml.NewEncoder(file)
		defer encoder.Close()

		if err := encoder.Encode(ConsoleConfig); err != nil {
			return nil, err
		}
	}

	err = cleanenv.ReadEnv(ConsoleConfig)
	if err != nil {
		return nil, err
	}

	return ConsoleConfig, nil
}
