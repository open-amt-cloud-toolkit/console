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
	wsv1 "github.com/open-amt-cloud-toolkit/console/internal/controller/ws/v1"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase"
	"github.com/open-amt-cloud-toolkit/console/pkg/httpserver"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
	"github.com/open-amt-cloud-toolkit/console/pkg/postgres"
)

var Version = "DEVELOPMENT"

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)
	log.Info("app - Run - version: " + Version)
	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use case
	usecases := usecase.NewUseCases(pg, log)

	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// HTTP Server
	handler := gin.New()

	defaultConfig := cors.DefaultConfig()
	defaultConfig.AllowOrigins = cfg.HTTP.AllowedOrigins
	defaultConfig.AllowHeaders = cfg.HTTP.AllowedHeaders

	handler.Use(cors.New(defaultConfig))
	v1.NewRouter(handler, log, *usecases)
	wsv1.RegisterRoutes(handler, log, usecases.Devices)
	httpServer := httpserver.New(handler, httpserver.Port("127.0.0.1", cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
