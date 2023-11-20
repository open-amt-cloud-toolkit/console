package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/jritsema/go-htmx-starter/internal"
	"github.com/jritsema/go-htmx-starter/internal/certificates"
	"github.com/jritsema/go-htmx-starter/internal/dashboard"
	"github.com/jritsema/go-htmx-starter/internal/devices"
	"github.com/jritsema/go-htmx-starter/internal/profiles"
	"github.com/jritsema/gotoolbox"
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

	_ = internal.NewIndex(router)

	//logging/tracing
	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	middleware := internal.Tracing(nextRequestID)(internal.Logging(logger)(router))
	port := gotoolbox.GetEnvWithDefault("PORT", "8080")
	logger.Println("listening on http://localhost:" + port)
	url := "http://localhost:" + port + "/devices"
	// Since ListenAndServe is blocking launching browser before the server is up.  Potential race condition that should be fixed.
	browserError := openBrowser(url)

	if browserError != nil {
		panic(browserError)
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
