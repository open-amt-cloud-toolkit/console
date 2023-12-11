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
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/bios"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/boot"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/computer"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/credential"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/ieee8021x"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/kvm"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/mediaaccess"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/physical"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/power"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/service"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/software"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/system"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/wifi"
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
		{Name: authorization.AMT_AuthorizationService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: environmentdetection.AMT_EnvironmentDetectionSettingData, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: ethernetport.AMT_EthernetPortSettings, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: general.AMT_GeneralSettings, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: managementpresence.AMT_ManagementPresenceRemoteSAP, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: publickey.AMT_PublicKeyCertificate, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: publickey.AMT_PublicKeyManagementService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: redirection.AMT_RedirectionService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: remoteaccess.AMT_RemoteAccessPolicyAppliesToMPS, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: remoteaccess.AMT_RemoteAccessPolicyRule, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: remoteaccess.AMT_RemoteAccessService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: setupandconfiguration.AMT_SetupAndConfigurationService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: tls.AMT_TLSSettingData, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: userinitiatedconnection.AMT_UserInitiatedConnectionService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: boot.CIM_BootConfigSetting, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: boot.CIM_BootService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: boot.CIM_BootSourceSetting, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: bios.CIM_BIOSElement, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: computer.CIM_ComputerSystemPackage, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: credential.CIM_CredentialContext, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: ieee8021x.CIM_IEEE8021xSettings, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: kvm.CIM_KVMRedirectionSAP, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: mediaaccess.CIM_MediaAccessDevice, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: physical.CIM_Card, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: physical.CIM_Chassis, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: physical.CIM_Chip, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: physical.CIM_PhysicalMemory, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: physical.CIM_PhysicalPackage, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: physical.CIM_Processor, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: power.CIM_PowerManagementService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: service.CIM_ServiceAvailableToElement, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: software.CIM_SoftwareIdentity, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: system.CIM_SystemPackaging, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: wifi.CIM_WiFiEndpointSettings, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: wifi.CIM_WiFiPort, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
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
	case authorization.AMT_AuthorizationService:
		var output authorization.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.AuthorizationService.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.AuthorizationService.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er authorization.Response
			er, err = wsman.AMT.AuthorizationService.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.AuthorizationService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case environmentdetection.AMT_EnvironmentDetectionSettingData:
		var output environmentdetection.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.EnvironmentDetectionSettingData.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.EnvironmentDetectionSettingData.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er environmentdetection.Response
			er, err = wsman.AMT.EnvironmentDetectionSettingData.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.EnvironmentDetectionSettingData.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case ethernetport.AMT_EthernetPortSettings:
		var output ethernetport.Response
		switch method {
		case "Get":
			selector := ethernetport.Selector{
				Name:  "InstanceID",
				Value: "Intel(r) AMT Ethernet Port Settings 0",
			}
			output, err = wsman.AMT.EthernetPortSettings.Get(selector)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.EthernetPortSettings.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er ethernetport.Response
			er, err = wsman.AMT.EthernetPortSettings.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.EthernetPortSettings.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case general.AMT_GeneralSettings:
		var output general.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.GeneralSettings.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.GeneralSettings.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er general.Response
			er, err = wsman.AMT.GeneralSettings.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.GeneralSettings.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case managementpresence.AMT_ManagementPresenceRemoteSAP:
		var output managementpresence.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.ManagementPresenceRemoteSAP.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.ManagementPresenceRemoteSAP.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er managementpresence.Response
			er, err = wsman.AMT.ManagementPresenceRemoteSAP.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.ManagementPresenceRemoteSAP.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case publickey.AMT_PublicKeyCertificate:
		var output publickey.ResponseCert
		switch method {
		case "Get":
			output, err = wsman.AMT.PublicKeyCertificate.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.PublicKeyCertificate.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er publickey.ResponseCert
			er, err = wsman.AMT.PublicKeyCertificate.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.PublicKeyCertificate.Pull(er.BodyCert.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case publickey.AMT_PublicKeyManagementService:
		var output publickey.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.PublicKeyManagementService.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.PublicKeyManagementService.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er publickey.Response
			er, err = wsman.AMT.PublicKeyManagementService.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.PublicKeyManagementService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case redirection.AMT_RedirectionService:
		var output redirection.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.RedirectionService.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.RedirectionService.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er redirection.Response
			er, err = wsman.AMT.RedirectionService.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.RedirectionService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case remoteaccess.AMT_RemoteAccessPolicyAppliesToMPS:
		var output remoteaccess.ResponseApplies
		switch method {
		case "Get":
			output, err = wsman.AMT.RemoteAccessPolicyAppliesToMPS.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.RemoteAccessPolicyAppliesToMPS.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er remoteaccess.ResponseApplies
			er, err = wsman.AMT.RemoteAccessPolicyAppliesToMPS.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.RemoteAccessPolicyAppliesToMPS.Pull(er.BodyApplies.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case remoteaccess.AMT_RemoteAccessPolicyRule:
		var output remoteaccess.ResponseRule
		switch method {
		case "Get":
			output, err = wsman.AMT.RemoteAccessPolicyRule.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.RemoteAccessPolicyRule.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er remoteaccess.ResponseRule
			er, err = wsman.AMT.RemoteAccessPolicyRule.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.RemoteAccessPolicyRule.Pull(er.BodyRule.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case remoteaccess.AMT_RemoteAccessService:
		var output remoteaccess.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.RemoteAccessService.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.RemoteAccessService.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er remoteaccess.Response
			er, err = wsman.AMT.RemoteAccessService.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.RemoteAccessService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case setupandconfiguration.AMT_SetupAndConfigurationService:
		var output setupandconfiguration.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.SetupAndConfigurationService.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.SetupAndConfigurationService.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er setupandconfiguration.Response
			er, err = wsman.AMT.SetupAndConfigurationService.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.SetupAndConfigurationService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case tls.AMT_TLSSettingData:
		var output tls.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.TLSSettingData.Get("Intel(r) AMT 802.3 TLS Settings")
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.TLSSettingData.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er tls.Response
			er, err = wsman.AMT.TLSSettingData.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.TLSSettingData.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case userinitiatedconnection.AMT_UserInitiatedConnectionService:
		var output userinitiatedconnection.Response
		switch method {
		case "Get":
			output, err = wsman.AMT.UserInitiatedConnectionService.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.AMT.UserInitiatedConnectionService.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er userinitiatedconnection.Response
			er, err = wsman.AMT.UserInitiatedConnectionService.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.AMT.UserInitiatedConnectionService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case bios.CIM_BIOSElement:
		var output bios.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.BIOSElement.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.BIOSElement.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er bios.Response
			er, err = wsman.CIM.BIOSElement.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.BIOSElement.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case boot.CIM_BootConfigSetting:
		var output boot.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.BootConfigSetting.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.BootConfigSetting.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er boot.Response
			er, err = wsman.CIM.BootConfigSetting.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.BootConfigSetting.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case boot.CIM_BootService:
		var output boot.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.BootService.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.BootService.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er boot.Response
			er, err = wsman.CIM.BootService.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.BootService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case boot.CIM_BootSourceSetting:
		var output boot.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.BootSourceSetting.Get("Intel(r) AMT: Force Hard-drive Boot")
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.BootSourceSetting.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er boot.Response
			er, err = wsman.CIM.BootSourceSetting.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.BootSourceSetting.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case computer.CIM_ComputerSystemPackage:
		var output computer.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.ComputerSystemPackage.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.ComputerSystemPackage.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er computer.Response
			er, err = wsman.CIM.ComputerSystemPackage.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.ComputerSystemPackage.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case credential.CIM_CredentialContext:
		var output credential.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.CredentialContext.Get("InstanceID")
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.CredentialContext.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er credential.Response
			er, err = wsman.CIM.CredentialContext.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.CredentialContext.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case ieee8021x.CIM_IEEE8021xSettings:
		var output ieee8021x.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.IEEE8021xSettings.Get("Intel(r) AMT: 8021X Settings")
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.IEEE8021xSettings.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er ieee8021x.Response
			er, err = wsman.CIM.IEEE8021xSettings.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.IEEE8021xSettings.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case kvm.CIM_KVMRedirectionSAP:
		var output kvm.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.KVMRedirectionSAP.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.KVMRedirectionSAP.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er kvm.Response
			er, err = wsman.CIM.KVMRedirectionSAP.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.KVMRedirectionSAP.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case mediaaccess.CIM_MediaAccessDevice:
		var output mediaaccess.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.MediaAccessDevice.Get("MEDIA DEV 0")
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.MediaAccessDevice.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er mediaaccess.Response
			er, err = wsman.CIM.MediaAccessDevice.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.MediaAccessDevice.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case physical.CIM_Card:
		var output physical.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.Card.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.Card.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er physical.Response
			er, err = wsman.CIM.Card.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.Card.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case physical.CIM_Chassis:
		var output physical.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.Chassis.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.Chassis.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er physical.Response
			er, err = wsman.CIM.Chassis.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.Chassis.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case physical.CIM_Chip:
		var output physical.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.Chip.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.Chip.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er physical.Response
			er, err = wsman.CIM.Chip.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.Chip.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case physical.CIM_PhysicalMemory:
		var output physical.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.PhysicalMemory.Get("0")
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.PhysicalMemory.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er physical.Response
			er, err = wsman.CIM.PhysicalMemory.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.PhysicalMemory.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case physical.CIM_PhysicalPackage:
		var output physical.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.PhysicalPackage.Get("0")
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.PhysicalPackage.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er physical.Response
			er, err = wsman.CIM.PhysicalPackage.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.PhysicalPackage.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case physical.CIM_Processor:
		var output physical.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.Processor.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.Processor.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er physical.Response
			er, err = wsman.CIM.Processor.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.Processor.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case power.CIM_PowerManagementService:
		var output power.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.PowerManagementService.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.PowerManagementService.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er power.Response
			er, err = wsman.CIM.PowerManagementService.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.PowerManagementService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case service.CIM_ServiceAvailableToElement:
		var output service.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.ServiceAvailableToElement.Get("CIM_PowerManagementService")
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.ServiceAvailableToElement.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er service.Response
			er, err = wsman.CIM.ServiceAvailableToElement.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.ServiceAvailableToElement.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case software.CIM_SoftwareIdentity:
		var output software.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.SoftwareIdentity.Get("AMT")
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.SoftwareIdentity.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er software.Response
			er, err = wsman.CIM.SoftwareIdentity.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.SoftwareIdentity.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case system.CIM_SystemPackaging:
		var output system.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.SystemPackaging.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.SystemPackaging.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er system.Response
			er, err = wsman.CIM.SystemPackaging.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.SystemPackaging.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case wifi.CIM_WiFiEndpointSettings:
		var output wifi.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.WiFiEndpointSettings.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.WiFiEndpointSettings.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er wifi.Response
			er, err = wsman.CIM.WiFiEndpointSettings.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.WiFiEndpointSettings.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	case wifi.CIM_WiFiPort:
		var output wifi.Response
		switch method {
		case "Get":
			output, err = wsman.CIM.WiFiPort.Get()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Enumerate":
			output, err = wsman.CIM.WiFiPort.Enumerate()
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		case "Pull":
			var er wifi.Response
			er, err = wsman.CIM.WiFiPort.Enumerate()
			if err != nil {
				return
			}
			output, err = wsman.CIM.WiFiPort.Pull(er.Body.EnumerateResponse.EnumerationContext)
			response.XMLInput = output.XMLInput
			if err != nil {
				response.XMLOutput = error.Error(err)
				return
			}
			response.XMLOutput = output.XMLOutput
			return
		}
	}
	return
}
