package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type deviceRoutes struct {
	t devices.Feature
	l logger.Interface
}

func newDeviceRoutes(handler *gin.RouterGroup, t devices.Feature, l logger.Interface) {
	r := &deviceRoutes{t, l}

	h := handler.Group("/devices")
	{
		h.GET("", r.get)
		h.GET("stats", r.getStats)
		h.GET("redirectstatus/:guid", r.redirectStatus)
		h.GET(":guid", r.get)
		h.GET("tags", r.getTags)
		h.POST("", r.insert)
		h.PATCH("", r.update)
		h.DELETE(":guid", r.delete)
	}
}

type DeviceCountResponse struct {
	Count int          `json:"totalCount"`
	Data  []dto.Device `json:"data"`
}
type DeviceStatResponse struct {
	TotalCount        int `json:"totalCount"`
	ConnectedCount    int `json:"connectedCount"`
	DisconnectedCount int `json:"disconnectedCount"`
}

// @Summary     Gets Device Count
// @Description Gets number of devices
// @ID          getStats
// @Tags  	    devices
// @Accept      json
// @Produce     json
// @Success     200 {object} DeviceCountResponse
// @Failure     500 {object} response
// @Router      /api/v1/admin/devices [get]
func (dr *deviceRoutes) getStats(c *gin.Context) {
	count, err := dr.t.GetCount(c.Request.Context(), "")
	if err != nil {
		dr.l.Error(err, "http - devices - v1 - getCount")
		errorResponse(c, err)

		return
	}

	countResponse := DeviceStatResponse{
		TotalCount: count,
	}

	c.JSON(http.StatusOK, countResponse)
}

// @Summary     Show Devices
// @Description Show all devices
// @ID          getDevices
// @Tags  	    devices
// @Accept      json
// @Produce     json
// @Success     200 {object} DeviceCountResponse
// @Failure     500 {object} response
// @Router      /api/v1/admin/devices [get]
func (dr *deviceRoutes) get(c *gin.Context) {
	var odata OData
	if err := c.ShouldBindQuery(&odata); err != nil {
		errorResponse(c, err)

		return
	}

	tags := c.Query("tags")

	var items []dto.Device

	var err error

	if tags != "" {
		items, err = dr.t.GetByTags(c.Request.Context(), tags, c.Query("method"), odata.Top, odata.Skip, "")
		if err != nil {
			dr.l.Error(err, "http - devices - v1 - get")
			errorResponse(c, err)

			return
		}
	} else {
		items, err = dr.t.Get(c.Request.Context(), odata.Top, odata.Skip, "")
		if err != nil {
			dr.l.Error(err, "http - devices - v1 - get")
			errorResponse(c, err)

			return
		}
	}

	if odata.Count {
		count, err := dr.t.GetCount(c.Request.Context(), "")
		if err != nil {
			dr.l.Error(err, "http - devices - v1 - get")
			errorResponse(c, err)

			return
		}

		countResponse := DeviceCountResponse{
			Count: count,
			Data:  items,
		}

		c.JSON(http.StatusOK, countResponse)
	} else {
		c.JSON(http.StatusOK, items)
	}
}

// @Summary     Add Devices
// @Description Add a devices
// @ID          insertDevice
// @Tags  	    devices
// @Accept      json
// @Produce     json
// @Success     200 {object} DeviceResponse
// @Failure     500 {object} response
// @Router      /api/v1/admin/devices [post]
func (dr *deviceRoutes) insert(c *gin.Context) {
	var device dto.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		errorResponse(c, err)

		return
	}

	newDevice, err := dr.t.Insert(c.Request.Context(), &device)
	if err != nil {
		dr.l.Error(err, "http - devices - v1 - insert")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, newDevice)
}

// @Summary     Edit Devices
// @Description Edit a devices
// @ID          updateDevice
// @Tags  	    devices
// @Accept      json
// @Produce     json
// @Success     200 {object} DeviceResponse
// @Failure     500 {object} response
// @Router      /api/v1/admin/devices [patch]
func (dr *deviceRoutes) update(c *gin.Context) {
	var device dto.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		errorResponse(c, err)

		return
	}

	updatedDevice, err := dr.t.Update(c.Request.Context(), &device)
	if err != nil {
		dr.l.Error(err, "http - devices - v1 - update")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, updatedDevice)
}

// @Summary     Remove Devices
// @Description Remove a device
// @ID          deleteDevice
// @Tags  	    devices
// @Accept      json
// @Produce     json
// @Success     204 {object} noContent
// @Failure     500 {object} response
// @Router      /api/v1/admin/devices [delete]
func (dr *deviceRoutes) delete(c *gin.Context) {
	guid := c.Param("guid")

	err := dr.t.Delete(c.Request.Context(), guid, "")
	if err != nil {
		dr.l.Error(err, "http - devices - v1 - delete")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (dr *deviceRoutes) redirectStatus(c *gin.Context) {
	_ = c.Param("guid")
	result := map[string]bool{
		"isSOLConnected":  false, // device.solConnect,
		"isIDERConnected": false, // device.iderConnect,
	}
	c.JSON(http.StatusOK, result)
}

// @Summary     Get Tags
// @Description Get Available Distinct Tags in the system
// @ID          getTags
// @Tags  	    devices
// @Accept      json
// @Produce     json
// @Success     204 {object} noContent
// @Failure     500 {object} response
// @Router      /api/v1/admin/devices/tags [get]
func (dr *deviceRoutes) getTags(c *gin.Context) {
	tags, err := dr.t.GetDistinctTags(c.Request.Context(), "")
	if err != nil {
		dr.l.Error(err, "http - devices - v1 - tags")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, tags)
}
