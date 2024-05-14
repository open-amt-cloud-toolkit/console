// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/open-amt-cloud-toolkit/console/config"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

//go:embed all:ui
var content embed.FS

var Config *config.Config

// NewRouter -.
// Swagger spec:
// @title       Console API for Device Management Toolkit
// @description Provides a single pane of glass for managing devices with Intel® Active Management Technology and other device technologies
// @version     1.0
// @host        localhost:8181
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.Usecases, cfg *config.Config) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	Config = cfg
	// Static files
	// Serve static assets (js, css, images, etc.)
	// Create subdirectory view of the embedded file system
	staticFiles, err := fs.Sub(content, "ui")
	if err != nil {
		log.Fatal(err)
	}

	// Set up HTTP server to handle requests
	handler.StaticFileFS("/", "./", http.FS(staticFiles)) // Serve static files from "/" route
	handler.StaticFileFS("/main.js", "./main.js", http.FS(staticFiles))
	handler.StaticFileFS("/polyfills.js", "./polyfills.js", http.FS(staticFiles))
	handler.StaticFileFS("/runtime.js", "./runtime.js", http.FS(staticFiles))
	handler.StaticFileFS("/styles.css", "./styles.css", http.FS(staticFiles))
	handler.StaticFileFS("/vendor.js", "./vendor.js", http.FS(staticFiles))
	handler.StaticFileFS("/favicon.ico", "./favicon.ico", http.FS(staticFiles))
	handler.StaticFileFS("/assets/logo.png", "./assets/logo.png", http.FS(staticFiles))

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// version info
	handler.GET("/version", LatestReleaseHandler)

	// Routers
	h2 := handler.Group("/api/v1")
	{
		newDeviceRoutes(h2, t.Devices, l)
		newAmtRoutes(h2, t.Devices, l)
	}

	h := handler.Group("/api/v1/admin")
	{
		newDomainRoutes(h, t.Domains, l)
		newCIRAConfigRoutes(h, t.CIRAConfigs, l)
		newProfileRoutes(h, t.Profiles, l)
		newWirelessConfigRoutes(h, t.WirelessProfiles, l)
		newIEEE8021xConfigRoutes(h, t.IEEE8021xProfiles, l)
	}

	// Catch-all route to serve index.html for any route not matched above to be handled by Angular
	handler.NoRoute(func(c *gin.Context) {
		c.FileFromFS("./", http.FS(staticFiles)) // Serve static files from "/" route
	})
}
