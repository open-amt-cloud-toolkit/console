package app_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/config"
	"github.com/open-amt-cloud-toolkit/console/internal/app"
)

func getFreePort() (string, error) {
	port, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(":%d", 8000+port.Int64()), nil
}

func teardown() {}

func TestRun(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	defer ctrl.Finish()

	mockDB := NewMockDB(ctrl)
	mockHTTPServer := NewMockHTTPServer(ctrl)

	port, err := getFreePort()
	if err != nil {
		t.Fatalf("Failed to get a free port: %v", err)
	}

	cfg := &config.Config{
		Log: config.Log{
			Level: "info",
		},
		DB: config.DB{
			URL:     "postgres://testuser:testpass@localhost/testdb",
			PoolMax: 10,
		},
		HTTP: config.HTTP{
			Port:           port,
			AllowedOrigins: []string{"*"},
			AllowedHeaders: []string{"Content-Type"},
		},
		App: config.App{
			Version: "DEVELOPMENT",
		},
	}

	tests := []struct {
		name       string
		setupMocks func()
		setupEnv   func()
		cfg        *config.Config
		expectFunc func(t *testing.T)
	}{
		{
			name: "Successful run and shutdown",
			setupMocks: func() {
				mockDB.EXPECT().Close().Return(nil).Times(1)
				mockHTTPServer.EXPECT().Notify().Return(make(chan error)).Times(1)
				mockHTTPServer.EXPECT().Shutdown().Return(nil).Times(1)
			},
			setupEnv: func() {
				os.Setenv("GIN_MODE", "release")
			},
			cfg: cfg,
			expectFunc: func(_ *testing.T) {
				go func() {
					defer teardown()
					app.Run(cfg)
				}()
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.setupEnv()
			tc.setupMocks()

			tc.expectFunc(t)
		})
	}
}
