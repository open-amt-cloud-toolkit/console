package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"

	v1 "github.com/open-amt-cloud-toolkit/console/internal/controller/http/v1"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type deviceManagementRoutes struct {
	d devices.Feature
	l logger.Interface
}

func NewAmtRoutes(handler *gin.RouterGroup, d devices.Feature, l logger.Interface) {
	r := &deviceManagementRoutes{d, l}

	h := handler.Group("/amt")
	{
		h.GET("version/:guid", r.getVersion)
		h.GET("features/:guid", r.getFeatures)
	}
}

func (r *deviceManagementRoutes) getVersion(c *gin.Context) {
	guid := c.Param("guid")

	_, v2, err := r.d.GetVersion(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v2 - GetVersion")
		v1.ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, v2)
}

func (r *deviceManagementRoutes) getFeatures(c *gin.Context) {
	guid := c.Param("guid")

	features, err := r.d.GetFeatures(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v2 - getFeatures")
		v1.ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, features)
}
