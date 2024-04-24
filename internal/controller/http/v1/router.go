// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/open-amt-cloud-toolkit/console/internal/usecase"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.Repositories) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	// Static files
	// Serve static assets (js, css, images, etc.)
	handler.StaticFile("/", "./ui/index.html")
	handler.StaticFile("/main.js", "./ui/main.js")
	handler.StaticFile("/polyfills.js", "./ui/polyfills.js")
	handler.StaticFile("/runtime.js", "./ui/runtime.js")
	handler.StaticFile("/styles.css", "./ui/styles.css")
	handler.StaticFile("/vendor.js", "./ui/vendor.js")
	handler.StaticFile("/favicon.ico", "./ui/favicon.ico")
	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h2 := handler.Group("/api/v1")
	{
		newDeviceRoutes(h2, t.Devices, l)
		newAmtRoutes(h2, t.DeviceManagement, t.Devices, l)
	}

	h := handler.Group("/api/v1/admin")
	{
		newDomainRoutes(h, t.Domains, l)
		newProfileRoutes(h, t.Profiles, l)
		newCIRAConfigRoutes(h, t.CIRAConfigs, l)
		newWirelessConfigRoutes(h, t.WirelessProfiles, l)
		newIEEE8021xConfigRoutes(h, t.IEEE8021xProfiles, l)
	}

	// Catch-all route to serve index.html for any route not matched above to be handled by Angular
	handler.NoRoute(func(c *gin.Context) {
		c.File("./ui/index.html")
	})
}
