package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/ieee8021xconfigs"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type ieee8021xConfigRoutes struct {
	t ieee8021xconfigs.Feature
	l logger.Interface
}

func newIEEE8021xConfigRoutes(handler *gin.RouterGroup, t ieee8021xconfigs.Feature, l logger.Interface) {
	r := &ieee8021xConfigRoutes{t, l}

	h := handler.Group("/ieee8021xconfigs")
	{
		h.GET("", r.get)
		h.GET(":profileName", r.getByName)
		h.POST("", r.insert)
		h.PATCH("", r.update)
		h.DELETE(":profileName", r.delete)
	}
}

type IEEE8021xConfigCountResponse struct {
	Count int                   `json:"totalCount"`
	Data  []dto.IEEE8021xConfig `json:"data"`
}

func (r *ieee8021xConfigRoutes) get(c *gin.Context) {
	var odata OData
	if err := c.ShouldBindQuery(&odata); err != nil {
		errorResponse(c, err)

		return
	}

	items, err := r.t.Get(c.Request.Context(), odata.Top, odata.Skip, "")
	if err != nil {
		r.l.Error(err, "http - IEEE8021x configs - v1 - getCount")
		errorResponse(c, err)

		return
	}

	if odata.Count {
		count, err := r.t.GetCount(c.Request.Context(), "")
		if err != nil {
			r.l.Error(err, "http - IEEE8021x configs - v1 - getCount")
			errorResponse(c, err)
		}

		countResponse := IEEE8021xConfigCountResponse{
			Count: count,
			Data:  items,
		}

		c.JSON(http.StatusOK, countResponse)
	} else {
		c.JSON(http.StatusOK, items)
	}
}

func (r *ieee8021xConfigRoutes) getByName(c *gin.Context) {
	configName := c.Param("profileName")

	config, err := r.t.GetByName(c.Request.Context(), configName, "")
	if err != nil {
		r.l.Error(err, "http - IEEE8021x configs - v1 - getByName")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, config)
}

func (r *ieee8021xConfigRoutes) insert(c *gin.Context) {
	var config dto.IEEE8021xConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		errorResponse(c, err)

		return
	}

	newConfig, err := r.t.Insert(c.Request.Context(), &config)
	if err != nil {
		r.l.Error(err, "http - IEEE8021x configs - v1 - insert")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusCreated, newConfig)
}

func (r *ieee8021xConfigRoutes) update(c *gin.Context) {
	var config dto.IEEE8021xConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		errorResponse(c, err)

		return
	}

	updatedConfig, err := r.t.Update(c.Request.Context(), &config)
	if err != nil {
		r.l.Error(err, "http - IEEE8021x configs - v1 - update")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, updatedConfig)
}

func (r *ieee8021xConfigRoutes) delete(c *gin.Context) {
	configName := c.Param("profileName")

	err := r.t.Delete(c.Request.Context(), configName, "")
	if err != nil {
		r.l.Error(err, "http - IEEE8021x configs - v1 - delete")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusNoContent, nil)
}
