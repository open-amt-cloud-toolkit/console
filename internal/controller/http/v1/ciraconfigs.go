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

func NewCIRAConfigRoutes(handler *gin.RouterGroup, t ciraconfigs.Feature, l logger.Interface) {
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

// @Summary     CIRA Configurations
// @Description Show all CIRA Configuration profiles
// @ID          getCiraConfigs
// @Tags  	    ciraconfig
// @Accept      json
// @Produce     json
// @Success     200 {object} CIRAConfigCountResponse
// @Failure     500 {object} response
// @Router      /api/v1/admin/ciraconfigs [get]
func (r *ciraConfigRoutes) get(c *gin.Context) {
	var odata OData
	if err := c.ShouldBindQuery(&odata); err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - getCount")
		ErrorResponse(c, err)

		return
	}

	configs, err := r.cira.Get(c.Request.Context(), odata.Top, odata.Skip, "")
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - getCount")
		ErrorResponse(c, err)

		return
	}

	if odata.Count {
		count, err := r.cira.GetCount(c.Request.Context(), "")
		if err != nil {
			r.l.Error(err, "http - CIRA configs - v1 - getCount")
			ErrorResponse(c, err)

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

// @Summary     CIRA Configuration
// @Description Show a CIRA Configuration profile
// @ID          getCiraConfig
// @Tags  	    ciraconfig
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.CIRAConfig
// @Failure     500 {object} response
// @Router      /api/v1/admin/ciraconfigs/:ciraConfigName [get]
func (r *ciraConfigRoutes) getByName(c *gin.Context) {
	configName := c.Param("ciraConfigName")

	foundConfig, err := r.cira.GetByName(c.Request.Context(), configName, "")
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - getByName")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, foundConfig)
}

// @Summary     Add CIRA Configuration
// @Description Add CIRA Configuration
// @ID          addCiraConfig
// @Tags        ciraconfig
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.CIRAConfig
// @Failure     500 {object} response
// @Router      /api/v1/admin/ciraconfigs [post]
func (r *ciraConfigRoutes) insert(c *gin.Context) {
	var config dto.CIRAConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - insert")
		ErrorResponse(c, err)

		return
	}

	newCiraConfig, err := r.cira.Insert(c.Request.Context(), &config)
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - insert")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusCreated, newCiraConfig)
}

// @Summary     Edit CIRA Configuration
// @Description Edit CIRA Configuration
// @ID          updateCiraConfig
// @Tags        ciraconfig
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.CIRAConfig
// @Failure     500 {object} response
// @Router      /api/v1/admin/ciraconfigs [patch]
func (r *ciraConfigRoutes) update(c *gin.Context) {
	var config dto.CIRAConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - update")
		ErrorResponse(c, err)

		return
	}

	updatedConfig, err := r.cira.Update(c.Request.Context(), &config)
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - update")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, updatedConfig)
}

// @Summary     Remove CIRA Configuration
// @Description Remove a CIRA Configuration profile
// @ID          removeCiraConfig
// @Tags  	    ciraconfig
// @Accept      json
// @Produce     json
// @Success     200 {object} nil
// @Failure     500 {object} response
// @Router      /api/v1/admin/ciraconfigs/:ciraConfigName [delete]
func (r *ciraConfigRoutes) delete(c *gin.Context) {
	configName := c.Param("ciraConfigName")

	err := r.cira.Delete(c.Request.Context(), configName, "")
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - delete")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusNoContent, nil)
}
