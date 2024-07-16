package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/amtexplorer"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type deviceManagementRoutes struct {
	d devices.Feature
	a amtexplorer.Feature
	l logger.Interface
}

func NewAmtRoutes(handler *gin.RouterGroup, d devices.Feature, amt amtexplorer.Feature, l logger.Interface) {
	r := &deviceManagementRoutes{d, amt, l}

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

		h.GET("explorer", r.getCallList)
		h.GET("explorer/:guid/:call", r.executeCall)
		h.GET("certificates/:guid", r.getCertificates)
	}
}

func (r *deviceManagementRoutes) getVersion(c *gin.Context) {
	guid := c.Param("guid")

	version, err := r.d.GetVersion(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - GetVersion")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, version)
}

func (r *deviceManagementRoutes) getFeatures(c *gin.Context) {
	guid := c.Param("guid")

	features, err := r.d.GetFeatures(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getFeatures")
		ErrorResponse(c, err)

		return
	}

	v1Features := map[string]interface{}{
		"redirection": features.Redirection,
		"KVM":         features.EnableKVM,
		"SOL":         features.EnableSOL,
		"IDER":        features.EnableIDER,
		"optInState":  features.OptInState,
		"userConsent": features.UserConsent,
	}

	c.JSON(http.StatusOK, v1Features)
}

func (r *deviceManagementRoutes) setFeatures(c *gin.Context) {
	guid := c.Param("guid")

	var features dto.Features

	if err := c.ShouldBindJSON(&features); err != nil {
		ErrorResponse(c, err)

		return
	}

	features, err := r.d.SetFeatures(c.Request.Context(), guid, features)
	if err != nil {
		r.l.Error(err, "http - v1 - setFeatures")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, features)
}

func (r *deviceManagementRoutes) getAlarmOccurrences(c *gin.Context) {
	guid := c.Param("guid")

	alarms, err := r.d.GetAlarmOccurrences(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getFeatures")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, alarms)
}

func (r *deviceManagementRoutes) createAlarmOccurrences(c *gin.Context) {
	guid := c.Param("guid")

	alarm := &dto.AlarmClockOccurrence{}
	if err := c.ShouldBindJSON(alarm); err != nil {
		ErrorResponse(c, err)

		return
	}

	alarmReference, err := r.d.CreateAlarmOccurrences(c.Request.Context(), guid, *alarm)
	if err != nil {
		r.l.Error(err, "http - v1 - createAlarmOccurrences")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusCreated, alarmReference)
}

func (r *deviceManagementRoutes) deleteAlarmOccurrences(c *gin.Context) {
	guid := c.Param("guid")

	alarm := dto.DeleteAlarmOccurrenceRequest{}
	if err := c.ShouldBindJSON(&alarm); err != nil {
		ErrorResponse(c, err)

		return
	}

	if alarm.InstanceID == nil {
		alarm.InstanceID = new(string)
	}

	err := r.d.DeleteAlarmOccurrences(c.Request.Context(), guid, *alarm.InstanceID)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteAlarmOccurrences")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (r *deviceManagementRoutes) getHardwareInfo(c *gin.Context) {
	guid := c.Param("guid")

	hwInfo, err := r.d.GetHardwareInfo(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getHardwareInfo")
		ErrorResponse(c, err)

		return
	}

	v1HwInfo := map[string]interface{}{
		"CIM_ComputerSystemPackage": map[string]interface{}{
			"response":  hwInfo.ComputerSystemPackage.PlatformGUID,
			"responses": hwInfo.ComputerSystemPackage.PlatformGUID,
		},
		"CIM_SystemPackaging": map[string]interface{}{
			"responses": []interface{}{hwInfo.SystemPackage},
		},
		"CIM_Chassis": map[string]interface{}{
			"response":  hwInfo.Chassis,
			"responses": []interface{}{},
		}, "CIM_Chip": map[string]interface{}{
			"responses": []interface{}{hwInfo.Chip},
		}, "CIM_Card": map[string]interface{}{
			"response":  hwInfo.Card,
			"responses": []interface{}{},
		}, "CIM_BIOSElement": map[string]interface{}{
			"response":  hwInfo.BIOSElement,
			"responses": []interface{}{},
		}, "CIM_Processor": map[string]interface{}{
			"responses": []interface{}{hwInfo.Processor},
		}, "CIM_PhysicalMemory": map[string]interface{}{
			"responses": hwInfo.PhysicalMemory,
		}, "CIM_MediaAccessDevice": map[string]interface{}{
			"responses": []interface{}{hwInfo.MediaAccessDevices},
		}, "CIM_PhysicalPackage": map[string]interface{}{
			"responses": []interface{}{hwInfo.PhysicalPackage},
		},
	}

	// v1HwInfo := map[string]interface{}{
	// 	"CIM_ComputerSystemPackage": map[string]interface{}{
	// 		"response":  hwInfo.CSPResult.Body.GetResponse,
	// 		"responses": hwInfo.CSPResult.Body.GetResponse,
	// 	},
	// 	"CIM_SystemPackaging": map[string]interface{}{
	// 		"responses": []interface{}{hwInfo.SPPullResult.Body.PullResponse.SystemPackageItems},
	// 	},
	// 	"CIM_Chassis": map[string]interface{}{
	// 		"response":  hwInfo.ChassisResult.Body.PackageResponse,
	// 		"responses": []interface{}{},
	// 	}, "CIM_Chip": map[string]interface{}{
	// 		"responses": []interface{}{hwInfo.ChipResult.Body.PackageResponse},
	// 	}, "CIM_Card": map[string]interface{}{
	// 		"response":  hwInfo.CardResult.Body.PackageResponse,
	// 		"responses": []interface{}{},
	// 	}, "CIM_BIOSElement": map[string]interface{}{
	// 		"response":  hwInfo.BiosResult.Body.GetResponse,
	// 		"responses": []interface{}{},
	// 	}, "CIM_Processor": map[string]interface{}{
	// 		"responses": []interface{}{hwInfo.ProcessorResult.Body.PackageResponse},
	// 	}, "CIM_PhysicalMemory": map[string]interface{}{
	// 		"responses": hwInfo.PhysicalMemoryResult.Body.PullResponse.MemoryItems,
	// 	}, "CIM_MediaAccessDevice": map[string]interface{}{
	// 		"responses": []interface{}{hwInfo.MediaAccessPullResult.Body.PullResponse.MediaAccessDevices},
	// 	}, "CIM_PhysicalPackage": map[string]interface{}{
	// 		"responses": []interface{}{hwInfo.PPPullResult.Body.PullResponse.PhysicalPackage},
	// 	},
	// }

	c.JSON(http.StatusOK, v1HwInfo)
}

func (r *deviceManagementRoutes) getPowerState(c *gin.Context) {
	guid := c.Param("guid")

	state, err := r.d.GetPowerState(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getPowerState")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, state)
}

func (r *deviceManagementRoutes) getPowerCapabilities(c *gin.Context) {
	guid := c.Param("guid")

	power, err := r.d.GetPowerCapabilities(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getPowerCapabilities")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, power)
}

func (r *deviceManagementRoutes) getGeneralSettings(c *gin.Context) {
	guid := c.Param("guid")

	generalSettings, err := r.d.GetGeneralSettings(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getGeneralSettings")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, generalSettings)
}

func (r *deviceManagementRoutes) cancelUserConsentCode(c *gin.Context) {
	guid := c.Param("guid")

	result, err := r.d.CancelUserConsent(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - cancelUserConsentCode")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, result)
}

func (r *deviceManagementRoutes) getUserConsentCode(c *gin.Context) {
	guid := c.Param("guid")

	response, err := r.d.GetUserConsentCode(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getUserConsentCode")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, response)
}

func (r *deviceManagementRoutes) sendConsentCode(c *gin.Context) {
	guid := c.Param("guid")

	var userConsent dto.UserConsent
	if err := c.ShouldBindJSON(&userConsent); err != nil {
		ErrorResponse(c, err)

		return
	}

	response, err := r.d.SendConsentCode(c.Request.Context(), userConsent, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - sendConsentCode")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, response)
}

func (r *deviceManagementRoutes) powerAction(c *gin.Context) {
	guid := c.Param("guid")

	var powerAction dto.PowerAction
	if err := c.ShouldBindJSON(&powerAction); err != nil {
		ErrorResponse(c, err)

		return
	}

	response, err := r.d.SendPowerAction(c.Request.Context(), guid, powerAction.Action)
	if err != nil {
		r.l.Error(err, "http - v1 - powerAction")
		ErrorResponse(c, err)

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
		ErrorResponse(c, err)

		return
	}

	auditLogs, err := r.d.GetAuditLog(c.Request.Context(), startIdx, guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getAuditLog")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, auditLogs)
}

func (r *deviceManagementRoutes) getEventLog(c *gin.Context) {
	guid := c.Param("guid")

	eventLogs, err := r.d.GetEventLog(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getEventLog")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, eventLogs)
}

func (r *deviceManagementRoutes) setBootOptions(c *gin.Context) {
	guid := c.Param("guid")

	var bootSetting dto.BootSetting
	if err := c.ShouldBindJSON(&bootSetting); err != nil {
		ErrorResponse(c, err)

		return
	}

	features, err := r.d.SetBootOptions(c.Request.Context(), guid, bootSetting)
	if err != nil {
		r.l.Error(err, "http - v1 - setBootOptions")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, features)
}

func (r *deviceManagementRoutes) getNetworkSettings(c *gin.Context) {
	guid := c.Param("guid")

	network, err := r.d.GetNetworkSettings(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getNetworkSettings")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, network)
}

// @Summary     Get Call List
// @Description Get a list of supported WSMAN calls
// @ID          getCallList
// @Tags  	    devices
// @Accept      json
// @Produce     json
// @Success     200 {object} DeviceCountResponse
// @Failure     500 {object} response
// @Router      /api/v1/devices [get]
func (r *deviceManagementRoutes) getCallList(c *gin.Context) {
	items := r.a.GetExplorerSupportedCalls()

	c.JSON(http.StatusOK, items)
}

// @Summary     Execute Call
// @Description Execute a call
// @ID          executeCall
// @Tags  	    amt
// @Accept      json
// @Produce     json
// @Success     200 {object} DeviceCountResponse
// @Failure     500 {object} response
// @Router      /api/v1/devices [get]
func (r *deviceManagementRoutes) executeCall(c *gin.Context) {
	guid := c.Param("guid")
	call := c.Param("call")

	result, err := r.a.ExecuteCall(c.Request.Context(), guid, call, "")
	if err != nil {
		r.l.Error(err, "http - explorer - v1 - executeCall")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, result)
}

func (r *deviceManagementRoutes) getCertificates(c *gin.Context) {
	guid := c.Param("guid")

	certs, err := r.d.GetCertificates(c.Request.Context(), guid)
	if err != nil {
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, certs)
}
