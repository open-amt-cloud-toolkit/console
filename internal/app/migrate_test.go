package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupEnv   func()
		expectFunc func(t *testing.T, err error)
	}{
		{
			name: "Postgres DB setup success",
			setupEnv: func() {
				os.Setenv("DB_URL", "postgres://testuser:testpass@localhost/testdb")
			},

			expectFunc: func(t *testing.T, err error) {
				t.Helper()

				require.NoError(t, err)
			},
		},
		{
			name: "SQLite DB setup success",
			setupEnv: func() {
				os.Unsetenv("DB_URL")
			},

			expectFunc: func(t *testing.T, err error) {
				t.Helper()

				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.setupEnv()

			err := Init()

			tc.expectFunc(t, err)
		})
	}
}
