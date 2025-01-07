package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"

	v1 "github.com/open-amt-cloud-toolkit/console/internal/controller/http/v1"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	_ "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
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
		h.POST("features/:guid", r.setFeatures)
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

// @Summary     Get Intel® AMT Features
// @Description Retrieves the current Intel® AMT Enable/Disable state for User Consent, Redirection, KVM, SOL, and IDE-R.
// @Description
// @Description optInState refers to the current Opt In State if the device has User Consent enabled. Valid values:
// @Description
// @Description 0 (Not Started) - No sessions in progress or user consent requested
// @Description 1 (Requested) - Request to AMT device for user consent code successful
// @Description 2 (Displayed) - AMT device displaying user consent code for 300 seconds (5 minutes) before timeout by default
// @Description 3 (Received) - User consent code was entered correctly, a redirection session can be started. Will expire after 120 seconds (2 minutes) and return to State 0 if no active redirection session (State 4)
// @Description 4 (In Session) - Active redirection session in progress
// @ID          getFeaturesV2
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.Features "notpaging"
// @Success    default {object} dto_v1.Features "200 application/json"
// @Router      /api/v2/amt/features/:guid [get]
func (r *deviceManagementRoutes) getFeatures(c *gin.Context) {
	guid := c.Param("guid")

	_, v2, err := r.d.GetFeatures(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v2 - getFeatures")
		v1.ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, v2)
}

func (r *deviceManagementRoutes) setFeatures(c *gin.Context) {
	guid := c.Param("guid")

	var features dto.Features

	if err := c.ShouldBindJSON(&features); err != nil {
		r.l.Error(err, "http - v2 - setFeatures")
		v1.ErrorResponse(c, err)

		return
	}

	_, v2, err := r.d.SetFeatures(c.Request.Context(), guid, features)
	if err != nil {
		r.l.Error(err, "http - v2 - setFeatures")
		v1.ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, v2)
}
