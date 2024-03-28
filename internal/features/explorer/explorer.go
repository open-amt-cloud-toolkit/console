package explorer

import (
	"encoding/xml"
	"strconv"
	"time"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/alarmclock"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/auditlog"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/authorization"
	amtboot "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/boot"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/environmentdetection"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/ethernetport"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/general"
	amtieee8021x "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/ieee8021x"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/kerberos"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/managementpresence"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/messagelog"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/mps"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/publickey"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/publicprivate"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/redirection"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/remoteaccess"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/setupandconfiguration"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/timesynchronization"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/tls"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/userinitiatedconnection"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/wifiportconfiguration"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/bios"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/boot"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/card"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/chassis"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/chip"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/computer"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/concrete"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/credential"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/ieee8021x"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/kvm"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/mediaaccess"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/physical"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/power"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/processor"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/service"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/software"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/system"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/wifi"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/client"
	ipsalarmclock "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/alarmclock"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/hostbasedsetup"
	ipsieee8021x "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/ieee8021x"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/optin"
)

func GetSupportedWsmanClasses(className string) []Class {
	var ClassList = []Class{
		{Name: alarmclock.AMT_AlarmClockService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "AddAlarm"}}},
		{Name: auditlog.AMT_AuditLog, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "ReadRecords"}}},
		{Name: authorization.AMT_AuthorizationService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "EnumerateUserAclEntries"}, {Name: "GetAclEnabledState"}, {Name: "GetAdminAclEntry"}, {Name: "GetAdminAclEntryStatus"}, {Name: "GetAdminNetAclEntryStatus"}, {Name: "GetUserAclEntryEx"}, {Name: "RemoveUserAclEntry"}, {Name: "SetAclEnabledState"}, {Name: "SetAdminAclEntryEx"}}},
		{Name: amtboot.AMT_BootCapabilities, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: amtboot.AMT_BootSettingData, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Put"}}},
		{Name: environmentdetection.AMT_EnvironmentDetectionSettingData, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Put"}}},
		{Name: ethernetport.AMT_EthernetPortSettings, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Put"}}},
		{Name: general.AMT_GeneralSettings, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: amtieee8021x.AMT_IEEE8021xCredentialContext, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: amtieee8021x.AMT_IEEE8021xProfile, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Put"}}},
		{Name: kerberos.AMT_KerberosSettingData, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "GetCredentialCacheState"}, {Name: "SetCredentialCacheState"}}},
		{Name: managementpresence.AMT_ManagementPresenceRemoteSAP, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Delete"}}},
		{Name: messagelog.AMT_MessageLog, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "GetRecords"}, {Name: "PositionToFirstRecord"}}},
		{Name: mps.AMT_MPSUsernamePassword, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Put"}}},
		{Name: publickey.AMT_PublicKeyCertificate, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Put"}, {Name: "Delete"}}},
		{Name: publickey.AMT_PublicKeyManagementService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Delete"}, {Name: "AddCertificate"}, {Name: "AddTrustedRootCertificate"}, {Name: "GenerateKeyPair"}, {Name: "GeneratePKCS10RequestEx"}, {Name: "AddKey"}}},
		{Name: publicprivate.AMT_PublicPrivateKeyPair, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Delete"}}},
		{Name: redirection.AMT_RedirectionService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Put"}, {Name: "RequestStateChange"}}},
		{Name: remoteaccess.AMT_RemoteAccessPolicyAppliesToMPS, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Put"}, {Name: "Delete"}}},
		{Name: remoteaccess.AMT_RemoteAccessPolicyRule, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Delete"}}},
		{Name: remoteaccess.AMT_RemoteAccessService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "AddMPS"}, {Name: "AddRemoteAccessPolicyRule"}}},
		{Name: setupandconfiguration.AMT_SetupAndConfigurationService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "GetUUID"}}},
		{Name: timesynchronization.AMT_TimeSynchronizationService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "SetHighAccuracyTimeSynch"}, {Name: "GetLowAccuracyTimeSynch"}}},
		{Name: tls.AMT_TLSCredentialContext, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: tls.AMT_TLSProtocolEndpointCollection, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: tls.AMT_TLSSettingData, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: userinitiatedconnection.AMT_UserInitiatedConnectionService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "RequestStateChange"}}},
		{Name: wifiportconfiguration.AMT_WiFiPortConfigurationService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Put"}, {Name: "AddWiFiSettings"}}},
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
		{Name: ipsalarmclock.IPS_AlarmClockOccurrence, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Delete"}}},
		{Name: hostbasedsetup.IPS_HostBasedSetupService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Setup"}}},
		{Name: ipsieee8021x.IPS_8021xCredentialContext, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}}},
		{Name: ipsieee8021x.IPS_IEEE8021xSettings, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "Put"}, {Name: "SetCertificates"}}},
		{Name: optin.IPS_OptInService, MethodList: []Method{{Name: "Get"}, {Name: "Enumerate"}, {Name: "Pull"}, {Name: "SendOptInCode"}, {Name: "StartOptIn"}, {Name: "CancelOptIn"}}},
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
	Lookup[alarmclock.AMT_AlarmClockService]["AddAlarm"] = Method{
		Execute: func(value string) (client.Message, error) {
			alarm := alarmclock.AlarmClockOccurrence{
				InstanceID:         "test",
				StartTime:          time.Now().Add(time.Hour + 1),
				Interval:           0,
				DeleteOnCompletion: true,
			}
			response, err := wsman.AMT.AlarmClockService.AddAlarm(alarm)
			return *response.Message, err
		},
	}
	// Audit Log
	Lookup[auditlog.AMT_AuditLog] = make(map[string]Method)
	Lookup[auditlog.AMT_AuditLog]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AuditLog.Get()
			return *response.Message, err
		},
	}
	Lookup[auditlog.AMT_AuditLog]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AuditLog.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[auditlog.AMT_AuditLog]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.AuditLog.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.AuditLog.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[auditlog.AMT_AuditLog]["ReadRecords"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AuditLog.ReadRecords(1)
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
	Lookup[authorization.AMT_AuthorizationService]["EnumerateUserAclEntries"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AuthorizationService.EnumerateUserAclEntries(1)
			return *response.Message, err
		},
	}
	Lookup[authorization.AMT_AuthorizationService]["GetAclEnabledState"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AuthorizationService.GetAclEnabledState(1)
			return *response.Message, err
		},
	}
	Lookup[authorization.AMT_AuthorizationService]["GetAdminAclEntry"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AuthorizationService.GetAdminAclEntry()
			return *response.Message, err
		},
	}
	Lookup[authorization.AMT_AuthorizationService]["GetAdminAclEntryStatus"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AuthorizationService.GetAdminAclEntryStatus()
			return *response.Message, err
		},
	}
	Lookup[authorization.AMT_AuthorizationService]["GetAdminNetAclEntryStatus"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AuthorizationService.GetAdminNetAclEntryStatus()
			return *response.Message, err
		},
	}
	Lookup[authorization.AMT_AuthorizationService]["GetUserAclEntryEx"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AuthorizationService.GetUserAclEntryEx(1)
			return *response.Message, err
		},
	}
	Lookup[authorization.AMT_AuthorizationService]["RemoveUserAclEntry"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AuthorizationService.RemoveUserAclEntry(1)
			return *response.Message, err
		},
	}
	Lookup[authorization.AMT_AuthorizationService]["SetAclEnabledState"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AuthorizationService.SetAclEnabledState(1, true)
			return *response.Message, err
		},
	}
	Lookup[authorization.AMT_AuthorizationService]["SetAdminAclEntryEx"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.AuthorizationService.SetAdminAclEntryEx("admin", "P@ssw0rd")
			return *response.Message, err
		},
	}
	// BootSettingData
	Lookup[amtboot.AMT_BootSettingData] = make(map[string]Method)
	Lookup[amtboot.AMT_BootSettingData]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.BootSettingData.Get()
			return *response.Message, err
		},
	}
	Lookup[amtboot.AMT_BootSettingData]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.BootSettingData.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[amtboot.AMT_BootSettingData]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.BootSettingData.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.BootSettingData.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[amtboot.AMT_BootSettingData]["Put"] = Method{
		Execute: func(value string) (client.Message, error) {
			bootSettingData := amtboot.BootSettingDataRequest{
				H:                      "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_BootSettingData",
				BIOSPause:              false,
				BIOSSetup:              false,
				BootMediaIndex:         0,
				ConfigurationDataReset: false,
				ElementName:            "Intel(r) AMT Boot Configuration Settings",
				EnforceSecureBoot:      false,
				FirmwareVerbosity:      0,
				ForcedProgressEvents:   false,
				IDERBootDevice:         0,
				InstanceID:             "Intel(r) AMT:BootSettingData 0",
				LockKeyboard:           false,
				LockPowerButton:        false,
				LockResetButton:        false,
				LockSleepButton:        false,
				RSEPassword:            "",
				ReflashBIOS:            false,
				SecureErase:            false,
				UseIDER:                false,
				UseSOL:                 false,
				UseSafeMode:            false,
				UserPasswordBypass:     false,
			}
			response, err := wsman.AMT.BootSettingData.Put(bootSettingData)
			return *response.Message, err
		},
	}
	// BootCapabilities
	Lookup[amtboot.AMT_BootCapabilities] = make(map[string]Method)
	Lookup[amtboot.AMT_BootCapabilities]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.BootCapabilities.Get()
			return *response.Message, err
		},
	}
	Lookup[amtboot.AMT_BootCapabilities]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.BootCapabilities.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[amtboot.AMT_BootCapabilities]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.BootCapabilities.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.BootCapabilities.Pull(er.Body.EnumerateResponse.EnumerationContext)
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
	Lookup[environmentdetection.AMT_EnvironmentDetectionSettingData]["Put"] = Method{
		Execute: func(value string) (client.Message, error) {
			edsd := environmentdetection.EnvironmentDetectionSettingDataRequest{
				H:                  "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_EnvironmentDetectionSettingData",
				DetectionAlgorithm: 0,
				ElementName:        "Intel(r) AMT Environment Detection Settings",
				InstanceID:         "Intel(r) AMT Environment Detection Settings",
				DetectionStrings:   []string{},
			}
			response, err := wsman.AMT.EnvironmentDetectionSettingData.Put(edsd)
			return *response.Message, err
		},
	}
	// EthernetPortSetting
	Lookup[ethernetport.AMT_EthernetPortSettings] = make(map[string]Method)
	Lookup[ethernetport.AMT_EthernetPortSettings]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.EthernetPortSettings.Get("Intel(r) AMT Ethernet Port Settings 0")
			return *response.Message, err
		},
	}
	Lookup[ethernetport.AMT_EthernetPortSettings]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.EthernetPortSettings.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[ethernetport.AMT_EthernetPortSettings]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.EthernetPortSettings.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.EthernetPortSettings.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[ethernetport.AMT_EthernetPortSettings]["Put"] = Method{
		Execute: func(value string) (client.Message, error) {
			ethernetPortSettings := ethernetport.SettingsRequest{
				H:              "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_EthernetPortSettings",
				DHCPEnabled:    true,
				ElementName:    "Intel(r) AMT Ethernet Port Settings",
				InstanceID:     "Intel(r) AMT Ethernet Port Settings 0",
				IpSyncEnabled:  true,
				SharedMAC:      true,
				SharedStaticIp: true,
			}
			response, err := wsman.AMT.EthernetPortSettings.Put(ethernetPortSettings.InstanceID, ethernetPortSettings)
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
	// IEEE8021XCredentialContext
	Lookup[amtieee8021x.AMT_IEEE8021xCredentialContext] = make(map[string]Method)
	Lookup[amtieee8021x.AMT_IEEE8021xCredentialContext]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.IEEE8021xCredentialContext.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[amtieee8021x.AMT_IEEE8021xCredentialContext]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.IEEE8021xCredentialContext.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.IEEE8021xCredentialContext.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// IEEE8021XProfile
	Lookup[amtieee8021x.AMT_IEEE8021xProfile] = make(map[string]Method)
	Lookup[amtieee8021x.AMT_IEEE8021xProfile]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.IEEE8021xProfile.Get()
			return *response.Message, err
		},
	}
	Lookup[amtieee8021x.AMT_IEEE8021xProfile]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.IEEE8021xProfile.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[amtieee8021x.AMT_IEEE8021xProfile]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.IEEE8021xProfile.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.IEEE8021xProfile.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[amtieee8021x.AMT_IEEE8021xProfile]["Put"] = Method{
		Execute: func(value string) (client.Message, error) {
			profileRequest := amtieee8021x.ProfileRequest{
				H:           "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_8021XProfile",
				ActiveInS0:  false,
				ElementName: "Intel(r) AMT 802.1x Profile",
				Enabled:     false,
				InstanceID:  "Intel(r) AMT 802.1x Profile 0",
			}
			response, err := wsman.AMT.IEEE8021xProfile.Put(profileRequest)
			return *response.Message, err
		},
	}
	// KerberosSettingData
	Lookup[kerberos.AMT_KerberosSettingData] = make(map[string]Method)
	Lookup[kerberos.AMT_KerberosSettingData]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.KerberosSettingData.Get()
			return *response.Message, err
		},
	}
	Lookup[kerberos.AMT_KerberosSettingData]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.KerberosSettingData.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[kerberos.AMT_KerberosSettingData]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.KerberosSettingData.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.KerberosSettingData.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[kerberos.AMT_KerberosSettingData]["GetCredentialCacheState"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.KerberosSettingData.GetCredentialCacheState()
			return *response.Message, err
		},
	}
	Lookup[kerberos.AMT_KerberosSettingData]["SetCredentialCacheState"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.KerberosSettingData.SetCredentialCacheState(false)
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
	Lookup[managementpresence.AMT_ManagementPresenceRemoteSAP]["Delete"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.ManagementPresenceRemoteSAP.Delete("Intel(r) AMT:Management Presence Server 0")
			return *response.Message, err
		},
	}
	// MessageLog
	Lookup[messagelog.AMT_MessageLog] = make(map[string]Method)
	Lookup[messagelog.AMT_MessageLog]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.MessageLog.Get()
			return *response.Message, err
		},
	}
	Lookup[messagelog.AMT_MessageLog]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.MessageLog.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[messagelog.AMT_MessageLog]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.MessageLog.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.MessageLog.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[messagelog.AMT_MessageLog]["GetRecords"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.MessageLog.GetRecords(1)
			return *response.Message, err
		},
	}
	Lookup[messagelog.AMT_MessageLog]["PositionToFirstRecord"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.MessageLog.PositionToFirstRecord()
			return *response.Message, err
		},
	}
	// MPS
	Lookup[mps.AMT_MPSUsernamePassword] = make(map[string]Method)
	Lookup[mps.AMT_MPSUsernamePassword]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.MPSUsernamePassword.Get()
			return *response.Message, err
		},
	}
	Lookup[mps.AMT_MPSUsernamePassword]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.MPSUsernamePassword.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[mps.AMT_MPSUsernamePassword]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.MPSUsernamePassword.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.MPSUsernamePassword.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[mps.AMT_MPSUsernamePassword]["Put"] = Method{
		Execute: func(value string) (client.Message, error) {
			mpsUsernamePassword := mps.MPSUsernamePasswordRequest{
				H:          "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_MPSUsernamePassword",
				InstanceID: "Intel(r) AMT:MPS Username Password 0",
				RemoteID:   "test",
				Secret:     "P@ssw0rd",
			}
			response, err := wsman.AMT.MPSUsernamePassword.Put(mpsUsernamePassword)
			return *response.Message, err
		},
	}
	// PublicKeyCertificate
	Lookup[publickey.AMT_PublicKeyCertificate] = make(map[string]Method)
	Lookup[publickey.AMT_PublicKeyCertificate]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.PublicKeyCertificate.Get("Intel(r) AMT Certificate: Handle: 0")
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
			response, err := wsman.AMT.PublicKeyCertificate.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[publickey.AMT_PublicKeyCertificate]["Delete"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.PublicKeyCertificate.Delete("Intel(r) AMT Certificate: Handle: 0")
			return *response.Message, err
		},
	}
	Lookup[publickey.AMT_PublicKeyCertificate]["Put"] = Method{
		Execute: func(value string) (client.Message, error) {
			X509Certificate := "MIIEOzCCAqOgAwIBAgIDAZMjMA0GCSqGSIb3DQEBDAUAMD0xFzAVBgNVBAMTDk1QU1Jvb3QtMGFmMWQ1MRAwDgYDVQQKEwd1bmtub3duMRAwDgYDVQQGEwd1bmtub3duMCAXDTIyMDkyNDEwNDUwOFoYDzIwNTMwOTI0MTA0NTA4WjA9MRcwFQYDVQQDEw5NUFNSb290LTBhZjFkNTEQMA4GA1UEChMHdW5rbm93bjEQMA4GA1UEBhMHdW5rbm93bjCCAaIwDQYJKoZIhvcNAQEBBQADggGPADCCAYoCggGBALz/oJNyWXlClSlteAieC8Uyd4A+tbn8b45k6LKiImhDmdz/xFo9xe0C9GNf7b42KVpg5WoH/sPhoClR9Tv5i1LnilT1SUir42fcm2NEV9dRcLsPd/RAQfz8u0D4zb3blnxE8isqzriNpG7kac35UidSr5ym8TZ3IwXx6JJuncGgfB0DFZADC/+dA74n3coykvWBYqLr6RI5pkAxvulkRlCsatJTJrvMUYJ51GI28jV56mIAc89sLrHqiSKCZBH9AcUrnZ/cB6ST/IikXpxy5wXBIvWT3VKVq75T/uIoCBEp5TLEn1EOYGqBBOCSQgmtmX7eVaB0s1+ppPW9w9a2zS45cHAtQ7tYvkkPv2dRhSzZdlk6HRXDP5wsF0aiflZCgbrjkq0SFC4e3Lo7XQX3FTNb0SOTZVTydupoMKkgJQTNlcosdu1ZzaIBl3eSkKkJZz2rUTssZC5tn9vcDd5vy3BzcGh5pvkgfAgN1sydqG7Ke1qCkNEzm11B/BsevatjjwIDAQABo0IwQDAMBgNVHRMEBTADAQH/MBEGCWCGSAGG+EIBAQQEAwIABzAdBgNVHQ4EFgQUCvHVQqerCid99eLApuLky9x6H5owDQYJKoZIhvcNAQEMBQADggGBAIzOyGV0hzsmH2biJlzwTZaHMxqS7boTFMkHw+KvzsI201tHqVmCoiQ8EHErBGLSoDOTDRgOUGOCA5XU5ie9OWupAGqKBSwIyAhmJMOzrzC4Gwpu8K1msoFJH30kx/V9purpbS3BRj0xfYXLa6IczbTg3E5IfTnZRJ9YuUtKQfI0P9c5U9CoKtddKn4+lRvOjFDoYfQGCJ7go3xjNCcGCVCjfkUhAVdbQ21DCRr6/YCZDWmjzZpL0p7UKF8roTiNuL/Z7gIXxch5HOmEWHY9uQ6K2MntuxAu0aK/mSD2kwmt/ECongdEGfUvhULLoPRQlQ2LnzcUQEgMECGQR5Yfy9jT0E8zdWDpc2tgVioNu6rEYKgp/GhG+sv7jv58pW82FRAV9xXtftW9+XDugC8tBJ6JHn0Q2v0QAflD2CEQVhWAY8bAqrbfTGUsaLfGL6kxV/qqssoMgLR8WhQ96T5le/4XGhQpbCHWIlctD6MwbrsunIAeQKp1Sc3DosY7DLq1MQ=="
			response, err := wsman.AMT.PublicKeyCertificate.Put("Intel(r) AMT Certificate: Handle: 0", X509Certificate)
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
	Lookup[publickey.AMT_PublicKeyManagementService]["Delete"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.PublicKeyManagementService.Delete("Intel(r) AMT Certificate: Handle: 1")
			return *response.Message, err
		},
	}
	Lookup[publickey.AMT_PublicKeyManagementService]["AddCertificate"] = Method{
		Execute: func(value string) (client.Message, error) {
			certificateBlob := "MIIEOzCCAqOgAwIBAgIDAZMjMA0GCSqGSIb3DQEBDAUAMD0xFzAVBgNVBAMTDk1QU1Jvb3QtMGFmMWQ1MRAwDgYDVQQKEwd1bmtub3duMRAwDgYDVQQGEwd1bmtub3duMCAXDTIyMDkyNDEwNDUwOFoYDzIwNTMwOTI0MTA0NTA4WjA9MRcwFQYDVQQDEw5NUFNSb290LTBhZjFkNTEQMA4GA1UEChMHdW5rbm93bjEQMA4GA1UEBhMHdW5rbm93bjCCAaIwDQYJKoZIhvcNAQEBBQADggGPADCCAYoCggGBALz/oJNyWXlClSlteAieC8Uyd4A+tbn8b45k6LKiImhDmdz/xFo9xe0C9GNf7b42KVpg5WoH/sPhoClR9Tv5i1LnilT1SUir42fcm2NEV9dRcLsPd/RAQfz8u0D4zb3blnxE8isqzriNpG7kac35UidSr5ym8TZ3IwXx6JJuncGgfB0DFZADC/+dA74n3coykvWBYqLr6RI5pkAxvulkRlCsatJTJrvMUYJ51GI28jV56mIAc89sLrHqiSKCZBH9AcUrnZ/cB6ST/IikXpxy5wXBIvWT3VKVq75T/uIoCBEp5TLEn1EOYGqBBOCSQgmtmX7eVaB0s1+ppPW9w9a2zS45cHAtQ7tYvkkPv2dRhSzZdlk6HRXDP5wsF0aiflZCgbrjkq0SFC4e3Lo7XQX3FTNb0SOTZVTydupoMKkgJQTNlcosdu1ZzaIBl3eSkKkJZz2rUTssZC5tn9vcDd5vy3BzcGh5pvkgfAgN1sydqG7Ke1qCkNEzm11B/BsevatjjwIDAQABo0IwQDAMBgNVHRMEBTADAQH/MBEGCWCGSAGG+EIBAQQEAwIABzAdBgNVHQ4EFgQUCvHVQqerCid99eLApuLky9x6H5owDQYJKoZIhvcNAQEMBQADggGBAIzOyGV0hzsmH2biJlzwTZaHMxqS7boTFMkHw+KvzsI201tHqVmCoiQ8EHErBGLSoDOTDRgOUGOCA5XU5ie9OWupAGqKBSwIyAhmJMOzrzC4Gwpu8K1msoFJH30kx/V9purpbS3BRj0xfYXLa6IczbTg3E5IfTnZRJ9YuUtKQfI0P9c5U9CoKtddKn4+lRvOjFDoYfQGCJ7go3xjNCcGCVCjfkUhAVdbQ21DCRr6/YCZDWmjzZpL0p7UKF8roTiNuL/Z7gIXxch5HOmEWHY9uQ6K2MntuxAu0aK/mSD2kwmt/ECongdEGfUvhULLoPRQlQ2LnzcUQEgMECGQR5Yfy9jT0E8zdWDpc2tgVioNu6rEYKgp/GhG+sv7jv58pW82FRAV9xXtftW9+XDugC8tBJ6JHn0Q2v0QAflD2CEQVhWAY8bAqrbfTGUsaLfGL6kxV/qqssoMgLR8Whq96T5le/4XGhQpbCHWIlctD6MwbrsunIAeQKp1Sc3DosY7DLq1MQ=="
			response, err := wsman.AMT.PublicKeyManagementService.AddCertificate(certificateBlob)
			return *response.Message, err
		},
	}
	Lookup[publickey.AMT_PublicKeyManagementService]["AddTrustedRootCertificate"] = Method{
		Execute: func(value string) (client.Message, error) {
			certificateBlob := "MIIEOzCCAqOgAwIBAgIDAZMjMA0GCSqGSIb3DQEBDAUAMD0xFzAVBgNVBAMTDk1QU1Jvb3QtMGFmMWQ1MRAwDgYDVQQKEwd1bmtub3duMRAwDgYDVQQGEwd1bmtub3duMCAXDTIyMDkyNDEwNDUwOFoYDzIwNTMwOTI0MTA0NTA4WjA9MRcwFQYDVQQDEw5NUFNSb290LTBhZjFkNTEQMA4GA1UEChMHdW5rbm93bjEQMA4GA1UEBhMHdW5rbm93bjCCAaIwDQYJKoZIhvcNAQEBBQADggGPADCCAYoCggGBALz/oJNyWXlClSlteAieC8Uyd4A+tbn8b45k6LKiImhDmdz/xFo9xe0C9GNf7b42KVpg5WoH/sPhoClR9Tv5i1LnilT1SUir42fcm2NEV9dRcLsPd/RAQfz8u0D4zb3blnxE8isqzriNpG7kac35UidSr5ym8TZ3IwXx6JJuncGgfB0DFBADC/+dA74n3coykvWBYqLr6RI5pkAxvulkRlCsatJTJrvMUYJ51GI28jV56mIAc89sLrHqiSKCZBH9AcUrnZ/cB6ST/IikXpxy5wXBIvWT3VKVq75T/uIoCBEp5TLEn1EOYGqBBOCSQgmtmX7eVaB0s1+ppPW9w9a2zS45cHAtQ7tYvkkPv2dRhSzZdlk6HRXDP5wsF0aiflZCgbrjkq0SFC4e3Lo7XQX3FTNb0SOTZVTydupoMKkgJQTNlcosdu1ZzaIBl3eSkKkJZz2rUTssZC5tn9vcDd5vy3BzcGh5pvkgfAgN1sydqG7Ke1qCkNEzm11B/BsevatjjwIDAQABo0IwQDAMBgNVHRMEBTADAQH/MBEGCWCGSAGG+EIBAQQEAwIABzAdBgNVHQ4EFgQUCvHVQqerCid99eLApuLky9x6H5owDQYJKoZIhvcNAQEMBQADggGBAIzOyGV0hzsmH2biJlzwTZaHMxqS7boTFMkHw+KvzsI201tHqVmCoiQ8EHErBGLSoDOTDRgOUGOCA5XU5ie9OWupAGqKBSwIyAhmJMOzrzC4Gwpu8K1msoFJH30kx/V9purpbS3BRj0xfYXLa6IczbTg3E5IfTnZRJ9YuUtKQfI0P9c5U9CoKtddKn4+lRvOjFDoYfQGCJ7go3xjNCcGCVCjfkUhAVdbQ21DCRr6/YCZDWmjzZpL0p7UKF8roTiNuL/Z7gIXxch5HOmEWHY9uQ6K2MntuxAu0aK/mSD2kwmt/ECongdEGfUvhULLoPRQlQ2LnzcUQEgMECGQR5Yfy9jT0E8zdWDpc2tgVioNu6rEYKgp/GhG+sv7jv58pW82FRAV9xXtftW9+XDugC8tBJ6JHn0Q2v0QAflD2CEQVhWAY8bAqrbfTGUsaLfGL6kxV/qqssoMgLR8Whq96T5le/4XGhQpbCHWIlctD6MwbrsunIAeQKp1Sc3DosY7DLq1MQ=="
			response, err := wsman.AMT.PublicKeyManagementService.AddTrustedRootCertificate(certificateBlob)
			return *response.Message, err
		},
	}
	Lookup[publickey.AMT_PublicKeyManagementService]["GenerateKeyPair"] = Method{
		Execute: func(value string) (client.Message, error) {
			var keyAlgo publickey.KeyAlgorithm = publickey.RSA
			var keyLength publickey.KeyLength = publickey.KeyLength2048
			response, err := wsman.AMT.PublicKeyManagementService.GenerateKeyPair(keyAlgo, keyLength)
			return *response.Message, err
		},
	}
	Lookup[publickey.AMT_PublicKeyManagementService]["GeneratePKCS10RequestEx"] = Method{
		Execute: func(value string) (client.Message, error) {
			keyPair := "Intel(r) AMT Key: Handle: 1"
			nullSignedCert := "MIIC+jCCAeKgAwIBAgIJAKSnJkwyL9w3MA0GCSqGSIb3DQEBCwUAMBcxFTATBgNVBAMMDHRlc3QgY2VydGlmaWNhdGUwHhcNMTkxMjA1MDMxODMzWhcNMjAxMjA0MDMxODMzWjAhMR8wHQYDVQQDDBZUZXN0IFRlc3QgQ2VydGlmaWNhdGUwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDJaF5HXldh0iHFlsVGQFQKs1wH+MV65tNgY6x3FZND9/DXkZJQoqqfCWk3iZm5zD9OTKxNsv4ObuKs7+2BCnmEzQVKW2t7NE9Dd7evknMdYDWTNvmF3GJxIJNc9TD1I3z5Z2v+8xjYiw2S6J9sFfcWpXjXf6yLd5vhUE9YRcFJ4TAmFdE2BxP8aPEUBwZK5Sjkm0cgyJzI9SJdqCh5VO3bHc6AaP60qEeVxL2jHh8U+W1BsBh7ckY9R4EVrylL6uGxFfkxLdzT/c2Zw6NzL+wMXW47ewxVrv5VEf8hJ33zTteJ24T9EwN4XGvq4I1Z32ynqEVqh2Y9eBS6+M/JAgMBAAGjUDBOMB0GA1UdDgQWBBTRCG/W0nw9nyWVv8uAF15c8V6uBzAfBgNVHSMEGDAWgBTRCG/W0nw9nyWVv8uAF15c8V6uBzAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQAVNR5aK4LNY4ptCcrBs/KcTHvIzJbh5owc0aMRR6nQQx0sGJlA1I8OizREl6FzexnMnLZ7Nq1E7tQD5Z43fVPRlym9tXfw6kIKG+RaofF7Q7DC0JuvCnvoThBFpgQFWzFNm9d3zi+PkzokYm7dH68qz5zbpziSYG+qeUvEMccZfb0l5Kpgv1qu/3hnybAC6ILbtMS8Ku9lONj+AJad0vSg1ddtE5hVqT0I3cX7eZEl73Q7z7CRNOE4Neyn0Wn4b8TOXXJ1A8TvWz4Kx86p2dcPUMGJGtnAx5kTX6je2lHVDnY3+3mQQ8wGxwzKhgoAkUYkva7r7nS5Hv5xuFjyEKSf"
			var signingAlgo publickey.SigningAlgorithm = publickey.SHA256RSA
			response, err := wsman.AMT.PublicKeyManagementService.GeneratePKCS10RequestEx(keyPair, nullSignedCert, signingAlgo)
			return *response.Message, err
		},
	}
	Lookup[publickey.AMT_PublicKeyManagementService]["AddKey"] = Method{
		Execute: func(value string) (client.Message, error) {
			keyBlob := "MIIEowIBAAKCAQEAz2OvfbD0fBj0OM6PdqYcHnoLOvzJ3B+QRvKByP9L0vZG6XmbRA6Wp9qUE8PhRJUab8Rg4DwKrA+gx88HVLELnQk+ZahPivfp6P/T9DFsoFtyc1ZJ7zLpnbS3kxLBhpxlVgBh0ozHvK8i8hVLFlhh7ZvKWh6pdyMblJBcHhOl54uZVPPB+HTufcUQlOfNt76KZMeoOQ=="
			response, err := wsman.AMT.PublicKeyManagementService.AddKey(keyBlob)
			return *response.Message, err
		},
	}
	// PublicPrivate
	Lookup[publicprivate.AMT_PublicPrivateKeyPair] = make(map[string]Method)
	Lookup[publicprivate.AMT_PublicPrivateKeyPair]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.PublicPrivateKeyPair.Get("Intel(r) AMT Key: Handle: 0")
			return *response.Message, err
		},
	}
	Lookup[publicprivate.AMT_PublicPrivateKeyPair]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.PublicPrivateKeyPair.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[publicprivate.AMT_PublicPrivateKeyPair]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.PublicPrivateKeyPair.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.PublicPrivateKeyPair.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[publicprivate.AMT_PublicPrivateKeyPair]["Delete"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.PublicPrivateKeyPair.Delete("0")
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
	Lookup[redirection.AMT_RedirectionService]["Put"] = Method{
		Execute: func(value string) (client.Message, error) {
			redirectionService := redirection.RedirectionRequest{}
			redirectionService.CreationClassName = "AMT_RedirectionService"
			redirectionService.ElementName = "Intel(r) AMT Redirection Service"
			redirectionService.EnabledState = redirection.IDERAndSOLAreDisabled
			redirectionService.ListenerEnabled = true
			redirectionService.Name = "Intel(r) AMT Redirection Service"
			redirectionService.SystemCreationClassName = "CIM_ComputerSystem"
			redirectionService.SystemName = "Intel(r) AMT"
			response, err := wsman.AMT.RedirectionService.Put(redirectionService)
			return *response.Message, err
		},
	}
	Lookup[redirection.AMT_RedirectionService]["RequestStateChange"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.RedirectionService.RequestStateChange(redirection.IDERAndSOLAreEnabled)
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
			response, err := wsman.AMT.RemoteAccessPolicyAppliesToMPS.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[remoteaccess.AMT_RemoteAccessPolicyAppliesToMPS]["Put"] = Method{
		Execute: func(value string) (client.Message, error) {
			remoteAccessPolicy := remoteaccess.RemoteAccessPolicyAppliesToMPSRequest{
				XMLName: xml.Name{Space: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_RemoteAccessPolicyAppliesToMPS", Local: "AMT_RemoteAccessPolicyAppliesToMPS"},
				H:       "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_RemoteAccessPolicyAppliesToMPS",
				ManagedElement: remoteaccess.ManagedElement{
					Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
					B:       "http://schemas.xmlsoap.org/ws/2004/08/addressing",
					ReferenceParameters: remoteaccess.ReferenceParameters{
						ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_ManagementPresenceRemoteSAP",
						C:           "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
						SelectorSet: remoteaccess.SelectorSet{
							Selectors: []remoteaccess.Selector{
								{
									Name: "CreationClassName",
									Text: "AMT_ManagementPresenceRemoteSAP",
								},
								{
									Name: "Name",
									Text: "Intel(r) AMT:Management Presence Server 0",
								},
								{
									Name: "SystemCreationClassName",
									Text: "CIM_ComputerSystem",
								},
								{
									Name: "SystemName",
									Text: "Intel(r) AMT",
								},
							},
						},
					},
				},
				OrderOfAccess: 0,
				MPSType:       remoteaccess.BothMPS,
				PolicySet: remoteaccess.PolicySet{
					Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
					B:       "http://schemas.xmlsoap.org/ws/2004/08/addressing",
					ReferenceParameters: remoteaccess.ReferenceParameters{
						ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_RemoteAccessPolicyRule",
						C:           "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
						SelectorSet: remoteaccess.SelectorSet{
							Selectors: []remoteaccess.Selector{
								{
									Name: "CreationClassName",
									Text: "AMT_RemoteAccessPolicyRule",
								},
								{
									Name: "PolicyRuleName",
									Text: "Periodic",
								},
								{
									Name: "SystemCreationClassName",
									Text: "CIM_ComputerSystem",
								},
								{
									Name: "SystemName",
									Text: "Intel(r) AMT",
								},
							},
						},
					},
				},
			}
			response, err := wsman.AMT.RemoteAccessPolicyAppliesToMPS.Put(&remoteAccessPolicy)
			return *response.Message, err
		},
	}
	Lookup[remoteaccess.AMT_RemoteAccessPolicyAppliesToMPS]["Delete"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.RemoteAccessPolicyAppliesToMPS.Delete("Name")
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
			response, err := wsman.AMT.RemoteAccessPolicyRule.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[remoteaccess.AMT_RemoteAccessPolicyRule]["Delete"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.RemoteAccessPolicyRule.Delete("Periodic")
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
	Lookup[remoteaccess.AMT_RemoteAccessService]["AddMPS"] = Method{
		Execute: func(value string) (client.Message, error) {
			mps := remoteaccess.AddMpServerRequest{
				H:          "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_RemoteAccessService",
				AccessInfo: "192.168.0.25",
				InfoFormat: remoteaccess.IPv4Address,
				Port:       4433,
				AuthMethod: remoteaccess.UsernamePasswordAuthentication,
				Username:   "test",
				Password:   "P@ssw0rd",
				CommonName: "192.168.0.25",
			}
			response, err := wsman.AMT.RemoteAccessService.AddMPS(mps)
			return *response.Message, err
		},
	}
	Lookup[remoteaccess.AMT_RemoteAccessService]["AddRemoteAccessPolicyRule"] = Method{
		Execute: func(value string) (client.Message, error) {
			remoteAccessPolicyRule := remoteaccess.RemoteAccessPolicyRuleRequest{
				Trigger:        remoteaccess.Periodic,
				TunnelLifeTime: 0,
				ExtendedData:   "AAAAAAAAABk=",
			}
			response, err := wsman.AMT.RemoteAccessService.AddRemoteAccessPolicyRule(remoteAccessPolicyRule, "Intel(r) AMT:Management Presence Server 0")
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
	// TimeSynchronization
	Lookup[timesynchronization.AMT_TimeSynchronizationService] = make(map[string]Method)
	Lookup[timesynchronization.AMT_TimeSynchronizationService]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.TimeSynchronizationService.Get()
			return *response.Message, err
		},
	}
	Lookup[timesynchronization.AMT_TimeSynchronizationService]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.TimeSynchronizationService.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[timesynchronization.AMT_TimeSynchronizationService]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.TimeSynchronizationService.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.TimeSynchronizationService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[timesynchronization.AMT_TimeSynchronizationService]["GetLowAccuracyTimeSynch"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.TimeSynchronizationService.GetLowAccuracyTimeSynch()
			return *response.Message, err
		},
	}
	Lookup[timesynchronization.AMT_TimeSynchronizationService]["SetHighAccuracyTimeSynch"] = Method{
		Execute: func(value string) (client.Message, error) {
			ta2 := time.Now().Unix()
			glatResponse, err := wsman.AMT.TimeSynchronizationService.GetLowAccuracyTimeSynch()
			if err != nil {
				return client.Message{}, err
			}
			ta1 := time.Now().Unix()
			response, err := wsman.AMT.TimeSynchronizationService.SetHighAccuracyTimeSynch(glatResponse.Body.GetLowAccuracyTimeSynchResponse.Ta0, int64(ta2), ta1)
			return *response.Message, err
		},
	}
	// TLSCredentialContext
	Lookup[tls.AMT_TLSCredentialContext] = make(map[string]Method)
	Lookup[tls.AMT_TLSCredentialContext]["Get"] = Method{
		Execute: func(instanceID string) (client.Message, error) {
			response, err := wsman.AMT.TLSCredentialContext.Get()
			return *response.Message, err
		},
	}
	Lookup[tls.AMT_TLSCredentialContext]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.TLSCredentialContext.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[tls.AMT_TLSCredentialContext]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.TLSCredentialContext.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.TLSCredentialContext.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// Lookup[tls.AMT_TLSCredentialContext]["Delete"] = Method{
	// 	Execute: func(value string) (client.Message, error) {
	// 		response, err := wsman.AMT.TLSCredentialContext.Delete("Handle 1")
	// 		return *response.Message, err
	// 	},
	// }
	// Lookup[tls.AMT_TLSCredentialContext]["Create"] = Method{
	// 	Execute: func(value string) (client.Message, error) {
	// 		response, err := wsman.AMT.TLSCredentialContext.Create("Intel(r) AMT Certificate: Handle: 1")
	// 		return *response.Message, err
	// 	},
	// }
	// TLSProtocolEndpointCollection
	Lookup[tls.AMT_TLSProtocolEndpointCollection] = make(map[string]Method)
	Lookup[tls.AMT_TLSProtocolEndpointCollection]["Get"] = Method{
		Execute: func(instanceID string) (client.Message, error) {
			response, err := wsman.AMT.TLSProtocolEndpointCollection.Get()
			return *response.Message, err
		},
	}
	Lookup[tls.AMT_TLSProtocolEndpointCollection]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.TLSProtocolEndpointCollection.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[tls.AMT_TLSProtocolEndpointCollection]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.TLSProtocolEndpointCollection.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.TLSProtocolEndpointCollection.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// TLSSettingData
	Lookup[tls.AMT_TLSSettingData] = make(map[string]Method)
	Lookup[tls.AMT_TLSSettingData]["Get"] = Method{
		Execute: func(instanceID string) (client.Message, error) {
			response, err := wsman.AMT.TLSSettingData.Get("Intel(r) AMT 802.3 TLS Settings")
			return *response.Message, err
		},
	}
	Lookup[tls.AMT_TLSSettingData]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.TLSSettingData.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[tls.AMT_TLSSettingData]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.TLSSettingData.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.TLSSettingData.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// Lookup[tls.AMT_TLSSettingData]["Put"] = Method{
	// 	Execute: func(instanceID string) (client.Message, error) {
	// 		tlsSettingData := tls.TLSSettingDataRequest{
	// 			ElementName:                "Intel(r) AMT 802.3 TLS Settings",
	// 			InstanceID:                 "Intel(r) AMT 802.3 TLS Settings",
	// 			Enabled:                    false,
	// 			AcceptNonSecureConnections: true,
	// 		}
	// 		response, err := wsman.AMT.TLSSettingData.Put("Intel(r) AMT 802.3 TLS Settings", tlsSettingData)
	// 		return *response.Message, err
	// 	},
	// }
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
	Lookup[userinitiatedconnection.AMT_UserInitiatedConnectionService]["RequestStateChange"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.UserInitiatedConnectionService.RequestStateChange(userinitiatedconnection.BIOSandOSInterfacesEnabled)
			return *response.Message, err
		},
	}
	// WiFiPortConfigurationService
	Lookup[wifiportconfiguration.AMT_WiFiPortConfigurationService] = make(map[string]Method)
	Lookup[wifiportconfiguration.AMT_WiFiPortConfigurationService]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.WiFiPortConfigurationService.Get()
			return *response.Message, err
		},
	}
	Lookup[wifiportconfiguration.AMT_WiFiPortConfigurationService]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.AMT.WiFiPortConfigurationService.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[wifiportconfiguration.AMT_WiFiPortConfigurationService]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.AMT.WiFiPortConfigurationService.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.AMT.WiFiPortConfigurationService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	Lookup[wifiportconfiguration.AMT_WiFiPortConfigurationService]["Put"] = Method{
		Execute: func(value string) (client.Message, error) {
			wifiPortConfiguration := wifiportconfiguration.WiFiPortConfigurationServiceRequest{
				CreationClassName:                  "AMT_WiFiPortConfigurationService",
				ElementName:                        "Intel(r) AMT WiFiPort Configuration Service",
				EnabledState:                       5,
				HealthState:                        5,
				LastConnectedSsidUnderMeControl:    "",
				Name:                               "Intel(r) AMT WiFi Port Configuration Service",
				NoHostCsmeSoftwarePolicy:           0,
				RequestedState:                     12,
				SystemCreationClassName:            "CIM_ComputerSystem",
				SystemName:                         "Intel(r) AMT",
				LocalProfileSynchronizationEnabled: 3,
			}
			response, err := wsman.AMT.WiFiPortConfigurationService.Put(wifiPortConfiguration)
			return *response.Message, err
		},
	}
	// Lookup[wifiportconfiguration.AMT_WiFiPortConfigurationService]["AddWiFiSettings"] = Method{
	// 	Execute: func(value string) (client.Message, error) {
	// 		wifiEndpointSettings := wifi.WiFiEndpointSettings_INPUT{
	// 			AuthenticationMethod: wifi.AuthenticationMethod_WPA2_IEEE8021x,
	// 			BSSType: wifi.BSSType_Infrastructure,
	// 			ElementName: "TestWiFiProfile",
	// 			EncryptionMethod: wifi.EncryptionMethod_CCMP,
	// 			Priority: 1,
	// 			SSID: "TestWiFiProfile",
	// 		}
	// 		ieee8021xSettings := models.IEEE8021xSettings{
	// 			ElementName: "Test8021x",
	// 			AuthenticationProtocol: models.AuthenticationProtocolEAPFAST_MSCHAPv2,
	// 			Username: "testusername",
	// 			Password: "testpassword",
	// 			Domain: "vprodemo.com",
	// 		}
	// 		response, err := wsman.AMT.WiFiPortConfigurationService.AddWiFiSettings(wifiEndpointSettings, &ieee8021xSettings, "TestWiFiProfile", "testusername", "testpassword")
	// 		return *response.Message, err
	// 	},
	// }
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
			response, err := wsman.CIM.SoftwareIdentity.Get("AMTApps")
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
	// IPS AlarmClock
	Lookup[ipsalarmclock.IPS_AlarmClockOccurrence] = make(map[string]Method)
	Lookup[ipsalarmclock.IPS_AlarmClockOccurrence]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.IPS.AlarmClockOccurrence.Get(value)
			return *response.Message, err
		},
	}
	Lookup[ipsalarmclock.IPS_AlarmClockOccurrence]["Delete"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.IPS.AlarmClockOccurrence.Delete("test")
			return *response.Message, err
		},
	}
	Lookup[ipsalarmclock.IPS_AlarmClockOccurrence]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.IPS.AlarmClockOccurrence.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[ipsalarmclock.IPS_AlarmClockOccurrence]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.IPS.AlarmClockOccurrence.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.IPS.AlarmClockOccurrence.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// HostbasedSetupService
	Lookup[hostbasedsetup.IPS_HostBasedSetupService] = make(map[string]Method)
	Lookup[hostbasedsetup.IPS_HostBasedSetupService]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.IPS.HostBasedSetupService.Get()
			return *response.Message, err
		},
	}
	Lookup[hostbasedsetup.IPS_HostBasedSetupService]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.IPS.HostBasedSetupService.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[hostbasedsetup.IPS_HostBasedSetupService]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.IPS.HostBasedSetupService.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.IPS.HostBasedSetupService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// Lookup[hostbasedsetup.IPS_HostBasedSetupService]["AddNextCertInChain"] = Method{
	// 	Execute: func(value string) (client.Message, error) {
	// 		response, err := wsman.IPS.HostBasedSetupService.AddNextCertInChain()
	// 		return *response.Message, err
	// 	},
	// }
	// Lookup[hostbasedsetup.IPS_HostBasedSetupService]["AdminSetup"] = Method{
	// 	Execute: func(value string) (client.Message, error) {
	// 		response, err := wsman.IPS.HostBasedSetupService.AdminSetup()
	// 		return *response.Message, err
	// 	},
	// }
	// Lookup[hostbasedsetup.IPS_HostBasedSetupService]["Setup"] = Method{
	// 	Execute: func(value string) (client.Message, error) {
	// 		response, err := wsman.IPS.HostBasedSetupService.Setup()
	// 		return *response.Message, err
	// 	},
	// }
	// Lookup[hostbasedsetup.IPS_HostBasedSetupService]["UpgradeClientToAdmin"] = Method{
	// 	Execute: func(value string) (client.Message, error) {
	// 		response, err := wsman.IPS.HostBasedSetupService.UpgradeClientToAdmin()
	// 		return *response.Message, err
	// 	},
	// }
	// IPS IEEE8021x CredentialContext
	Lookup[ipsieee8021x.IPS_8021xCredentialContext] = make(map[string]Method)
	Lookup[ipsieee8021x.IPS_8021xCredentialContext]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.IPS.IEEE8021xCredentialContext.Get()
			return *response.Message, err
		},
	}
	Lookup[ipsieee8021x.IPS_8021xCredentialContext]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.IPS.IEEE8021xCredentialContext.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[ipsieee8021x.IPS_8021xCredentialContext]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.IPS.IEEE8021xCredentialContext.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.IPS.IEEE8021xCredentialContext.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// IPS IEEE8021x Settings
	Lookup[ipsieee8021x.IPS_IEEE8021xSettings] = make(map[string]Method)
	Lookup[ipsieee8021x.IPS_IEEE8021xSettings]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.IPS.IEEE8021xSettings.Get()
			return *response.Message, err
		},
	}
	Lookup[ipsieee8021x.IPS_IEEE8021xSettings]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.IPS.IEEE8021xSettings.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[ipsieee8021x.IPS_IEEE8021xSettings]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.IPS.IEEE8021xSettings.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.IPS.IEEE8021xSettings.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// Lookup[ipsieee8021x.IPS_IEEE8021xSettings]["Put"] = Method{
	// 	Execute: func(value string) (client.Message, error) {
	// 		response, err := wsman.IPS.IEEE8021xSettings.Put()
	// 		return *response.Message, err
	// 	},
	// }
	// Lookup[ipsieee8021x.IPS_IEEE8021xSettings]["SetCertificates"] = Method{
	// 	Execute: func(value string) (client.Message, error) {
	// 		response, err := wsman.IPS.IEEE8021xSettings.SetCertificates()
	// 		return *response.Message, err
	// 	},
	// }
	// OptInService
	Lookup[optin.IPS_OptInService] = make(map[string]Method)
	Lookup[optin.IPS_OptInService]["Get"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.IPS.OptInService.Get()
			return *response.Message, err
		},
	}
	Lookup[optin.IPS_OptInService]["Enumerate"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.IPS.OptInService.Enumerate()
			return *response.Message, err
		},
	}
	Lookup[optin.IPS_OptInService]["Pull"] = Method{
		Execute: func(value string) (client.Message, error) {
			er, err := wsman.IPS.OptInService.Enumerate()
			if err != nil {
				return client.Message{}, err
			}
			response, err := wsman.IPS.OptInService.Pull(er.Body.EnumerateResponse.EnumerationContext)
			return *response.Message, err
		},
	}
	// Lookup[optin.IPS_OptInService]["SendOptInCode"] = Method{
	// 	Execute: func(value string) (client.Message, error) {
	// 		response, err := wsman.IPS.OptInService.SendOptInCode()
	// 		return *response.Message, err
	// 	},
	// }
	Lookup[optin.IPS_OptInService]["StartOptIn"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.IPS.OptInService.StartOptIn()
			return *response.Message, err
		},
	}
	Lookup[optin.IPS_OptInService]["CancelOptIn"] = Method{
		Execute: func(value string) (client.Message, error) {
			response, err := wsman.IPS.OptInService.CancelOptIn()
			return *response.Message, err
		},
	}
}
