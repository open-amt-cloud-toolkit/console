package explorer

import (
	"strconv"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/alarmclock"
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
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/card"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/chassis"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/chip"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/computer"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/concrete"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/credential"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/ieee8021x"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/kvm"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/mediaaccess"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/physical"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/power"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/processor"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/service"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/software"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/system"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/cim/wifi"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/client"
)

func GetSupportedWsmanClasses(className string) []Class {
	var ClassList = []Class{
		{Name: authorization.AMT_AuthorizationService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: alarmclock.AMT_AlarmClockService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
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
		{Name: setupandconfiguration.AMT_SetupAndConfigurationService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "GetUUID"}}},
		{Name: tls.AMT_TLSSettingData, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: userinitiatedconnection.AMT_UserInitiatedConnectionService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: bios.CIM_BIOSElement, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: boot.CIM_BootConfigSetting, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: boot.CIM_BootService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: boot.CIM_BootSourceSetting, MethodList: []Method{{Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: computer.CIM_ComputerSystemPackage, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: concrete.CIM_ConcreteDependency, MethodList: []Method{{Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: credential.CIM_CredentialContext, MethodList: []Method{{Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: ieee8021x.CIM_IEEE8021xSettings, MethodList: []Method{{Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: kvm.CIM_KVMRedirectionSAP, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: mediaaccess.CIM_MediaAccessDevice, MethodList: []Method{{Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: card.CIM_Card, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: chassis.CIM_Chassis, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: chip.CIM_Chip, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: physical.CIM_PhysicalMemory, MethodList: []Method{{Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: physical.CIM_PhysicalPackage, MethodList: []Method{{Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: processor.CIM_Processor, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: power.CIM_PowerManagementService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "RequestPowerStateChange"}}},
		{Name: service.CIM_ServiceAvailableToElement, MethodList: []Method{{Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: software.CIM_SoftwareIdentity, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: system.CIM_SystemPackaging, MethodList: []Method{{Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: wifi.CIM_WiFiEndpointSettings, MethodList: []Method{{Name: "Enumerate"}, {Name: "Pull"}}},
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

func MakeWsmanCall(class string, method string, param string) (response Response, err error) {

	selectedClass := GetSupportedWsmanClasses(class)

	output, err := Lookup[selectedClass[0].Name][method].Execute(param)
	response.XMLInput = output.XMLInput
	if err != nil {
		response.XMLOutput = error.Error(err)
		return
	}
	response.XMLOutput = output.XMLOutput
	return
}

var Lookup map[string]map[string]Method

func Init(wsman wsman.Messages) {
	Lookup = make(map[string]map[string]Method)

	// Alarm Clock
	Lookup[alarmclock.AMT_AlarmClockService] = make(map[string]Method)
	Lookup[alarmclock.AMT_AlarmClockService]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AlarmClockService.Get()
			return *response.Message, err
		},
	}
	Lookup[alarmclock.AMT_AlarmClockService]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AlarmClockService.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[alarmclock.AMT_AlarmClockService]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.AlarmClockService.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.AlarmClockService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// AuthorizationService
	Lookup[authorization.AMT_AuthorizationService] = make(map[string]Method)
	Lookup[authorization.AMT_AuthorizationService]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AuthorizationService.Get()
			return *response.Message, err
		},
	}
	Lookup[authorization.AMT_AuthorizationService]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AuthorizationService.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[authorization.AMT_AuthorizationService]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.AuthorizationService.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.AuthorizationService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// EnvironementDetectionSettingData
	Lookup[environmentdetection.AMT_EnvironmentDetectionSettingData] = make(map[string]Method)
	Lookup[environmentdetection.AMT_EnvironmentDetectionSettingData]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.EnvironmentDetectionSettingData.Get()
			return *response.Message, err
		},
	}
	Lookup[environmentdetection.AMT_EnvironmentDetectionSettingData]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.EnvironmentDetectionSettingData.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[environmentdetection.AMT_EnvironmentDetectionSettingData]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.EnvironmentDetectionSettingData.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.EnvironmentDetectionSettingData.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// GeneralSetting
	Lookup[general.AMT_GeneralSettings] = make(map[string]Method)
	Lookup[general.AMT_GeneralSettings]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.GeneralSettings.Get()
			return *response.Message, err
		},
	}
	Lookup[general.AMT_GeneralSettings]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.GeneralSettings.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[general.AMT_GeneralSettings]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.GeneralSettings.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.GeneralSettings.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// ManagementPresenceRemoteSAP
	Lookup[managementpresence.AMT_ManagementPresenceRemoteSAP] = make(map[string]Method)
	Lookup[managementpresence.AMT_ManagementPresenceRemoteSAP]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.ManagementPresenceRemoteSAP.Get()
			return *response.Message, err
		},
	}
	Lookup[managementpresence.AMT_ManagementPresenceRemoteSAP]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.ManagementPresenceRemoteSAP.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[managementpresence.AMT_ManagementPresenceRemoteSAP]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.ManagementPresenceRemoteSAP.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.ManagementPresenceRemoteSAP.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// PublicKeyCertificate
	Lookup[publickey.AMT_PublicKeyCertificate] = make(map[string]Method)
	Lookup[publickey.AMT_PublicKeyCertificate]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.PublicKeyCertificate.Get()
			return *response.Message, err
		},
	}
	Lookup[publickey.AMT_PublicKeyCertificate]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.PublicKeyCertificate.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[publickey.AMT_PublicKeyCertificate]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.PublicKeyCertificate.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.PublicKeyCertificate.Pull(er.BodyCert.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// PublicKeyManagementService
	Lookup[publickey.AMT_PublicKeyManagementService] = make(map[string]Method)
	Lookup[publickey.AMT_PublicKeyManagementService]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.PublicKeyManagementService.Get()
			return *response.Message, err
		},
	}
	Lookup[publickey.AMT_PublicKeyManagementService]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.PublicKeyManagementService.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[publickey.AMT_PublicKeyManagementService]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.PublicKeyManagementService.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.PublicKeyManagementService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// RedirectionService
	Lookup[redirection.AMT_RedirectionService] = make(map[string]Method)
	Lookup[redirection.AMT_RedirectionService]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.RedirectionService.Get()
			return *response.Message, err
		},
	}
	Lookup[redirection.AMT_RedirectionService]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.RedirectionService.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[redirection.AMT_RedirectionService]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.RedirectionService.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.RedirectionService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// RemoteAccessPolicyAppliesToMPS
	Lookup[remoteaccess.AMT_RemoteAccessPolicyAppliesToMPS] = make(map[string]Method)
	Lookup[remoteaccess.AMT_RemoteAccessPolicyAppliesToMPS]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.RemoteAccessPolicyAppliesToMPS.Get()
			return *response.Message, err
		},
	}
	Lookup[remoteaccess.AMT_RemoteAccessPolicyAppliesToMPS]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.RemoteAccessPolicyAppliesToMPS.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[remoteaccess.AMT_RemoteAccessPolicyAppliesToMPS]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.RemoteAccessPolicyAppliesToMPS.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.RemoteAccessPolicyAppliesToMPS.Pull(er.BodyApplies.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// RemoteAccessPolicyRule
	Lookup[remoteaccess.AMT_RemoteAccessPolicyRule] = make(map[string]Method)
	Lookup[remoteaccess.AMT_RemoteAccessPolicyRule]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.RemoteAccessPolicyRule.Get()
			return *response.Message, err
		},
	}
	Lookup[remoteaccess.AMT_RemoteAccessPolicyRule]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.RemoteAccessPolicyRule.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[remoteaccess.AMT_RemoteAccessPolicyRule]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.RemoteAccessPolicyRule.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.RemoteAccessPolicyRule.Pull(er.BodyRule.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// RemoteAccessService
	Lookup[remoteaccess.AMT_RemoteAccessService] = make(map[string]Method)
	Lookup[remoteaccess.AMT_RemoteAccessService]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.RemoteAccessService.Get()
			return *response.Message, err
		},
	}
	Lookup[remoteaccess.AMT_RemoteAccessService]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.RemoteAccessService.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[remoteaccess.AMT_RemoteAccessService]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.RemoteAccessService.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.RemoteAccessService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// SetupAndConfigurationService
	Lookup[setupandconfiguration.AMT_SetupAndConfigurationService] = make(map[string]Method)
	Lookup[setupandconfiguration.AMT_SetupAndConfigurationService]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.SetupAndConfigurationService.Get()
			return *response.Message, err
		},
	}
	Lookup[setupandconfiguration.AMT_SetupAndConfigurationService]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.SetupAndConfigurationService.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[setupandconfiguration.AMT_SetupAndConfigurationService]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.SetupAndConfigurationService.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.SetupAndConfigurationService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[setupandconfiguration.AMT_SetupAndConfigurationService]["GetUUID"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.SetupAndConfigurationService.GetUuid()
			return *response.Message, err
		},
	}
	// TLSSettingData
	Lookup[tls.AMT_TLSCredentialContext] = make(map[string]Method)
	Lookup[tls.AMT_TLSCredentialContext]["Get"] = Method{
		Execute: func(instanceID string) (client.Message, error) {
			response, err := wsman.AMT.TLSSettingData.Get(instanceID)
			return *response.Message, err
		},
	}
	Lookup[tls.AMT_TLSCredentialContext]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.TLSSettingData.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[tls.AMT_TLSCredentialContext]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.TLSSettingData.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.TLSSettingData.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// UserInitiatedConnectionService
	Lookup[userinitiatedconnection.AMT_UserInitiatedConnectionService] = make(map[string]Method)
	Lookup[userinitiatedconnection.AMT_UserInitiatedConnectionService]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.UserInitiatedConnectionService.Get()
			return *response.Message, err
		},
	}
	Lookup[userinitiatedconnection.AMT_UserInitiatedConnectionService]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.UserInitiatedConnectionService.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[userinitiatedconnection.AMT_UserInitiatedConnectionService]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.UserInitiatedConnectionService.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.UserInitiatedConnectionService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// BootConfigSetting
	Lookup[boot.CIM_BootConfigSetting] = make(map[string]Method)
	Lookup[boot.CIM_BootConfigSetting]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.BootConfigSetting.Get()
			return *response.Message, err
		},
	}
	Lookup[boot.CIM_BootConfigSetting]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.BootConfigSetting.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[boot.CIM_BootConfigSetting]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.BootConfigSetting.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.BootConfigSetting.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// BootService
	Lookup[boot.CIM_BootService] = make(map[string]Method)
	Lookup[boot.CIM_BootService]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.BootService.Get()
			return *response.Message, err
		},
	}
	Lookup[boot.CIM_BootService]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.BootService.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[boot.CIM_BootService]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.BootService.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.BootService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// BootSourceSetting
	Lookup[boot.CIM_BootSourceSetting] = make(map[string]Method)
	Lookup[boot.CIM_BootSourceSetting]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.BootSourceSetting.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[boot.CIM_BootSourceSetting]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.BootSourceSetting.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.BootSourceSetting.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// BIOSElement
	Lookup[bios.CIM_BIOSElement] = make(map[string]Method)
	Lookup[bios.CIM_BIOSElement]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.BIOSElement.Get()
			return *response.Message, err
		},
	}
	Lookup[bios.CIM_BIOSElement]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.BIOSElement.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[bios.CIM_BIOSElement]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.BIOSElement.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.BIOSElement.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// ComputerSystemPackage
	Lookup[computer.CIM_ComputerSystemPackage] = make(map[string]Method)
	Lookup[computer.CIM_ComputerSystemPackage]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.ComputerSystemPackage.Get()
			return *response.Message, err
		},
	}
	Lookup[computer.CIM_ComputerSystemPackage]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.ComputerSystemPackage.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[computer.CIM_ComputerSystemPackage]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.ComputerSystemPackage.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.ComputerSystemPackage.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// Concrete
	Lookup[concrete.CIM_ConcreteDependency] = make(map[string]Method)
	Lookup[concrete.CIM_ConcreteDependency]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.ConcreteDependency.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[concrete.CIM_ConcreteDependency]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.ConcreteDependency.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.ConcreteDependency.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// CredentialContext
	Lookup[credential.CIM_CredentialContext] = make(map[string]Method)
	Lookup[credential.CIM_CredentialContext]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.CredentialContext.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[credential.CIM_CredentialContext]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.CredentialContext.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.CredentialContext.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// IEEE802.1XSettings
	Lookup[ieee8021x.CIM_IEEE8021xSettings] = make(map[string]Method)
	Lookup[ieee8021x.CIM_IEEE8021xSettings]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.IEEE8021xSettings.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[ieee8021x.CIM_IEEE8021xSettings]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.IEEE8021xSettings.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.IEEE8021xSettings.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// KVMRedirectionSAP
	Lookup[kvm.CIM_KVMRedirectionSAP] = make(map[string]Method)
	Lookup[kvm.CIM_KVMRedirectionSAP]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.KVMRedirectionSAP.Get()
			return *response.Message, err
		},
	}
	Lookup[kvm.CIM_KVMRedirectionSAP]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.KVMRedirectionSAP.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[kvm.CIM_KVMRedirectionSAP]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.KVMRedirectionSAP.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.KVMRedirectionSAP.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// MediaAccessDevice
	Lookup[mediaaccess.CIM_MediaAccessDevice] = make(map[string]Method)
	Lookup[mediaaccess.CIM_MediaAccessDevice]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.MediaAccessDevice.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[mediaaccess.CIM_MediaAccessDevice]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.MediaAccessDevice.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.MediaAccessDevice.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// PhysicalCard
	Lookup[card.CIM_Card] = make(map[string]Method)
	Lookup[card.CIM_Card]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.Card.Get()
			return *response.Message, err
		},
	}
	Lookup[card.CIM_Card]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.Card.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[card.CIM_Card]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.Card.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.Card.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// PhysicalChassis
	Lookup[chassis.CIM_Chassis] = make(map[string]Method)
	Lookup[chassis.CIM_Chassis]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.Chassis.Get()
			return *response.Message, err
		},
	}
	Lookup[chassis.CIM_Chassis]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.Chassis.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[chassis.CIM_Chassis]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.Chassis.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.Chassis.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// PhysicalChip
	Lookup[chip.CIM_Chip] = make(map[string]Method)
	Lookup[chip.CIM_Chip]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.Chip.Get()
			return *response.Message, err
		},
	}
	Lookup[chip.CIM_Chip]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.Chip.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[chip.CIM_Chip]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.Chip.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.Chip.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// PhysicalMemory
	Lookup[physical.CIM_PhysicalMemory] = make(map[string]Method)
	Lookup[physical.CIM_PhysicalMemory]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.PhysicalMemory.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[physical.CIM_PhysicalMemory]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.PhysicalMemory.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.PhysicalMemory.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// PhysicalPackage
	Lookup[physical.CIM_PhysicalPackage] = make(map[string]Method)
	Lookup[physical.CIM_PhysicalPackage]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.PhysicalPackage.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[physical.CIM_PhysicalPackage]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.PhysicalPackage.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.PhysicalPackage.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// PhysicalProcessor
	Lookup[processor.CIM_Processor] = make(map[string]Method)
	Lookup[processor.CIM_Processor]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.Processor.Get()
			return *response.Message, err
		},
	}
	Lookup[processor.CIM_Processor]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.Processor.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[processor.CIM_Processor]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.Processor.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.Processor.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// PowerManagementService
	Lookup[power.CIM_PowerManagementService] = make(map[string]Method)
	Lookup[power.CIM_PowerManagementService]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.PowerManagementService.Get()
			return *response.Message, err
		},
	}
	Lookup[power.CIM_PowerManagementService]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.PowerManagementService.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[power.CIM_PowerManagementService]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.PowerManagementService.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.PowerManagementService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[power.CIM_PowerManagementService]["RequestPowerStateChange"] = Method{
		Execute: func(value string) (client.Message, error) {
			powerState, err := strconv.Atoi(value)
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.PowerManagementService.RequestPowerStateChange(power.PowerState(powerState))
			return *response.Message, err
		},
	}
	// ServiceAvailableToElement
	Lookup[service.CIM_ServiceAvailableToElement] = make(map[string]Method)
	Lookup[service.CIM_ServiceAvailableToElement]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.ServiceAvailableToElement.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[service.CIM_ServiceAvailableToElement]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.ServiceAvailableToElement.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.ServiceAvailableToElement.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// SoftwareIdentity
	Lookup[software.CIM_SoftwareIdentity] = make(map[string]Method)
	Lookup[software.CIM_SoftwareIdentity]["Get"] = Method{
		Execute: func(instanceID string) (client.Message, error) {
			selector := software.Selector{
				Name:  "InstanceID",
				Value: "AMTApps",
			}
			response, err := wsman.CIM.SoftwareIdentity.Get(selector)
			return *response.Message, err
		},
	}
	Lookup[software.CIM_SoftwareIdentity]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.SoftwareIdentity.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[software.CIM_SoftwareIdentity]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.SoftwareIdentity.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.SoftwareIdentity.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// SystemPackaging
	Lookup[system.CIM_SystemPackaging] = make(map[string]Method)
	Lookup[system.CIM_SystemPackaging]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.SystemPackaging.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[system.CIM_SystemPackaging]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.SystemPackaging.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.SystemPackaging.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// WiFiEndpointSettings
	Lookup[wifi.CIM_WiFiEndpointSettings] = make(map[string]Method)
	Lookup[wifi.CIM_WiFiEndpointSettings]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.WiFiEndpointSettings.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[wifi.CIM_WiFiEndpointSettings]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.WiFiEndpointSettings.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.WiFiEndpointSettings.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// WiFiPort
	Lookup[wifi.CIM_WiFiPort] = make(map[string]Method)
	Lookup[wifi.CIM_WiFiPort]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.WiFiPort.Get()
			return *response.Message, err
		},
	}
	Lookup[wifi.CIM_WiFiPort]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.CIM.WiFiPort.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[wifi.CIM_WiFiPort]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.CIM.WiFiPort.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.CIM.WiFiPort.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
}
