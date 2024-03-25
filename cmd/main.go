package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/jritsema/gotoolbox"
	"github.com/open-amt-cloud-toolkit/console/internal"
	"github.com/open-amt-cloud-toolkit/console/internal/certificates"
	"github.com/open-amt-cloud-toolkit/console/internal/dashboard"
	"github.com/open-amt-cloud-toolkit/console/internal/devices"
	"github.com/open-amt-cloud-toolkit/console/internal/i18n"
	"github.com/open-amt-cloud-toolkit/console/internal/profiles"
	"go.etcd.io/bbolt"
)

var (
	//go:embed css/output.css
	css embed.FS
	//go:embed assets/logo.png
	logo embed.FS
)

func main() {

	//exit process immediately upon sigterm
	handleSigTerms()

	//add routes
	router := http.NewServeMux()

	router.Handle("/css/output.css", http.FileServer(http.FS(css)))
	router.Handle("/assets/logo.png", http.FileServer(http.FS(logo)))
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bbolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_ = devices.NewDevices(db, router)
	_ = certificates.NewCertificates(router)
	_ = profiles.NewProfiles(db, router)
	_ = dashboard.NewDashboard(router)
	it := internal.NewIndex(router)

	//logging/tracing
	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	middleware := internal.Tracing(nextRequestID)(internal.Logging(logger)(router))

	// Setup localization
	translations, err := i18n.LoadTranslations()
	if err != nil {
		logger.Println("failed loading translations", err)
		os.Exit(1)
	}

	if err := i18n.SetupTranslations(translations); err != nil {
		logger.Println("failed setting up translations", err)
		os.Exit(1)
	}

	port := gotoolbox.GetEnvWithDefault("PORT", "8085")
	logger.Println("listening on http://localhost:" + port)

	it.Dev = flag.Bool("dev", false, "Set to true to enable development mode")
	flag.Parse()
	if *it.Dev {
		fmt.Println("Development mode enabled")
	} else {
		url := "http://localhost:" + port + "/devices"
		// Since ListenAndServe is blocking launching browser before the server is up.  Potential race condition that should be fixed.
		browserError := openBrowser(url)

		if browserError != nil {
			panic(browserError)
		}
	}

	if err := http.ListenAndServe("localhost:"+port, middleware); err != nil {
		logger.Println("http.ListenAndServe():", err)
		os.Exit(1)

	}
}

func handleSigTerms() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("received SIGTERM, exiting")
		os.Exit(1)
	}()
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
