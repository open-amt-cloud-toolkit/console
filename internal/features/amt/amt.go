package amt

import (
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/setupandconfiguration"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/power"
)

func ProvisioningModeLookup(mode setupandconfiguration.ProvisioningModeValue) string {
	valueMap := map[setupandconfiguration.ProvisioningModeValue]string{
		1: "Admin Control Mode",
		4: "Client Control Mode",
	}

	result, ok := valueMap[mode]
	if !ok {
		result = "invalid provisioning mode"
	}

	return result
}

func ProvisioningStateLookup(state setupandconfiguration.ProvisioningStateValue) string {
	valueMap := map[setupandconfiguration.ProvisioningStateValue]string{
		0: "Pre-Provisioning",
		1: "In Provisioning",
		2: "Post Provisioning",
	}

	result, ok := valueMap[state]
	if !ok {
		result = "invalid provisoining state"
	}

	return result
}

func PowerControlLookup(value int) string {
	valueMap := map[int]string{
		2:  "On",
		3:  "Sleep - Light",
		4:  "Sleep - Deep",
		5:  "Power Cycle Off - Soft",
		6:  "Power Off - Hard",
		7:  "Hibernate",
		8:  "Power Off - Soft",
		9:  "Power Cycle Off - Hard",
		10: "Master Bus Reset",
		11: "Diagnostic Interrupt NMI",
		12: "Power Off - Soft Graceful",
		13: "Power Off - Hard Graceful",
		14: "Master Bus Reset - Graceful",
		15: "Power Cycle Off - Soft Graceful",
		16: "Power Cycle Off - Hard Graceful",
	}

	result, ok := valueMap[value]
	if !ok {
		result = "invalid power control value"
	}

	return result
}

func PowerControlReturnValue(value int) string {
	valueMap := map[int]string{
		0:    "Completed with No Error",
		1:    "Not Supported",
		2:    "Unknown or Unspecified Error",
		3:    "Cannot complete within Timeout Period",
		4:    "Failed",
		5:    "Invalid Parameter",
		6:    "In Use",
		4096: "Method Parameters Checked - Job Started",
		4097: "Invalid State Transition",
		4098: "Use of Timeout Parameter Not Supported",
		4099: "Busy",
	}

	result, ok := valueMap[value]
	if !ok {
		result = "invalid power control value"
	}

	return result
}

func PowerStateLookup(value int) string {
	valueMap := map[int]string{
		2:  "On",
		3:  "Sleep - Light",
		4:  "Sleep - Deep",
		5:  "Power Cycle Off - Soft",
		6:  "Power Off - Hard",
		7:  "Hibernate",
		8:  "Power Off - Soft",
		9:  "Power Cycle Off - Hard",
		10: "Master Bus Reset",
		11: "Diagnostic Interrupt NMI",
		12: "Power Off - Soft Graceful",
		13: "Power Off - Hard Graceful",
		14: "Master Bus Reset - Graceful",
		15: "Power Cycle Off - Soft Graceful",
		16: "Power Cycle Off - Hard Graceful",
	}

	result, ok := valueMap[value]
	if !ok {
		result = "invalid power control value"
	}

	return result
}

func CreateWsmanConnection(amtConnectionParameters AMTConnectionParameters) wsman.Messages {
	cp := wsman.ClientParameters{
		Target:            amtConnectionParameters.Target,
		Username:          amtConnectionParameters.Username,
		Password:          amtConnectionParameters.Password,
		UseDigest:         true,
		UseTLS:            amtConnectionParameters.UseTLS,
		SelfSignedAllowed: amtConnectionParameters.SelfSignedAllowed,
	}
	wsman := wsman.NewMessages(cp)
	return wsman
}

func GetDeviceDetails(wsman wsman.Messages) (amtDeviceDetails AMTSpecific) {
	gs, err := wsman.AMT.GeneralSettings.Get()
	if err != nil {
		amtDeviceDetails.Errors = append(amtDeviceDetails.Errors, err)
	}
	amtDeviceDetails.GeneralSettings = gs.Body.GetResponse

	scs, err := wsman.AMT.SetupAndConfigurationService.Get()
	if err != nil {
		amtDeviceDetails.Errors = append(amtDeviceDetails.Errors, err)
	}
	amtDeviceDetails.SetupAndConfigurationService = scs.Body.GetResponse

	uuid, err := wsman.AMT.SetupAndConfigurationService.GetUuid()
	if err != nil {
		amtDeviceDetails.Errors = append(amtDeviceDetails.Errors, err)
	}
	amtDeviceDetails.UUID, err = uuid.DecodeUUID()
	if err != nil {
		amtDeviceDetails.Errors = append(amtDeviceDetails.Errors, err)
	}

	eth0, err := wsman.AMT.EthernetPortSettings.Get(0)
	if err != nil {
		amtDeviceDetails.Errors = append(amtDeviceDetails.Errors, err)
	}
	amtDeviceDetails.EthernetSettings.Wired = eth0.Body.GetAndPutResponse

	eth1, err := wsman.AMT.EthernetPortSettings.Get(1)
	if err != nil {
		amtDeviceDetails.Errors = append(amtDeviceDetails.Errors, err)
	}
	amtDeviceDetails.EthernetSettings.Wireless = eth1.Body.GetAndPutResponse
	return amtDeviceDetails
}

func GetPowerState(wsman wsman.Messages) (powerState string, err error) {
	er, err := wsman.CIM.ServiceAvailableToElement.Enumerate()
	if err != nil {
		return
	}
	response, err := wsman.CIM.ServiceAvailableToElement.Pull(er.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return
	}
	powerState = PowerControlLookup(response.Body.PullResponse.AssociatedPowerManagementService[0].PowerState)
	return
}

func ChangePowerState(wsman wsman.Messages, powerState power.PowerState) (powerResponse string, err error) {
	response, err := wsman.CIM.PowerManagementService.RequestPowerStateChange(powerState)
	if err != nil {
		return "", err
	}
	powerResponse = PowerControlReturnValue(response.Body.RequestPowerStateChangeResponse.ReturnValue)
	return powerResponse, nil
}

func GetPowerStateValue(technology string, value string) power.PowerState {
	switch value {
	case "on":
		if technology == "amt" {
			return power.PowerOn // 2
		}
	case "off":
		if technology == "amt" {
			return power.PowerOffHard // 8
		}
	case "reboot":
		if technology == "amt" {
			return power.MasterBusReset // 10
		}
	case "powercycle":
		if technology == "amt" {
			return power.PowerCycleOffHard // 5
		}
	}
	return 0
}
