package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/open-amt-cloud-toolkit/console/config"
	"github.com/open-amt-cloud-toolkit/console/internal/app"
)

// Function pointers for better testability.
var (
	initializeConfigFunc = config.NewConfig
	initializeAppFunc    = app.Init
	runAppFunc           = app.Run
)

func main() {
	cfg, err := initializeConfigFunc()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	err = initializeAppFunc()
	if err != nil {
		log.Fatalf("App init error: %s", err)
	}

	// root, privateKey, err := certificates.GenerateRootCertificate(true, cfg.CommonName, "US", "open-amt-cloud-toolkit", true)
	// if err != nil {
	// 	log.Fatalf("Error generating root certificate: %s", err)
	// }
	// certificates.IssueWebServerCertificate(certificates.CertAndKeyType{Cert: root, Key: privateKey}, false, cfg.CommonName, "US", "open-amt-cloud-toolkit", true)

	if os.Getenv("GIN_MODE") != "debug" {
		go func() {
			browserError := openBrowser("http://localhost:"+cfg.HTTP.Port, runtime.GOOS)
			if browserError != nil {
				panic(browserError)
			}
		}()
	}

	runAppFunc(cfg)
}

// CommandExecutor is an interface to allow for mocking exec.Command in tests.
type CommandExecutor interface {
	Execute(name string, arg ...string) error
}

// RealCommandExecutor is a real implementation of CommandExecutor.
type RealCommandExecutor struct{}

func (e *RealCommandExecutor) Execute(name string, arg ...string) error {
	return exec.Command(name, arg...).Start()
}

// Global command executor, can be replaced in tests.
var cmdExecutor CommandExecutor = &RealCommandExecutor{}

func openBrowser(url, currentOS string) error {
	var cmd string

	var args []string

	switch currentOS {
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

	return cmdExecutor.Execute(cmd, args...)
}
