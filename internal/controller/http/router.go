// Package v1 implements routing paths. Each services in own file.
package http

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
	v1 "github.com/open-amt-cloud-toolkit/console/internal/controller/http/v1"
	v2 "github.com/open-amt-cloud-toolkit/console/internal/controller/http/v2"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

//go:embed all:ui
var content embed.FS

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
	handler.StaticFileFS("/flUhRq6tzZclQEJ-Vdg-IuiaDsNc.woff2", "./flUhRq6tzZclQEJ-Vdg-IuiaDsNc.woff2", http.FS(staticFiles))
	handler.StaticFileFS("/KFOlCnqEu92Fr1MmEU9fBBc4.woff2", "./KFOlCnqEu92Fr1MmEU9fBBc4.woff2", http.FS(staticFiles))
	handler.StaticFileFS("/KFOlCnqEu92Fr1MmSU5fBBc4.woff2", "./KFOlCnqEu92Fr1MmSU5fBBc4.woff2", http.FS(staticFiles))
	handler.StaticFileFS("/KFOmCnqEu92Fr1Mu4mxK.woff2", "./KFOmCnqEu92Fr1Mu4mxK.woff2", http.FS(staticFiles))
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
	vr := v1.NewVersionRoute(cfg)
	handler.GET("/version", vr.LatestReleaseHandler)

	// Routers
	h2 := handler.Group("/api/v1")
	{
		v1.NewDeviceRoutes(h2, t.Devices, l)
		v1.NewAmtRoutes(h2, t.Devices, t.AMTExplorer, l)
	}

	h := handler.Group("/api/v1/admin")
	{
		v1.NewDomainRoutes(h, t.Domains, l)
		v1.NewCIRAConfigRoutes(h, t.CIRAConfigs, l)
		v1.NewProfileRoutes(h, t.Profiles, l)
		v1.NewWirelessConfigRoutes(h, t.WirelessProfiles, l)
		v1.NewIEEE8021xConfigRoutes(h, t.IEEE8021xProfiles, l)
	}

	h3 := handler.Group("/api/v2")
	{
		v2.NewAmtRoutes(h3, t.Devices, l)
	}

	// Catch-all route to serve index.html for any route not matched above to be handled by Angular
	handler.NoRoute(func(c *gin.Context) {
		c.FileFromFS("./", http.FS(staticFiles)) // Serve static files from "/" route
	})
}
