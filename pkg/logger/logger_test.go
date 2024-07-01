package logger

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mu sync.Mutex

type loggerTest struct {
	name         string
	logLevel     string
	testFunction func(t *testing.T, log *Logger, buf *bytes.Buffer)
}

func TestLogger(t *testing.T) {
	t.Parallel()

	tests := []loggerTest{
		{
			name:     "Debug level logging",
			logLevel: "debug",
			testFunction: func(t *testing.T, log *Logger, buf *bytes.Buffer) {
				t.Helper()
				log.Debug("debug message")
				assert.Contains(t, buf.String(), "debug message")
			},
		},
		{
			name:     "Info level logging",
			logLevel: "info",
			testFunction: func(t *testing.T, log *Logger, buf *bytes.Buffer) {
				t.Helper()
				log.Info("info message")
				assert.Contains(t, buf.String(), "info message")
			},
		},
		{
			name:     "Warn level logging",
			logLevel: "warn",
			testFunction: func(t *testing.T, log *Logger, buf *bytes.Buffer) {
				t.Helper()
				log.Warn("warn message")
				assert.Contains(t, buf.String(), "warn message")
			},
		},
		{
			name:     "Error level logging",
			logLevel: "error",
			testFunction: func(t *testing.T, log *Logger, buf *bytes.Buffer) {
				t.Helper()
				log.Error("error message")
				assert.Contains(t, buf.String(), "error message")
			},
		},
		{
			name:     "Fatal level logging",
			logLevel: "fatal",
			testFunction: func(t *testing.T, log *Logger, _ *bytes.Buffer) {
				t.Helper()
				if os.Getenv("BE_CRASHER") == "1" {
					log.Fatal("fatal message")

					return
				}
				cmd := exec.Command(os.Args[0], "-test.run=TestLogger") // #nosec
				cmd.Env = append(os.Environ(), "BE_CRASHER=1")
				err := cmd.Run()
				var exitError *exec.ExitError
				if errors.As(err, &exitError) && !exitError.Success() {
					return
				}
				t.Fatalf("process ran with err %v, want exit status 1", err)
			},
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			logger := zerolog.New(&buf).With().Timestamp().Logger()
			log := &Logger{logger: &logger}

			tc.testFunction(t, log, &buf)
		})
	}
}

func TestNewLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level         string
		expectedLevel zerolog.Level
	}{
		{"debug", zerolog.DebugLevel},
		{"info", zerolog.InfoLevel},
		{"warn", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"invalid", zerolog.InfoLevel},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(fmt.Sprintf("LogLevel_%s", tc.level), func(t *testing.T) {
			t.Parallel()

			mu.Lock()
			defer mu.Unlock()

			log := New(tc.level)
			require.NotNil(t, log)
			assert.Equal(t, tc.expectedLevel, log.localLevel)
		})
	}
}
