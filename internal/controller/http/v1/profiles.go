package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/profiles"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type profileRoutes struct {
	t profiles.Feature
	l logger.Interface
}

func newProfileRoutes(handler *gin.RouterGroup, t profiles.Feature, l logger.Interface) {
	r := &profileRoutes{t, l}

	h := handler.Group("/profiles")
	{
		h.GET("", r.get)
		h.GET(":name", r.getByName)
		h.POST("", r.insert)
		h.PATCH("", r.update)
		h.DELETE(":name", r.delete)
	}
}

type ProfileCountResponse struct {
	Count int              `json:"totalCount"`
	Data  []entity.Profile `json:"data"`
}

// @Summary     Show Profiles
// @Description Show all profiles
// @ID          profiles
// @Tags  	    profiles
// @Accept      json
// @Produce     json
// @Success     200 {object} ProfileCountResponse
// @Failure     500 {object} response
// @Router      /api/v1/admin/profiles [get]
func (pr *profileRoutes) get(c *gin.Context) {
	var odata OData
	if err := c.ShouldBindQuery(&odata); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	items, err := pr.t.Get(c.Request.Context(), odata.Top, odata.Skip, "")
	if err != nil {
		pr.l.Error(err, "http - profiles - v1 - getCount")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	if odata.Count {
		count, err := pr.t.GetCount(c.Request.Context(), "")
		if err != nil {
			pr.l.Error(err, "http - profiles - v1 - getCount")
			errorResponse(c, http.StatusInternalServerError, "database problems")
		}

		countResponse := ProfileCountResponse{
			Count: count,
			Data:  items,
		}

		c.JSON(http.StatusOK, countResponse)
	} else {
		c.JSON(http.StatusOK, items)
	}
}

func (pr *profileRoutes) getByName(c *gin.Context) {
	name := c.Param("name")

	item, err := pr.t.GetByName(c.Request.Context(), name, "")
	if err != nil {
		pr.l.Error(err, "http - profiles - v1 - getByName")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusOK, item)
}

func (pr *profileRoutes) insert(c *gin.Context) {
	var profile entity.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	_, err := pr.t.Insert(c.Request.Context(), &profile)
	if err != nil {
		pr.l.Error(err, "http - profiles - v1 - insert")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusOK, profile)
}

func (pr *profileRoutes) update(c *gin.Context) {
	var profile entity.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	updateSuccessful, err := pr.t.Update(c.Request.Context(), &profile)
	if err != nil || !updateSuccessful {
		pr.l.Error(err, "http - profiles - v1 - update")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusOK, profile)
}

func (pr *profileRoutes) delete(c *gin.Context) {
	name := c.Param("name")

	deleteSuccessful, err := pr.t.Delete(c.Request.Context(), name, "")
	if err != nil {
		pr.l.Error(err, "http - profiles - v1 - delete")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusOK, deleteSuccessful)
}
