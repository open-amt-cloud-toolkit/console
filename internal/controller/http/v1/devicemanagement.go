package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type deviceManagementRoutes struct {
	d devices.Feature
	l logger.Interface
}

func newAmtRoutes(handler *gin.RouterGroup, d devices.Feature, l logger.Interface) {
	r := &deviceManagementRoutes{d, l}

	h := handler.Group("/amt")
	{
		h.GET("version/:guid", r.getVersion)

		h.GET("features/:guid", r.getFeatures)
		h.POST("features/:guid", r.setFeatures)

		h.GET("alarmOccurrences/:guid", r.getAlarmOccurrences)
		h.POST("alarmOccurrences/:guid", r.createAlarmOccurrences)
		h.DELETE("alarmOccurrences/:guid", r.deleteAlarmOccurrences)

		h.GET("hardwareInfo/:guid", r.getHardwareInfo)
		h.GET("power/state/:guid", r.getPowerState)
		h.POST("power/action/:guid", r.powerAction)
		h.POST("power/bootOptions/:guid", r.setBootOptions)
		h.POST("power/bootoptions/:guid", r.setBootOptions)
		h.GET("power/capabilities/:guid", r.getPowerCapabilities)

		h.GET("log/audit/:guid", r.getAuditLog)
		h.GET("log/event/:guid", r.getEventLog)
		h.GET("generalSettings/:guid", r.getGeneralSettings)

		h.GET("userConsentCode/cancel/:guid", r.cancelUserConsentCode)
		h.GET("userConsentCode/:guid", r.getUserConsentCode)
		h.POST("userConsentCode/:guid", r.sendConsentCode)

		h.GET("networkSettings/:guid", r.getNetworkSettings)
	}
}

func (r *deviceManagementRoutes) getVersion(c *gin.Context) {
	guid := c.Param("guid")

	version, err := r.d.GetVersion(c, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - GetVersion")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, version)
}

func (r *deviceManagementRoutes) getFeatures(c *gin.Context) {
	guid := c.Param("guid")

	features, err := r.d.GetFeatures(c, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getFeatures")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, features)
}

func (r *deviceManagementRoutes) setFeatures(c *gin.Context) {
	guid := c.Param("guid")

	var features dto.Features

	if err := c.ShouldBindJSON(&features); err != nil {
		errorResponse(c, err)

		return
	}

	features, err := r.d.SetFeatures(c, guid, features)
	if err != nil {
		r.l.Error(err, "http - v1 - setFeatures")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, features)
}

func (r *deviceManagementRoutes) getAlarmOccurrences(c *gin.Context) {
	guid := c.Param("guid")

	alarms, err := r.d.GetAlarmOccurrences(c, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getFeatures")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, alarms)
}

func (r *deviceManagementRoutes) createAlarmOccurrences(c *gin.Context) {
	guid := c.Param("guid")

	alarm := dto.AlarmClockOccurrence{}
	if err := c.ShouldBindJSON(&alarm); err != nil {
		errorResponse(c, err)

		return
	}

	alarmReference, err := r.d.CreateAlarmOccurrences(c, guid, alarm)
	if err != nil {
		r.l.Error(err, "http - v1 - createAlarmOccurrences")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, alarmReference)
}

func (r *deviceManagementRoutes) deleteAlarmOccurrences(c *gin.Context) {
	guid := c.Param("guid")

	alarm := dto.AlarmClockOccurrence{}
	if err := c.ShouldBindJSON(&alarm); err != nil {
		errorResponse(c, err)

		return
	}

	err := r.d.DeleteAlarmOccurrences(c, guid, alarm.InstanceID)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteAlarmOccurrences")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (r *deviceManagementRoutes) getHardwareInfo(c *gin.Context) {
	guid := c.Param("guid")

	hwInfo, err := r.d.GetHardwareInfo(c, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getHardwareInfo")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, hwInfo)
}

func (r *deviceManagementRoutes) getPowerState(c *gin.Context) {
	guid := c.Param("guid")

	state, err := r.d.GetPowerState(c, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getPowerState")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, state)
}

func (r *deviceManagementRoutes) getPowerCapabilities(c *gin.Context) {
	guid := c.Param("guid")

	power, err := r.d.GetPowerCapabilities(c, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getPowerCapabilities")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, power)
}

func (r *deviceManagementRoutes) getGeneralSettings(c *gin.Context) {
	guid := c.Param("guid")

	generalSettings, err := r.d.GetGeneralSettings(c, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getGeneralSettings")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, generalSettings)
}

func (r *deviceManagementRoutes) cancelUserConsentCode(c *gin.Context) {
	guid := c.Param("guid")

	result, err := r.d.CancelUserConsent(c, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - cancelUserConsentCode")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, result)
}

func (r *deviceManagementRoutes) getUserConsentCode(c *gin.Context) {
	guid := c.Param("guid")

	response, err := r.d.GetUserConsentCode(c, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getUserConsentCode")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, response)
}

func (r *deviceManagementRoutes) sendConsentCode(c *gin.Context) {
	guid := c.Param("guid")

	var userConsent dto.UserConsent
	if err := c.ShouldBindJSON(&userConsent); err != nil {
		errorResponse(c, err)

		return
	}

	response, err := r.d.SendConsentCode(c, userConsent, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - sendConsentCode")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, response)
}

func (r *deviceManagementRoutes) powerAction(c *gin.Context) {
	guid := c.Param("guid")

	var powerAction dto.PowerAction
	if err := c.ShouldBindJSON(&powerAction); err != nil {
		errorResponse(c, err)

		return
	}

	response, err := r.d.SendPowerAction(c, guid, powerAction.Action)
	if err != nil {
		r.l.Error(err, "http - v1 - powerAction")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, response)
}

func (r *deviceManagementRoutes) getAuditLog(c *gin.Context) {
	guid := c.Param("guid")

	startIndex := c.Query("startIndex")

	startIdx, err := strconv.Atoi(startIndex)
	if err != nil {
		r.l.Error(err, "http - v1 - getAuditLog")
		errorResponse(c, err)

		return
	}

	auditLogs, err := r.d.GetAuditLog(c, startIdx, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getAuditLog")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, auditLogs)
}

func (r *deviceManagementRoutes) getEventLog(c *gin.Context) {
	guid := c.Param("guid")

	eventLogs, err := r.d.GetEventLog(c, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getEventLog")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, eventLogs)
}

func (r *deviceManagementRoutes) setBootOptions(c *gin.Context) {
	guid := c.Param("guid")

	var bootSetting dto.BootSetting
	if err := c.ShouldBindJSON(&bootSetting); err != nil {
		errorResponse(c, err)

		return
	}

	features, err := r.d.SetBootOptions(c, guid, bootSetting)
	if err != nil {
		r.l.Error(err, "http - v1 - setBootOptions")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, features)
}

func (r *deviceManagementRoutes) getNetworkSettings(c *gin.Context) {
	guid := c.Param("guid")

	network, err := r.d.GetNetworkSettings(c, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getNetworkSettings")
		errorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, network)
}
