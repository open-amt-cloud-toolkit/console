package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
	"github.com/open-amt-cloud-toolkit/console/pkg/postgres"
)

type ieee8021xConfigRoutes struct {
	t usecase.IEEE8021xProfile
	l logger.Interface
}

func newIEEE8021xConfigRoutes(handler *gin.RouterGroup, t usecase.IEEE8021xProfile, l logger.Interface) {
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
	Count int                      `json:"totalCount"`
	Data  []entity.IEEE8021xConfig `json:"data"`
}

func (r *ieee8021xConfigRoutes) get(c *gin.Context) {
	var odata OData
	if err := c.ShouldBindQuery(&odata); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	items, err := r.t.Get(c.Request.Context(), odata.Top, odata.Skip, "")
	if err != nil {
		r.l.Error(err, "http - IEEE8021x configs - v1 - getCount")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	if odata.Count {
		count, err := r.t.GetCount(c.Request.Context(), "")
		if err != nil {
			r.l.Error(err, "http - IEEE8021x configs - v1 - getCount")
			errorResponse(c, http.StatusInternalServerError, "database problems")
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
		if strings.Contains(err.Error(), "Not Found") {
			r.l.Error(err, "IEEE8021x Config "+configName+" not found")
			errorResponse(c, http.StatusNotFound, "database problems")
		} else {
			r.l.Error(err, "http - IEEE8021x configs - v1 - getByName")
			errorResponse(c, http.StatusInternalServerError, "database problems")
		}

		return
	}

	c.JSON(http.StatusOK, config)
}

func (r *ieee8021xConfigRoutes) insert(c *gin.Context) {
	var config entity.IEEE8021xConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	version, err := r.t.Insert(c.Request.Context(), &config)
	if err != nil {
		r.l.Error(err, "http - IEEE8021x configs - v1 - insert")

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

func (r *ieee8021xConfigRoutes) update(c *gin.Context) {
	var config entity.IEEE8021xConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	updated, err := r.t.Update(c.Request.Context(), &config)
	if err != nil {
		r.l.Error(err, "http - IEEE8021x configs - v1 - update")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	if !updated {
		errorResponse(c, http.StatusNotFound, "not found")

		return
	}

	updatedConfig, err := r.t.GetByName(c, config.ProfileName, config.TenantID)
	if err != nil {
		r.l.Error(err, "http - IEEE8021x configs - v1 - getByName")

		return
	}

	c.JSON(http.StatusOK, updatedConfig)
}

func (r *ieee8021xConfigRoutes) delete(c *gin.Context) {
	configName := c.Param("profileName")

	configs, err := r.t.Delete(c.Request.Context(), configName, "")
	if err != nil {
		r.l.Error(err, "http - IEEE8021x configs - v1 - delete")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusNoContent, configs)
}
