package devices

import (
	"fmt"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/authorization"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/environmentdetection"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/ethernetport"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/general"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/managementpresence"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/publickey"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/redirection"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/remoteaccess"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/setupandconfiguration"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/tls"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/userinitiatedconnection"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/boot"
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

func GetDeviceUUID(wsman wsman.Messages) (uuid string, err error) {
	response, err := wsman.AMT.SetupAndConfigurationService.GetUuid()
	if err != nil {
		return
	}

	uuid, err = response.DecodeUUID()
	if err != nil {
		return
	}
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

func ChangePowerState(wsman wsman.Messages, powerState power.PowerState) (response power.Response, err error) {
	response, err = wsman.CIM.PowerManagementService.RequestPowerStateChange(powerState)
	if err != nil {
		return power.Response{}, err
	}
	return response, nil
}

func getPowerStateValue(technology string, value string) power.PowerState {
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

type Class struct {
	Name       string
	MethodList []Method
}

type Method struct {
	Name string
}

type WsmanMethods struct {
	MethodList []Method
}

func GetSupportedWsmanClasses(className string) []Class {
	var ClassList = []Class{
		{Name: "AMT_AuthorizationService", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "AMT_EnvironmentDetectionSettingData", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "AMT_EthernetPortSettings", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "AMT_GeneralSettings", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "AMT_ManagementPresenceRemoteSAP", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "AMT_PublicKeyCertificate", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "AMT_PublicKeyManagementService", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "AMT_RedirectionService", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "AMT_RemoteAccessPolicyAppliesToMPS", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "AMT_RemoteAccessPolicyRule", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "AMT_RemoteAccessService", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "AMT_TLSSettingData", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "AMT_UserInitiatedConnectionService", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "CIM_BootConfigSetting", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "CIM_BootService", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: "CIM_BootSourceSetting", MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
	}
	if className == "" {
		return ClassList
	}

	for _, class := range ClassList {
		if class.Name == className {
			return []Class{class}
		}
	}
	return []Class{}
}

type Response struct {
	XMLInput  string
	XMLOutput string
}

func MakeWsmanCall(device Device, class string, method string) (response Response, err error) {
	clientParameters := wsman.ClientParameters{
		Target:            device.Address,
		Username:          device.Username,
		Password:          device.Password,
		UseDigest:         true,
		UseTLS:            device.UseTLS,
		SelfSignedAllowed: device.SelfSignedAllowed,
	}
	wsman := wsman.NewMessages(clientParameters)
	selectedClass := GetSupportedWsmanClasses(class)
	switch selectedClass[0].Name {
	case "AMT_AuthorizationService":
		var output authorization.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.AuthorizationService.Get()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.AuthorizationService.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er authorization.Response
			er, err = wsman.AMT.AuthorizationService.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.AuthorizationService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "AMT_EnvironmentDetectionSettingData":
		var output environmentdetection.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.EnvironmentDetectionSettingData.Get()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.EnvironmentDetectionSettingData.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er environmentdetection.Response
			er, err = wsman.AMT.EnvironmentDetectionSettingData.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.EnvironmentDetectionSettingData.Pull(er.Body.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "AMT_EthernetPortSettings":
		var output ethernetport.Response
		switch method {
		case "Get":
			selector := ethernetport.Selector{
				Name:  "InstanceID",
				Value: "Intel(r) AMT Ethernet Port Settings 0",
			}
			output, err = wsman.AMT.EthernetPortSettings.Get(selector)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.EthernetPortSettings.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er ethernetport.Response
			er, err = wsman.AMT.EthernetPortSettings.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.EthernetPortSettings.Pull(er.Body.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "AMT_GeneralSettings":
		var output general.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.GeneralSettings.Get()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.GeneralSettings.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er general.Response
			er, err = wsman.AMT.GeneralSettings.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.GeneralSettings.Pull(er.Body.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "AMT_ManagementPresenceRemoteSAP":
		var output managementpresence.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.ManagementPresenceRemoteSAP.Get()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.ManagementPresenceRemoteSAP.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er managementpresence.Response
			er, err = wsman.AMT.ManagementPresenceRemoteSAP.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.ManagementPresenceRemoteSAP.Pull(er.Body.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "AMT_PublicKeyCertificate":
		var output publickey.ResponseCert
		switch method {
		case "Get":
			output, err = wsman.AMT.PublicKeyCertificate.Get()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.PublicKeyCertificate.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er publickey.ResponseCert
			er, err = wsman.AMT.PublicKeyCertificate.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.PublicKeyCertificate.Pull(er.BodyCert.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "AMT_PublicKeyManagementService":
		var output publickey.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.PublicKeyManagementService.Get()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.PublicKeyManagementService.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er publickey.Response
			er, err = wsman.AMT.PublicKeyManagementService.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.PublicKeyManagementService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "AMT_RedirectionService":
		var output redirection.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.RedirectionService.Get()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.RedirectionService.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er redirection.Response
			er, err = wsman.AMT.RedirectionService.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.RedirectionService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "AMT_RemoteAccessPolicyAppliesToMPS":
		var output remoteaccess.ResponseApplies
		switch method {
		case "Get":
			output, err = wsman.AMT.RemoteAccessPolicyAppliesToMPS.Get()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.RemoteAccessPolicyAppliesToMPS.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er remoteaccess.ResponseApplies
			er, err = wsman.AMT.RemoteAccessPolicyAppliesToMPS.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.RemoteAccessPolicyAppliesToMPS.Pull(er.BodyApplies.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "AMT_RemoteAccessPolicyRule":
		var output remoteaccess.ResponseRule
		switch method {
		case "Get":
			output, err = wsman.AMT.RemoteAccessPolicyRule.Get()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.RemoteAccessPolicyRule.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er remoteaccess.ResponseRule
			er, err = wsman.AMT.RemoteAccessPolicyRule.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.RemoteAccessPolicyRule.Pull(er.BodyRule.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "AMT_RemoteAccessService":
		var output remoteaccess.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.RemoteAccessService.Get()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.RemoteAccessService.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er remoteaccess.Response
			er, err = wsman.AMT.RemoteAccessService.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.RemoteAccessService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "AMT_TLSSettingData":
		var output tls.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.TLSSettingData.Get("Intel(r) AMT 802.3 TLS Settings")
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.TLSSettingData.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er tls.Response
			er, err = wsman.AMT.TLSSettingData.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.TLSSettingData.Pull(er.Body.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "AMT_UserInitiatedConnectionService":
		var output userinitiatedconnection.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.UserInitiatedConnectionService.Get()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.UserInitiatedConnectionService.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er userinitiatedconnection.Response
			er, err = wsman.AMT.UserInitiatedConnectionService.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.UserInitiatedConnectionService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "CIM_BootConfigSetting":
		var output boot.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.BootConfigSetting.Get()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.BootConfigSetting.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er boot.Response
			er, err = wsman.CIM.BootConfigSetting.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.BootConfigSetting.Pull(er.Body.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "CIM_BootService":
		var output boot.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.BootService.Get()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.BootService.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er boot.Response
			er, err = wsman.CIM.BootService.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.BootService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	case "CIM_BootSourceSetting":
		var output boot.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.BootSourceSetting.Get("Intel(r) AMT: Force Hard-drive Boot")
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.BootSourceSetting.Enumerate()
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er boot.Response
			er, err = wsman.CIM.BootSourceSetting.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.BootSourceSetting.Pull(er.Body.EnumerateResponse.EnumerationContext)
			if err != nil {
				return
			}
			response.XMLInput = output.XMLInput
			response.XMLOutput = output.XMLOutput
			return
		}
	}
	return
}
