package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/wificonfigs"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

var ErrValidationWifiConfig = dto.NotValidError{Console: consoleerrors.CreateConsoleError("WifiConfigsAPI")}

type WirelessConfigRoutes struct {
	t wificonfigs.Feature
	l logger.Interface
}

func NewWirelessConfigRoutes(handler *gin.RouterGroup, t wificonfigs.Feature, l logger.Interface) {
	r := &WirelessConfigRoutes{t, l}

	if binding.Validator != nil {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			err := v.RegisterValidation("authforieee8021x", dto.ValidateAuthandIEEE)
			if err != nil {
				l.Error(err, "failed to register validation")
			}
		}
	}

	h := handler.Group("/wirelessconfigs")
	{
		h.GET("", r.get)
		h.GET(":profileName", r.getByName)
		h.POST("", r.insert)
		h.PATCH("", r.update)
		h.DELETE(":profileName", r.delete)
	}
}

// @Summary     Get All Wireless Configs
// @Description Retrieves all of the Wireless configuration profiles from the database. Will not return the password field to protect the privacy of this asset. 
// @ID          getWifiConfigs
// @Tags  	    wireless
// @Accept      json
// @Produce     json
// @success     200 {object} dto.WirelessConfigCountResponse
// @Failure     500 {object} response
// @Router      /api/v1/admin/wireless/ [get]
func (r *WirelessConfigRoutes) get(c *gin.Context) {
	var odata OData
	if err := c.ShouldBindQuery(&odata); err != nil {
		validationErr := ErrValidationWifiConfig.Wrap("get", "ShouldBindQuery", err)
		ErrorResponse(c, validationErr)

		return
	}

	items, err := r.t.Get(c.Request.Context(), odata.Top, odata.Skip, "")
	if err != nil {
		r.l.Error(err, "http - wireless configs - v1 - getCount")
		ErrorResponse(c, err)

		return
	}

	if odata.Count {
		count, err := r.t.GetCount(c.Request.Context(), "")
		if err != nil {
			r.l.Error(err, "http - wireless configs - v1 - getCount")
			ErrorResponse(c, err)
		}

		countResponse := dto.WirelessConfigCountResponse{
			Count: count,
			Data:  items,
		}

		c.JSON(http.StatusOK, countResponse)
	} else {
		c.JSON(http.StatusOK, items)
	}
}

// @Summary     Get a Wireless Config
// @Description Retrieves the specific Wireless configuration profile from the database. Will not return the password field to protect the privacy of this asset. 
// @ID          getWifiConfig
// @Tags  	    wireless
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.WirelessConfig
// @Failure     500 {object} response
// @Router      /api/v1/admin/wireless/:profileName [get]
func (r *WirelessConfigRoutes) getByName(c *gin.Context) {
	profileName := c.Param("profileName")

	config, err := r.t.GetByName(c.Request.Context(), profileName, "")
	if err != nil {
		r.l.Error(err, "http - wireless configs - v1 - getByName")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, config)
}

// @Summary     Create a Wireless Config
// @Description Creates a new Wireless configuration profile. The PSK passphrase is stored in a secrets manager and is only used during configuration to set the Wi-Fi credentials in the AMT device. 
// @ID          addWifiConfig
// @Tags  	    wireless
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.WirelessConfig
// @Failure     500 {object} response
// @Router      /api/v1/admin/wireless/ [post]
func (r *WirelessConfigRoutes) insert(c *gin.Context) {
	var config dto.WirelessConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		validationErr := ErrValidationWifiConfig.Wrap("insert", "ShouldBindJSON", err)
		ErrorResponse(c, validationErr)

		return
	}

	insertedConfig, err := r.t.Insert(c.Request.Context(), &config)
	if err != nil {
		r.l.Error(err, "http - wireless configs - v1 - insert")

		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusCreated, insertedConfig)
}

// @Summary     Edit a Wireless Config
// @Description Edits an existing Wireless configuration profile.
// @Description 
// @Description The profileName can not be changed.
// @Description 
// @Description Version must be provided to ensure the correct profile is edited.
// @ID          editWifiConfig
// @Tags  	    wireless
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.WirelessConfig
// @Failure     500 {object} response
// @Router      /api/v1/admin/wireless/ [patch]
func (r *WirelessConfigRoutes) update(c *gin.Context) {
	var config dto.WirelessConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		validationErr := ErrValidationWifiConfig.Wrap("update", "ShouldBindJSON", err)
		ErrorResponse(c, validationErr)

		return
	}

	updatedWirelessConfig, err := r.t.Update(c.Request.Context(), &config)
	if err != nil {
		r.l.Error(err, "http - wireless configs - v1 - update")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, updatedWirelessConfig)
}

// @Summary     Remove a Wireless Config
// @Description Removes the specific Wireless configuration profile from the database. 
// @ID          deleteWifiConfig
// @Tags  	    wireless
// @Accept      json
// @Produce     json
// @Success     200 {object} nil
// @Failure     500 {object} response
// @Router      /api/v1/admin/wireless/:profileName [delete]
func (r *WirelessConfigRoutes) delete(c *gin.Context) {
	configName := c.Param("profileName")

	err := r.t.Delete(c.Request.Context(), configName, "")
	if err != nil {
		r.l.Error(err, "http - wireless configs - v1 - delete")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusNoContent, nil)
}
