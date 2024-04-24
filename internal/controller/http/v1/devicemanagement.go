package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/alarmclock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type deviceManagementRoutes struct {
	t usecase.DeviceManagement
	d usecase.Device
	l logger.Interface
}

func newAmtRoutes(handler *gin.RouterGroup, t usecase.DeviceManagement, d usecase.Device, l logger.Interface) {
	r := &deviceManagementRoutes{t, d, l}

	h := handler.Group("/amt")
	{
		h.GET("version/:guid", r.getVersion)

		h.GET("features/:guid", r.getFeatures)
		h.POST("features/:guid", r.getVersion)

		h.GET("alarmOccurrences/:guid", r.getAlarmOccurrences)
		h.POST("alarmOccurrences/:guid", r.getVersion)
		h.DELETE("alarmOccurrences/:guid", r.getVersion)

		h.GET("hardwareInfo/:guid", r.getHardwareInfo)
		h.GET("power/state/:guid", r.getPowerState)
		h.POST("power/action/:guid", r.getVersion)
		h.POST("power/bootOptions/:guid", r.getVersion)
		h.GET("power/capabilities/:guid", r.getPowerCapabilities)

		h.GET("log/audit/:guid", r.getVersion)
		h.GET("log/event/:guid", r.getVersion)
		h.GET("generalSettings/:guid", r.getGeneralSettings)

		h.GET("userConsentCode/cancel/:guid", r.cancelUserConsentCode)
		h.GET("userConsentCode/:guid", r.getUserConsentCode)
		h.POST("userConsentCode/:guid", r.sendConsentCode)
	}
}

func (r *deviceManagementRoutes) getVersion(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - getAmtVersion")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.t.SetupWsmanClient(item, true)

	version, err := r.t.GetAMTVersion()
	if err != nil {
		r.l.Error(err, "http - v1 - getAmtVersion")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, version)
}

func (r *deviceManagementRoutes) getFeatures(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - getFeatures")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.t.SetupWsmanClient(item, true)

	features, err := r.t.GetFeatures()
	if err != nil {
		r.l.Error(err, "http - v1 - getFeatures")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, features)
}

func (r *deviceManagementRoutes) getAlarmOccurrences(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - getAlarmOccurrences")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.t.SetupWsmanClient(item, true)

	alarms, err := r.t.GetAlarmOccurrences()
	if err != nil {
		r.l.Error(err, "http - v1 - getAlarmOccurrences")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	if alarms == nil {
		alarms = []alarmclock.AlarmClockOccurrence{}
	}

	c.JSON(http.StatusOK, alarms)
}

func (r *deviceManagementRoutes) getHardwareInfo(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - getHardwareInfo")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.t.SetupWsmanClient(item, true)

	hwInfo, err := r.t.GetHardwareInfo()
	if err != nil {
		r.l.Error(err, "http - v1 - getHardwareInfo")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, hwInfo)
}

func (r *deviceManagementRoutes) getPowerState(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - getPowerState")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.t.SetupWsmanClient(item, true)

	features, err := r.t.GetPowerState()
	if err != nil {
		r.l.Error(err, "http - v1 - getPowerState")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, features)
}

func (r *deviceManagementRoutes) getPowerCapabilities(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - getPowerCapabilities")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.t.SetupWsmanClient(item, true)

	features, err := r.t.GetPowerCapabilities()
	if err != nil {
		r.l.Error(err, "http - v1 - getPowerCapabilities")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, features)
}

func (r *deviceManagementRoutes) getGeneralSettings(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - getGeneralSettings")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.t.SetupWsmanClient(item, true)

	features, err := r.t.GetGeneralSettings()
	if err != nil {
		r.l.Error(err, "http - v1 - getGeneralSettings")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, features)
}

func (r *deviceManagementRoutes) cancelUserConsentCode(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - cancelUserConsentCode")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.t.SetupWsmanClient(item, true)

	features, err := r.t.CancelUserConsent()
	if err != nil {
		r.l.Error(err, "http - v1 - cancelUserConsentCode")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, features)
}

func (r *deviceManagementRoutes) getUserConsentCode(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - getUserConsentCode")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.t.SetupWsmanClient(item, true)

	features, err := r.t.GetUserConsentCode()
	if err != nil {
		r.l.Error(err, "http - v1 - getUserConsentCode")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, features)
}

func (r *deviceManagementRoutes) sendConsentCode(c *gin.Context) {
	var userConsent dto.UserConsent
	if err := c.ShouldBindJSON(&userConsent); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - sendConsentCode")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.t.SetupWsmanClient(item, true)

	features, err := r.t.SendConsentCode(userConsent.ConsentCode)
	if err != nil {
		r.l.Error(err, "http - v1 - sendConsentCode")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, features)
}
