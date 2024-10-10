package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/security"

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

	if os.Getenv("GIN_MODE") != "debug" {
		go func() {
			browserError := openBrowser("http://localhost:"+cfg.HTTP.Port, runtime.GOOS)
			if browserError != nil {
				panic(browserError)
			}
		}()
	}

	toolkitCrypto := security.Crypto{}
	// check if key is provided by user (unencrypted storage of key)
	if cfg.EncryptionKey == "" {
		secureStorage := security.NewKeyRingStorage("device-management-toolkit")

		// check if key exist in credential store
		cfg.EncryptionKey, err = secureStorage.GetKeyValue("default-security-key")
		if err != nil {
			if err.Error() == "The specified item could not be found in the keyring" {
				fmt.Print("\033[31mWarning: Key Not Found, Generate new key? -- This will prevent access to existing data? Y/N: \033[0m")

				var response string

				_, err = fmt.Scanln(&response)
				if err != nil {
					log.Fatal(err)

					return
				}

				if response != "Y" && response != "y" {
					log.Fatal("Exiting without generating a new key.")

					return
				}

				// generate 32 byte key for encryption
				cfg.EncryptionKey = toolkitCrypto.GenerateKey()

				err = secureStorage.SetKeyValue("default-security-key", cfg.EncryptionKey)
				if err != nil {
					log.Fatal(err)

					return
				}
			} else {
				log.Fatal(err)

				return
			}
		}
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
