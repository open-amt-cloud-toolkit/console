package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
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
		h.GET("diskInfo/:guid", r.getDiskInfo)
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
		h.GET("tls/:guid", r.getTLSSettingData)
	}
}

// @Summary     Get Intel® AMT Version
// @Description Retrieves hardware version information for Intel® AMT and the current activation state.
// @ID          getVersion
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     500 {object} response
// @Router      /api/v1/amt/version/:guid [get]
func (r *deviceManagementRoutes) getVersion(c *gin.Context) {
	guid := c.Param("guid")

	versionv1, _, err := r.d.GetVersion(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - GetVersion")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, versionv1)
}

// @Summary     Get Intel® AMT Features
// @Description Retrieves the current Intel® AMT Enable/Disable state for User Consent, Redirection, KVM, SOL, and IDE-R.
// @Description
// @Description optInState refers to the current Opt In State if the device has User Consent enabled. Valid values:
// @Description
// @Description 0 (Not Started) - No sessions in progress or user consent requested
// @Description 1 (Requested) - Request to AMT device for user consent code successful
// @Description 2 (Displayed) - AMT device displaying user consent code for 300 seconds (5 minutes) before timeout by default
// @Description 3 (Received) - User consent code was entered correctly, a redirection session can be started. Will expire after 120 seconds (2 minutes) and return to State 0 if no active redirection session (State 4)
// @Description 4 (In Session) - Active redirection session in progress
// @ID          getFeatures
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} dto_v1.Features
// @Failure     500 {object} response
// @Router      /api/v1/amt/features/:guid [get]
func (r *deviceManagementRoutes) getFeatures(c *gin.Context) {
	guid := c.Param("guid")

	features, _, err := r.d.GetFeatures(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getFeatures")
		ErrorResponse(c, err)

		return
	}

	v1Features := map[string]interface{}{
		"redirection":  features.Redirection,
		"KVM":          features.EnableKVM,
		"SOL":          features.EnableSOL,
		"IDER":         features.EnableIDER,
		"optInState":   features.OptInState,
		"userConsent":  features.UserConsent,
		"kvmAvailable": features.KVMAvailable,
	}

	c.JSON(http.StatusOK, v1Features)
}

// @Summary     Set Intel® AMT Features
// @Description Retrieves the current Intel® AMT Enable/Disable state for User Consent, Redirection, KVM, SOL, and IDE-R.
// @ID          setFeatures
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.Features
// @Failure     500 {object} response
// @Router      /api/v1/amt/features/:guid [post]
func (r *deviceManagementRoutes) setFeatures(c *gin.Context) {
	guid := c.Param("guid")

	var features dto.Features
	if err := c.ShouldBindJSON(&features); err != nil {
		ErrorResponse(c, err)

		return
	}

	features, _, err := r.d.SetFeatures(c.Request.Context(), guid, features)
	if err != nil {
		r.l.Error(err, "http - v1 - setFeatures")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, features)
}

// @Summary     Get Alarm Clock Occurences
// @Description Retrieves all of the current Alarm Clock occurences for the device
// @ID          getAlarm
// @Tags  	    device management
// @Accept      json
// @Produce     json
// /@Success     200 {object} []alarmclock.AlarmClockOccurrence
// @Failure     500 {object} response
// @Router      /api/v1/amt/alarmOccurrences/:guid [get]
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

// @Summary     Set new Alarm Clock Occurence
// @Description Create a new Alarm Clock occurence to wake device for AMT device.
// @ID          setAlarm
// @Tags  	    device management
// @Accept      json
// @Produce     json
// /@Success     200 {object} alarmclock.AddAlarmOutput
// @Failure     500 {object} response
// @Router      /api/v1/amt/alarmOccurrences/:guid [post]
func (r *deviceManagementRoutes) createAlarmOccurrences(c *gin.Context) {
	guid := c.Param("guid")

	alarm := &dto.AlarmClockOccurrenceInput{}
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

// @Summary     Remove Alarm Clock Occurence
// @Description Delete named Alarm Clock occurence from the device.
// @ID          deleteAlarm
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} nil
// @Failure     500 {object} response
// @Router      /api/v1/amt/alarmOccurrences/:guid [delete]
func (r *deviceManagementRoutes) deleteAlarmOccurrences(c *gin.Context) {
	guid := c.Param("guid")

	alarm := dto.DeleteAlarmOccurrenceRequest{}
	if err := c.ShouldBindJSON(&alarm); err != nil {
		ErrorResponse(c, err)

		return
	}

	err := r.d.DeleteAlarmOccurrences(c.Request.Context(), guid, alarm.Name)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteAlarmOccurrences")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary     Get Hardware Information
// @Description Retrieve hardware information such as processor or storage.
// @ID          getHardwareInfo
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} interface{}
// @Failure     500 {object} response
// @Router      /api/v1/amt/hardwareInfo/:guid [get]
func (r *deviceManagementRoutes) getHardwareInfo(c *gin.Context) {
	guid := c.Param("guid")

	hwInfo, err := r.d.GetHardwareInfo(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getHardwareInfo")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, hwInfo)
}

func (r *deviceManagementRoutes) getDiskInfo(c *gin.Context) {
	guid := c.Param("guid")

	diskInfo, err := r.d.GetDiskInfo(c.Request.Context(), guid)
	if err != nil {
		r.l.Error(err, "http - v1 - getHardwareInfo")
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, diskInfo)
}

// @Summary     Get Power State
// @Description Retrieve current power state of Intel® AMT device, returns a number that maps to a device power state. Possible power state values:
// @Description
// @Description 2 = On - corresponding to ACPI state G0 or S0 or D0
// @Description 3 = Sleep - Light, corresponding to ACPI state G1, S1/S2, or D1
// @Description 4 = Sleep - Deep, corresponding to ACPI state G1, S3, or D2
// @Description 6 = Off - Hard, corresponding to ACPI state G3, S5, or D3
// @Description 7 = Hibernate (Off - Soft), corresponding to ACPI state S4, where the state of the managed element is preserved and will be recovered upon powering on
// @Description 8 = Off - Soft, corresponding to ACPI state G2, S5, or D3
// @Description 9 = Power Cycle (Off-Hard), corresponds to the managed element reaching the ACPI state G3 followed by ACPI state S0
// @Description 13 = Off - Hard Graceful, equivalent to Off Hard but preceded by a request to the managed element to perform an orderly shutdown
// @ID          getPowerState
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     500 {object} response
// @Router      /api/v1/amt/power/state/:guid [get]
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

// @Summary     Get Power Capabilities
// @Description View what OOB power actions are available for that device.
// @ID          getPowerCapabilities
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     500 {object} response
// @Router      /api/v1/amt/power/capabilities/:guid [get]
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

// @Summary     Get General Settings
// @Description Retrieve the Intel® AMT general settings.
// @ID          getGeneralSettings
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} interface{}
// @Failure     500 {object} response
// @Router      /api/v1/amt/generalSettings/:guid [get]
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

// @Summary     Cancel User Consent Code
// @Description Cancel six digit user consent code previously generated on client device
// @ID          cancelUserConsentCode
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} interface{}
// @Failure     500 {object} response
// @Router      /api/v1/amt/userConsentCode/cancel/:guid [get]
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

// @Summary     Get User Consent Code
// @Description If optInState is 0, it will request for a new user consent code
// @ID          getUserConsentCode
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} interface{}
// @Failure     500 {object} response
// @Router      /api/v1/amt/userConsentCode/:guid [get]
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

// @Summary     Send User Consent Code
// @Description Send the user consent code displayed on the client device.
// @ID          sendUserConsentCode
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} optin.SendOptInCode_OUTPUT
// @Failure     500 {object} response
// @Router      /api/v1/amt/userConsentCode/:guid [post]
func (r *deviceManagementRoutes) sendConsentCode(c *gin.Context) {
	guid := c.Param("guid")

	var userConsent dto.UserConsentCode
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

// @Summary     Perform OOB Power Action (1 - 99)
// @Description Perform an OOB power actions numbered 1 thru 99.
// @Description Execute a GET /power/capabilities/{guid} call first to get the list of available power actions. See AMT Power States for ALL potential power actions.
// @ID          sendPowerAction
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} power.PowerActionResponse
// @Failure     500 {object} response
// @Router      /api/v1/amt/power/action/:guid [post]
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

// @Summary     Get Intel® AMT Audit Log
// @Description Returns Intel® AMT Audit Log data in blocks of 10 records for a specified guid. Reference AMT SDK for definition of property return codes.
// @ID          getAuditLog
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.AuditLog
// @Failure     500 {object} response
// @Router      /api/v1/amt/log/audit/:guid [get]
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

// @Summary     Get Intel® AMT Event Log
// @Description Return sensor and hardware event data from the Intel® AMT event log. Reference AMT SDK for definition of property return codes.
// @ID          getEventLog
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.EventLog
// @Failure     500 {object} response
// @Router      /api/v1/amt/log/event/:guid [get]
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

// @Summary     Set OOB Power Action (100+)
// @Description Perform an OOB power actions numbered 100+.
// @Description Execute a GET /power/capabilities/{guid} call first to get the list of available power actions. See AMT Power States for ALL potential power actions.
// @ID          setBootSettings
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} power.PowerActionResponse
// @Failure     500 {object} response
// @Router      /api/v1/amt/power/bootoptions/:guid [post]
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

// @Summary     Get Intel® AMT Network Settings
// @Description Return network settings.
// @ID          getNetworkSettings
// @Tags  	    device management
// @Accept      json
// @Produce     json
// @Success     200 {object} interface{}
// @Failure     500 {object} response
// @Router      /api/v1/amt/networkSettings/:guid [get]
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

func (r *deviceManagementRoutes) getTLSSettingData(c *gin.Context) {
	guid := c.Param("guid")

	tlsSettingData, err := r.d.GetTLSSettingData(c.Request.Context(), guid)
	if err != nil {
		ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, tlsSettingData)
}
