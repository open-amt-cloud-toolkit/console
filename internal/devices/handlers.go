package devices

import (
	"fmt"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/ethernetport"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/general"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/setupandconfiguration"
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
		result = "invalid provisoining state"
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

// func CreateFirmwareConnection() firmware.Messages {
// 	return firmware.NewMessages()
// }

func CreateWsmanConnection(d Device) wsman.Messages {
	cp := wsman.ClientParameters{
		Target:            d.Address,
		Username:          d.Username,
		Password:          d.Password,
		UseDigest:         true,
		UseTLS:            d.UseTLS,
		SelfSignedAllowed: d.SelfSignedAllowed,
	}
	wsman := wsman.NewMessages(cp)
	return wsman
}

func GetGeneralSettings(wsman wsman.Messages) (gs general.GeneralSettings, err error) {
	response, err := wsman.AMT.GeneralSettings.Get()
	if err != nil {
		return
	}
	gs = response.Body.AMTGeneralSettings
	return
}

func GetEthernetSettings(wsman wsman.Messages, eth int) (ep ethernetport.EthernetPort, err error) {
	selector := ethernetport.Selector{
		Name:  "InstanceID",
		Value: fmt.Sprintf("Intel(r) AMT Ethernet Port Settings %d", eth),
	}

	response, err := wsman.AMT.EthernetPortSettings.Get(selector)
	if err != nil {
		return
	}
	ep = response.Body.EthernetPort

	return ep, err
}

func GetSetupAndConfigurationService(wsman wsman.Messages) (sc setupandconfiguration.Setup, err error) {
	response, err := wsman.AMT.SetupAndConfigurationService.Get()
	if err != nil {
		return
	}
	sc = response.Body.Setup
	return
}

type PowerState string

const (
	PowerOn                   PowerState = "Power On"
	SleepLight                PowerState = "Sleep Light (OS)"
	SleepDeep                 PowerState = "Sleep Deep (OS)"
	PowerCycleOffSoft         PowerState = "Soft Power Cycle (OS Graceful)"
	PowerOffHard              PowerState = "Hard Power Off"
	Hibernate                 PowerState = "Hibernate (OS)"
	PowerOffSoft              PowerState = "Soft Power Off (OS Graceful)"
	PowerCycleOffHard         PowerState = "Hard Power Cycle"
	MasterBusReset            PowerState = "Master Bus Reset"
	DiagnosticInterruptNMI    PowerState = "Diagnostic Interrupt NMI"
	PowerOffSoftGraceful      PowerState = "Soft Power Off (OS Graceful)"
	PowerCycleOffHardGraceful PowerState = "Hard Power Cycle (OS Graceful)"
)

// func ChangePowerState(wsman wsman.Messages, powerState power.PowerState) (response power.Response, err error) {
// 	response, err = wsman.CIM.PowerManagementService.RequestPowerStateChange(powerState)
// 	if err != nil {
// 		return power.Response{}, err
// 	}
// 	return response, nil
// }
