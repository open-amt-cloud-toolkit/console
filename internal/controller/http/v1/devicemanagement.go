package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/boot"
	cimBoot "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/boot"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/power"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/alarmclock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type deviceManagementRoutes struct {
	dm usecase.DeviceManagement
	d  usecase.Device
	l  logger.Interface
}

func newAmtRoutes(handler *gin.RouterGroup, dm usecase.DeviceManagement, d usecase.Device, l logger.Interface) {
	r := &deviceManagementRoutes{dm, d, l}

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
		h.GET("power/capabilities/:guid", r.getPowerCapabilities)

		h.GET("log/audit/:guid", r.getAuditLog)
		h.GET("log/event/:guid", r.getEventLog)
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

	r.dm.SetupWsmanClient(item, true)

	version, err := r.dm.GetAMTVersion()
	if err != nil {
		r.l.Error(err, "http - v1 - getAmtVersion")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	data, err := r.dm.GetSetupAndConfiguration()
	if err != nil {
		r.l.Error(err, "http - v1 - getAmtVersion")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	response := map[string]interface{}{
		"CIM_SoftwareIdentity": map[string]interface{}{
			"responses": version,
		},
		"AMT_SetupAndConfigurationService": map[string]interface{}{
			"response": data[0],
		},
	}

	c.JSON(http.StatusOK, response)
}

func (r *deviceManagementRoutes) getFeatures(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - getFeatures")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.dm.SetupWsmanClient(item, true)

	features, err := r.dm.GetFeatures()
	if err != nil {
		r.l.Error(err, "http - v1 - getFeatures")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, features)
}

func (r *deviceManagementRoutes) setFeatures(c *gin.Context) {
	var features dto.Features
	if err := c.ShouldBindJSON(&features); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - setFeatures")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.dm.SetupWsmanClient(item, true)

	features, err = r.dm.SetFeatures(features)
	if err != nil {
		r.l.Error(err, "http - v1 - setFeatures")
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

	r.dm.SetupWsmanClient(item, true)

	alarms, err := r.dm.GetAlarmOccurrences()
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

func (r *deviceManagementRoutes) createAlarmOccurrences(c *gin.Context) {
	alarm := dto.AlarmClockOccurrence{}
	if err := c.ShouldBindJSON(&alarm); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - createAlarmOccurrences")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.dm.SetupWsmanClient(item, true)

	alarmReference, err := r.dm.CreateAlarmOccurrences(alarm.InstanceID, alarm.StartTime, alarm.Interval, alarm.DeleteOnCompletion)
	if err != nil {
		r.l.Error(err, "http - v1 - createAlarmOccurrences")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, alarmReference)
}

func (r *deviceManagementRoutes) deleteAlarmOccurrences(c *gin.Context) {
	alarm := dto.AlarmClockOccurrence{}
	if err := c.ShouldBindJSON(&alarm); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - deleteAlarmOccurrences")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.dm.SetupWsmanClient(item, true)

	err = r.dm.DeleteAlarmOccurrences(alarm.InstanceID)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteAlarmOccurrences")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (r *deviceManagementRoutes) getHardwareInfo(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - getHardwareInfo")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.dm.SetupWsmanClient(item, true)

	hwInfo, err := r.dm.GetHardwareInfo()
	if err != nil {
		r.l.Error(err, "http - v1 - getHardwareInfo")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, hwInfo)
}

func (r *deviceManagementRoutes) getPowerState(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil {
		r.l.Error(err, "http - v1 - getPowerState")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	if item.GUID == "" {
		errorResponse(c, http.StatusNotFound, "device not found")

		return
	}

	r.dm.SetupWsmanClient(item, true)

	features, err := r.dm.GetPowerState()
	if err != nil {
		r.l.Error(err, "http - v1 - getPowerState")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, features)
}

var MinAMTVersion = 9

func (r *deviceManagementRoutes) getPowerCapabilities(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil {
		r.l.Error(err, "http - v1 - getPowerCapabilities")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	if item.GUID == "" {
		errorResponse(c, http.StatusNotFound, "device not found")

		return
	}

	r.dm.SetupWsmanClient(item, true)

	version, err := r.dm.GetAMTVersion()
	if err != nil {
		r.l.Error(err, "http - v1 - getPowerCapabilities")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	capabilities, err := r.dm.GetPowerCapabilities()
	if err != nil {
		r.l.Error(err, "http - v1 - getPowerCapabilities")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	amtversion := 0

	for _, v := range version {
		if v.InstanceID == "AMT" {
			splitversion := strings.Split(v.VersionString, ".")
			amtversion, err = strconv.Atoi(splitversion[0])

			if err != nil {
				r.l.Error(err, "http - v1 - getPowerCapabilities")
				errorResponse(c, http.StatusInternalServerError, "error converting version")

				return
			}
		}
	}

	response := map[string]interface{}{
		"Power up":    2,
		"Power cycle": 5,
		"Power down":  8,
		"Reset":       10,
	}

	if amtversion > MinAMTVersion {
		response["Soft-off"] = 12
		response["Soft-reset"] = 14
		response["Sleep"] = 4
		response["Hibernate"] = 7
	}

	if capabilities.BIOSSetup {
		response["Power up to BIOS"] = 100
		response["Reset to BIOS"] = 101
	}

	if capabilities.SecureErase {
		response["Reset to Secure Erase"] = 104
	}

	response["Reset to IDE-R Floppy"] = 200
	response["Power on to IDE-R Floppy"] = 201
	response["Reset to IDE-R CDROM"] = 202
	response["Power on to IDE-R CDROM"] = 203

	if capabilities.ForceDiagnosticBoot {
		response["Power on to diagnostic"] = 300
		response["Reset to diagnostic"] = 301
	}

	response["Reset to PXE"] = 400
	response["Power on to PXE"] = 401

	c.JSON(http.StatusOK, response)
}

func (r *deviceManagementRoutes) getGeneralSettings(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - getGeneralSettings")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.dm.SetupWsmanClient(item, true)

	features, err := r.dm.GetGeneralSettings()
	if err != nil {
		r.l.Error(err, "http - v1 - getGeneralSettings")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	response := map[string]interface{}{
		"Body": features,
	}

	c.JSON(http.StatusOK, response)
}

func (r *deviceManagementRoutes) cancelUserConsentCode(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - cancelUserConsentCode")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.dm.SetupWsmanClient(item, true)

	features, err := r.dm.CancelUserConsent()
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

	r.dm.SetupWsmanClient(item, true)

	features, err := r.dm.GetUserConsentCode()
	if err != nil {
		r.l.Error(err, "http - v1 - getUserConsentCode")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	response := map[string]interface{}{
		"Body": features,
	}

	c.JSON(http.StatusOK, response)
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

	r.dm.SetupWsmanClient(item, true)

	features, err := r.dm.SendConsentCode(userConsent.ConsentCode)
	if err != nil {
		r.l.Error(err, "http - v1 - sendConsentCode")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, features)
}

func (r *deviceManagementRoutes) powerAction(c *gin.Context) {
	var powerAction dto.PowerAction
	if err := c.ShouldBindJSON(&powerAction); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - powerAction")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.dm.SetupWsmanClient(item, true)

	features, err := r.dm.SendPowerAction(powerAction.Action)
	if err != nil {
		r.l.Error(err, "http - v1 - powerAction")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, features)
}

func (r *deviceManagementRoutes) getAuditLog(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - getAuditLog")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.dm.SetupWsmanClient(item, true)

	startIndex := c.Query("startIndex")

	startIdx, err := strconv.Atoi(startIndex)
	if err != nil {
		r.l.Error(err, "http - v1 - getAuditLog")
		errorResponse(c, http.StatusInternalServerError, "error converting start index")

		return
	}

	auditlogoutput, err := r.dm.GetAuditLog(startIdx)
	if err != nil {
		r.l.Error(err, "http - v1 - getAuditLog")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, auditlogoutput)
}

func (r *deviceManagementRoutes) getEventLog(c *gin.Context) {
	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - getEventLog")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.dm.SetupWsmanClient(item, true)

	eventLogOutput, err := r.dm.GetEventLog()
	if err != nil {
		r.l.Error(err, "http - v1 - getEventLog")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, eventLogOutput)
}

func (r *deviceManagementRoutes) setBootOptions(c *gin.Context) {
	var bootSetting dto.BootSetting
	if err := c.ShouldBindJSON(&bootSetting); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	item, err := r.d.GetByID(c.Request.Context(), c.Param("guid"), "")
	if err != nil || item.GUID == "" {
		r.l.Error(err, "http - v1 - setBootOptions")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	r.dm.SetupWsmanClient(item, true)

	// bootData, err := r.t.GetBootData()
	// if err != nil {
	// 	r.l.Error(err, "http - v1 - setBootOptions")
	// 	errorResponse(c, http.StatusInternalServerError, "amt problems")

	// 	return
	// }
	newData := boot.BootSettingDataRequest{
		UseSOL:                 bootSetting.UseSOL,
		UseSafeMode:            false,
		ReflashBIOS:            false,
		BIOSSetup:              bootSetting.Action < 104,
		BIOSPause:              false,
		LockPowerButton:        false,
		LockResetButton:        false,
		LockKeyboard:           false,
		LockSleepButton:        false,
		UserPasswordBypass:     false,
		ForcedProgressEvents:   false,
		FirmwareVerbosity:      0,
		ConfigurationDataReset: false,
		UseIDER:                bootSetting.Action > 199 || bootSetting.Action < 300,
		EnforceSecureBoot:      false,
		BootMediaIndex:         0,
		SecureErase:            false,
		RPEEnabled:             false,
		PlatformErase:          false,
	}

	if bootSetting.Action == 202 || bootSetting.Action == 203 {
		newData.IDERBootDevice = 1 // boot on ider
	} else {
		newData.IDERBootDevice = 0 // boot on floppy
	}
	// force boot mode
	_, err = r.dm.SetBootConfigRole(1)
	if err != nil {
		r.l.Error(err, "http - v1 - setBootOptions")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	if bootSetting.Action == 400 || bootSetting.Action == 401 { // pxe boots
		_, err = r.dm.ChangeBootOrder(string(cimBoot.PXE)) // "Intel(r) AMT: Force PXE Boot"
	} else if bootSetting.Action == 202 || bootSetting.Action == 203 {
		_, err = r.dm.ChangeBootOrder(string(cimBoot.CD)) // "Intel(r) AMT: Force CD/DVD Boot"
	}

	if err != nil {
		r.l.Error(err, "http - v1 - changeBootOrder")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	_, err = r.dm.SetBootData(newData)
	if err != nil {
		r.l.Error(err, "http - v1 - setBootOptions")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	if bootSetting.Action == 101 || bootSetting.Action == 200 || bootSetting.Action == 202 || bootSetting.Action == 301 || bootSetting.Action == 400 {
		bootSetting.Action = int(power.MasterBusReset) // reset
	} else {
		bootSetting.Action = int(power.PowerOn) // power on
	}

	features, err := r.dm.SendPowerAction(bootSetting.Action)
	if err != nil {
		r.l.Error(err, "http - v1 - setBootOptions")
		errorResponse(c, http.StatusInternalServerError, "amt problems")

		return
	}

	c.JSON(http.StatusOK, features)
}
