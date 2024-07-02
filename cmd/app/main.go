package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/open-amt-cloud-toolkit/console/config"
	"github.com/open-amt-cloud-toolkit/console/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	err = app.Init()
	if err != nil {
		log.Fatalf("App init error: %s", err)
	}

	if os.Getenv("GIN_MODE") != "debug" {
		go func() {
			browserError := openBrowser("http://localhost:" + cfg.HTTP.Port)

			if browserError != nil {
				panic(browserError)
			}
		}()
	}
	// Run
	app.Run(cfg)
}

func openBrowser(url string) error {
	var cmd string

	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}
