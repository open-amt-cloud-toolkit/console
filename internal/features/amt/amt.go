package amt

import (
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/ethernetport"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/power"
)

func ProvisioningModeLookup(mode int) string {
	valueMap := map[int]string{
		1: "Admin Control Mode",
		4: "Client Control Mode",
	}

	result, ok := valueMap[mode]
	if !ok {
		result = "invalid provisioning mode"
	}

	return result
}

func ProvisioningStateLookup(state int) string {
	valueMap := map[int]string{
		0: "Pre-Provisioning",
		1: "In Provisioning",
		2: "Post Provisioning",
	}

	result, ok := valueMap[state]
	if !ok {
		result = "invalid provisioning state"
	}

	return result
}

func PowerControlLookup(value int) string {
	valueMap := map[int]string{
		2:  "PowerOn",
		3:  "SleepLight",
		4:  "SleepDeep",
		5:  "PowerCycleOffSoft",
		6:  "PowerOffHard",
		7:  "Hibernate",
		8:  "PowerOffSoft",
		9:  "PowerCycleOffHard",
		10: "MasterBusReset",
		11: "DiagnosticInterruptNMI",
		12: "PowerOffSoftGraceful",
		13: "PowerOffHardGraceful",
		14: "MasterBusResetGraceful",
		15: "PowerCycleOffSoftGraceful",
		16: "PowerCycleOffHardGraceful",
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
	amtDeviceDetails.GeneralSettings = gs.Body.AMTGeneralSettings

	scs, err := wsman.AMT.SetupAndConfigurationService.Get()
	if err != nil {
		amtDeviceDetails.Errors = append(amtDeviceDetails.Errors, err)
	}
	amtDeviceDetails.SetupAndConfigurationService = scs.Body.Setup

	uuid, err := wsman.AMT.SetupAndConfigurationService.GetUuid()
	if err != nil {
		amtDeviceDetails.Errors = append(amtDeviceDetails.Errors, err)
	}
	amtDeviceDetails.UUID, err = uuid.DecodeUUID()
	if err != nil {
		amtDeviceDetails.Errors = append(amtDeviceDetails.Errors, err)
	}
	wiredSelector := ethernetport.Selector{
		Name:  "InstanceID",
		Value: "Intel(r) AMT Ethernet Port Settings 0",
	}
	eth0, err := wsman.AMT.EthernetPortSettings.Get(wiredSelector)
	if err != nil {
		amtDeviceDetails.Errors = append(amtDeviceDetails.Errors, err)
	}
	amtDeviceDetails.EthernetSettings.Wired = eth0.Body.EthernetPort
	wirelessSelector := ethernetport.Selector{
		Name:  "InstanceID",
		Value: "Intel(r) AMT Ethernet Port Settings 1",
	}
	eth1, err := wsman.AMT.EthernetPortSettings.Get(wirelessSelector)
	if err != nil {
		amtDeviceDetails.Errors = append(amtDeviceDetails.Errors, err)
	}
	amtDeviceDetails.EthernetSettings.Wireless = eth1.Body.EthernetPort
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

func ChangePowerState(wsman wsman.Messages, powerState power.PowerState) (response power.Response, err error) {
	response, err = wsman.CIM.PowerManagementService.RequestPowerStateChange(powerState)
	if err != nil {
		return power.Response{}, err
	}
	return response, nil
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
