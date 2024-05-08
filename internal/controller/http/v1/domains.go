package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/domains"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type domainRoutes struct {
	t domains.Feature
	l logger.Interface
}

func newDomainRoutes(handler *gin.RouterGroup, t domains.Feature, l logger.Interface) {
	r := &domainRoutes{t, l}

	h := handler.Group("/domains")
	{
		h.GET("", r.get)
		h.GET(":name", r.getByName)
		h.POST("", r.insert)
		h.PATCH("", r.update)
		h.DELETE(":name", r.delete)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("storageformat", dto.StorageFormatValidation)
		if err != nil {
			l.Error(err, "error registering validation")
		}
	}
}

type DomainCountResponse struct {
	Count int          `json:"totalCount"`
	Data  []dto.Domain `json:"data"`
}

// @Summary     Show Domains
// @Description Show all domains
// @ID          domains
// @Tags  	    domains
// @Accept      json
// @Produce     json
// @Success     200 {object} DomainCountResponse
// @Failure     500 {object} response
// @Router      /api/v1/admin/domains [get]
func (r *domainRoutes) get(c *gin.Context) {
	var odata OData
	if err := c.ShouldBindQuery(&odata); err != nil {
		errorResponse(c, err)

		return
	}

	items, err := r.t.Get(c.Request.Context(), odata.Top, odata.Skip, "")
	if err != nil {
		r.l.Error(err, "http - v1 - getCount")
		errorResponse(c, err)

		return
	}

	if odata.Count {
		count, err := r.t.GetCount(c.Request.Context(), "")
		if err != nil {
			r.l.Error(err, "http - v1 - getCount")
			errorResponse(c, err)
		}

		countResponse := DomainCountResponse{
			Count: count,
			Data:  items,
		}

		c.JSON(http.StatusOK, countResponse)
	} else {
		c.JSON(http.StatusOK, items)
	}
}

// @Summary     Show Domains
// @Description Show domain by name
// @ID          domains
// @Tags  	    domains
// @Accept      json
// @Produce     json
// @Success     200 {object} DomainCountResponse
// @Failure     500 {object} response
// @Router      /api/v1/admin/domains/:name [get]
func (r *domainRoutes) getByName(c *gin.Context) {
	name := c.Param("name")

	item, err := r.t.GetByName(c.Request.Context(), name, "")
	if err != nil {
		r.l.Error(err, "http - v1 - getByName")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary     Add Domain
// @Description Add Domain
// @ID          domains
// @Tags  	    domains
// @Accept      json
// @Produce     json
// @Success     200 {object} DomainResponse
// @Failure     500 {object} response
// @Router      /api/v1/admin/domains [post]
func (r *domainRoutes) insert(c *gin.Context) {
	var domain dto.Domain
	if err := c.ShouldBindJSON(&domain); err != nil {
		errorResponse(c, err)

		return
	}

	newDomain, err := r.t.Insert(c.Request.Context(), &domain)
	if err != nil {
		r.l.Error(err, "http - v1 - insert")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusCreated, newDomain)
}

// @Summary     Edit Domain
// @Description Edit a Domain
// @ID          updateDomain
// @Tags  	    domains
// @Accept      json
// @Produce     json
// @Success     200 {object} DomainResponse
// @Failure     500 {object} response
// @Router      /api/v1/admin/Domains [patch]
func (r *domainRoutes) update(c *gin.Context) {
	var domain dto.Domain
	if err := c.ShouldBindJSON(&domain); err != nil {
		errorResponse(c, err)

		return
	}

	updatedDomain, err := r.t.Update(c.Request.Context(), &domain)
	if err != nil {
		r.l.Error(err, "http - v1 - update")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, updatedDomain)
}

// @Summary     Remove Domains
// @Description Remove a Domain
// @ID          deleteDomain
// @Tags  	    domains
// @Accept      json
// @Produce     json
// @Success     204 {object} noContent
// @Failure     500 {object} response
// @Router      /api/v1/admin/domains [delete]
func (r *domainRoutes) delete(c *gin.Context) {
	name := c.Param("name")

	err := r.t.Delete(c.Request.Context(), name, "")
	if err != nil {
		r.l.Error(err, "http - v1 - delete")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusNoContent, nil)
}
