package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/ciraconfigs"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type ciraConfigRoutes struct {
	cira ciraconfigs.Feature
	l    logger.Interface
}

func newCIRAConfigRoutes(handler *gin.RouterGroup, t ciraconfigs.Feature, l logger.Interface) {
	r := &ciraConfigRoutes{t, l}

	h := handler.Group("/ciraconfigs")
	{
		h.GET("", r.get)
		h.GET(":ciraConfigName", r.getByName)
		h.POST("", r.insert)
		h.PATCH("", r.update)
		h.DELETE(":ciraConfigName", r.delete)
	}
}

type CIRAConfigCountResponse struct {
	Count int              `json:"totalCount"`
	Data  []dto.CIRAConfig `json:"data"`
}

func (r *ciraConfigRoutes) get(c *gin.Context) {
	var odata OData
	if err := c.ShouldBindQuery(&odata); err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - getCount")
		errorResponse(c, err)

		return
	}

	configs, err := r.cira.Get(c.Request.Context(), odata.Top, odata.Skip, "")
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - getCount")
		errorResponse(c, err)

		return
	}

	if odata.Count {
		count, err := r.cira.GetCount(c.Request.Context(), "")
		if err != nil {
			r.l.Error(err, "http - CIRA configs - v1 - getCount")
			errorResponse(c, err)

			return
		}

		countResponse := CIRAConfigCountResponse{
			Count: count,
			Data:  configs,
		}

		c.JSON(http.StatusOK, countResponse)
	} else {
		c.JSON(http.StatusOK, configs)
	}
}

func (r *ciraConfigRoutes) getByName(c *gin.Context) {
	configName := c.Param("ciraConfigName")

	foundConfig, err := r.cira.GetByName(c.Request.Context(), configName, "")
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - getByName")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, foundConfig)
}

func (r *ciraConfigRoutes) insert(c *gin.Context) {
	var config dto.CIRAConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - insert")
		errorResponse(c, err)

		return
	}

	newCiraConfig, err := r.cira.Insert(c.Request.Context(), &config)
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - insert")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusCreated, newCiraConfig)
}

func (r *ciraConfigRoutes) update(c *gin.Context) {
	var config dto.CIRAConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - update")
		errorResponse(c, err)

		return
	}

	updatedConfig, err := r.cira.Update(c.Request.Context(), &config)
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - update")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, updatedConfig)
}

func (r *ciraConfigRoutes) delete(c *gin.Context) {
	configName := c.Param("ciraConfigName")

	err := r.cira.Delete(c.Request.Context(), configName, "")
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - delete")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusNoContent, nil)
}
