package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jritsema/go-htmx-starter/internal"
	"github.com/jritsema/go-htmx-starter/internal/companies"
	"github.com/jritsema/gotoolbox"
)

var (
	//go:embed css/output.css
	css embed.FS
)

func main() {

	//exit process immediately upon sigterm
	handleSigTerms()

	//add routes
	router := http.NewServeMux()

	router.Handle("/css/output.css", http.FileServer(http.FS(css)))

	_ = companies.NewCompanies(router)
	_ = internal.NewIndex(router)

	//logging/tracing
	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	middleware := internal.Tracing(nextRequestID)(internal.Logging(logger)(router))

	port := gotoolbox.GetEnvWithDefault("PORT", "8080")
	logger.Println("listening on http://localhost:" + port)
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
