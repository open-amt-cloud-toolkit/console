// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/open-amt-cloud-toolkit/console/config"
	v1 "github.com/open-amt-cloud-toolkit/console/internal/controller/http/v1"
	"github.com/open-amt-cloud-toolkit/console/internal/controller/tcp/cira"
	wsv1 "github.com/open-amt-cloud-toolkit/console/internal/controller/ws/v1"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase"
	"github.com/open-amt-cloud-toolkit/console/pkg/db"
	"github.com/open-amt-cloud-toolkit/console/pkg/httpserver"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

var Version = "DEVELOPMENT"

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)
	cfg.App.Version = Version
	log.Info("app - Run - version: " + cfg.App.Version)
	// Repository
	database, err := db.New(cfg.DB.URL, db.MaxPoolSize(cfg.DB.PoolMax))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - db.New: %w", err))
	}
	defer database.Close()

	// Use case
	usecases := usecase.NewUseCases(database, log)

	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// HTTP Server
	handler := gin.New()

	defaultConfig := cors.DefaultConfig()
	defaultConfig.AllowOrigins = cfg.HTTP.AllowedOrigins
	defaultConfig.AllowHeaders = cfg.HTTP.AllowedHeaders

	handler.Use(cors.New(defaultConfig))
	v1.NewRouter(handler, log, *usecases, cfg)
	wsv1.RegisterRoutes(handler, log, usecases.Devices)

	ciraServer, err := cira.NewServer("config/cert.pem", "config/key.pem")
	if err != nil {
		log.Fatal("CIRA Server failed: %v", err)
	}

	httpServer := httpserver.New(handler, httpserver.Port("", cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case ciraErr := <-ciraServer.Notify():
		log.Error(fmt.Errorf("app - Run - ciraServer.Notify: %w", ciraErr))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
	err = ciraServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - ciraServer.Shutdown: %w", err))
	}
}
