package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type deviceRoutes struct {
	t usecase.Device
	l logger.Interface
}

func newDeviceRoutes(handler *gin.RouterGroup, t usecase.Device, l logger.Interface) {
	r := &deviceRoutes{t, l}

	h := handler.Group("/devices")
	{
		h.GET("", r.get)
		h.GET(":guid", r.get)
		h.POST("", r.insert)
		h.PATCH("", r.update)
		h.DELETE(":guid", r.delete)
	}
}

type DeviceCountResponse struct {
	Count int             `json:"totalAccount"`
	Data  []entity.Device `json:"data"`
}

// @Summary     Show Devices
// @Description Show all devices
// @ID          devices
// @Tags  	    devices
// @Accept      json
// @Produce     json
// @Success     200 {object} DeviceCountResponse
// @Failure     500 {object} response
// @Router      /api/v1/admin/devices [get]
func (dr *deviceRoutes) get(c *gin.Context) {
	var odata OData
	if err := c.ShouldBindQuery(&odata); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	items, err := dr.t.Get(c.Request.Context(), odata.Top, odata.Skip, "")
	if err != nil {
		dr.l.Error(err, "http - devices - v1 - getCount")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	if odata.Count {
		count, err := dr.t.GetCount(c.Request.Context(), "")
		if err != nil {
			dr.l.Error(err, "http - devices - v1 - getCount")
			errorResponse(c, http.StatusInternalServerError, "database problems")

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

func (dr *deviceRoutes) insert(c *gin.Context) {
	var device entity.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	_, err := dr.t.Insert(c.Request.Context(), &device)
	if err != nil {
		dr.l.Error(err, "http - devices - v1 - insert")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusOK, device)
}

func (dr *deviceRoutes) update(c *gin.Context) {
	var device entity.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	updateSuccessful, err := dr.t.Update(c.Request.Context(), &device)
	if err != nil || !updateSuccessful {
		dr.l.Error(err, "http - devices - v1 - update")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusOK, device)
}

func (dr *deviceRoutes) delete(c *gin.Context) {
	var device entity.Device
	if err := c.ShouldBindUri(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	deleteSuccessful, err := dr.t.Delete(c.Request.Context(), device.GUID, "")
	if err != nil {
		dr.l.Error(err, "http - devices - v1 - delete")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusOK, deleteSuccessful)
}
