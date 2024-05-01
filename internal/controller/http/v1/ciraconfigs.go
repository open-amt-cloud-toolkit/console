package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/ciraconfigs"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
	"github.com/open-amt-cloud-toolkit/console/pkg/postgres"
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
	Count int                 `json:"totalCount"`
	Data  []entity.CIRAConfig `json:"data"`
}

func (r *ciraConfigRoutes) get(c *gin.Context) {
	var odata OData
	if err := c.ShouldBindQuery(&odata); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	configs, err := r.cira.Get(c.Request.Context(), odata.Top, odata.Skip, "")
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - getCount")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	if odata.Count {
		count, err := r.cira.GetCount(c.Request.Context(), "")
		if err != nil {
			r.l.Error(err, "http - CIRA configs - v1 - getCount")
			errorResponse(c, http.StatusInternalServerError, "database problems")
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
		if strings.Contains(err.Error(), "Not Found") {
			r.l.Error(err, "CIRA Config "+configName+" not found")
			errorResponse(c, http.StatusNotFound, "database problems")
		} else {
			r.l.Error(err, "http - CIRA configs - v1 - getByName")
			errorResponse(c, http.StatusInternalServerError, "database problems")
		}

		return
	}

	c.JSON(http.StatusOK, foundConfig)
}

func (r *ciraConfigRoutes) insert(c *gin.Context) {
	var config entity.CIRAConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	version, err := r.cira.Insert(c.Request.Context(), &config)
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - insert")

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgres.UniqueViolation {
				errorResponse(c, http.StatusBadRequest, pgErr.Message)
			}

			return
		}

		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusCreated, version)
}

func (r *ciraConfigRoutes) update(c *gin.Context) {
	var config entity.CIRAConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	updated, err := r.cira.Update(c.Request.Context(), &config)
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - update")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	if !updated {
		errorResponse(c, http.StatusNotFound, "not found")

		return
	}

	updatedConfig, err := r.cira.GetByName(c, config.ConfigName, config.TenantID)
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - getByName")

		return
	}

	c.JSON(http.StatusOK, updatedConfig)
}

func (r *ciraConfigRoutes) delete(c *gin.Context) {
	configName := c.Param("ciraConfigName")

	deleted, err := r.cira.Delete(c.Request.Context(), configName, "")
	if err != nil {
		r.l.Error(err, "http - CIRA configs - v1 - delete")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	if !deleted {
		errorResponse(c, http.StatusNotFound, "not found")
	}

	c.JSON(http.StatusNoContent, deleted)
}
