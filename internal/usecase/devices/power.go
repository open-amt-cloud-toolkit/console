package devices

import (
	"context"
	"strconv"
	"strings"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/boot"
	cimBoot "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/boot"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/power"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/software"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/utils"
)

func (uc *UseCase) SendPowerAction(c context.Context, guid string, action int) (power.PowerActionResponse, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return power.PowerActionResponse{}, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	response, err := uc.device.SendPowerAction(action)
	if err != nil {
		return power.PowerActionResponse{}, err
	}

	return response, nil
}

func (uc *UseCase) GetPowerState(c context.Context, guid string) (map[string]interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	state, err := uc.device.GetPowerState()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"powerstate": state[0].PowerState,
	}, nil
}

func (uc *UseCase) GetPowerCapabilities(c context.Context, guid string) (map[string]interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	version, err := uc.device.GetAMTVersion()
	if err != nil {
		return nil, err
	}

	capabilities, err := uc.device.GetPowerCapabilities()
	if err != nil {
		return nil, err
	}

	amtversion, err := parseVersion(version)
	if err != nil {
		return nil, utils.ErrParseVersion
	}

	response := determinePowerCapabilities(amtversion, capabilities)

	return response, nil
}

func determinePowerCapabilities(amtversion int, capabilities boot.BootCapabilitiesResponse) map[string]interface{} {
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

	return response
}

func (uc *UseCase) SetBootOptions(c context.Context, guid string, bootSetting dto.BootSetting) (power.PowerActionResponse, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return power.PowerActionResponse{}, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

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

	// boot on ider
	// boot on floppy
	determineIDERBootDevice(bootSetting, &newData)
	// force boot mode
	_, err = uc.device.SetBootConfigRole(1)
	if err != nil {
		return power.PowerActionResponse{}, err
	}

	bootSource := getBootSource(bootSetting)
	if bootSource != "" {
		_, err = uc.device.ChangeBootOrder(bootSource)
		if err != nil {
			return power.PowerActionResponse{}, err
		}
	}

	_, err = uc.device.SetBootData(newData)
	if err != nil {
		return power.PowerActionResponse{}, err
	}

	// reset
	// power on
	determineBootAction(&bootSetting)

	powerActionResult, err := uc.device.SendPowerAction(bootSetting.Action)
	if err != nil {
		return power.PowerActionResponse{}, err
	}

	return powerActionResult, nil
}

func determineIDERBootDevice(bootSetting dto.BootSetting, newData *boot.BootSettingDataRequest) {
	if bootSetting.Action == 202 || bootSetting.Action == 203 {
		newData.IDERBootDevice = 1
	} else {
		newData.IDERBootDevice = 0
	}
}

// "Intel(r) AMT: Force PXE Boot".
// "Intel(r) AMT: Force CD/DVD Boot".
func getBootSource(bootSetting dto.BootSetting) string {
	if bootSetting.Action == 400 || bootSetting.Action == 401 {
		return string(cimBoot.PXE)
	} else if bootSetting.Action == 202 || bootSetting.Action == 203 {
		return string(cimBoot.CD)
	}

	return ""
}

func determineBootAction(bootSetting *dto.BootSetting) {
	if bootSetting.Action == 101 || bootSetting.Action == 200 || bootSetting.Action == 202 || bootSetting.Action == 301 || bootSetting.Action == 400 {
		bootSetting.Action = int(power.MasterBusReset)
	} else {
		bootSetting.Action = int(power.PowerOn)
	}
}

func parseVersion(version []software.SoftwareIdentity) (int, error) {
	amtversion := 0

	var err error

	for _, v := range version {
		if v.InstanceID == "AMT" {
			splitversion := strings.Split(v.VersionString, ".")

			amtversion, err = strconv.Atoi(splitversion[0])
			if err != nil {
				return 0, err
			}
		}
	}

	return amtversion, nil
}
