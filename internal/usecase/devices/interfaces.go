package devices

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/alarmclock"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/auditlog"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/authorization"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/boot"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/environmentdetection"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/ethernetport"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/general"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/ieee8021x"
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
	cimBoot "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/boot"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/card"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/chassis"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/chip"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/computer"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/concrete"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/credential"
	cimIEEE8021x "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/ieee8021x"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/kvm"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/mediaaccess"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/physical"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/power"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/processor"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/service"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/software"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/system"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/wifi"
	ipsAlarmClock "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/alarmclock"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/hostbasedsetup"
	ipsIEEE8021x "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/ieee8021x"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/optin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	wsmanAPI "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
)

type (
	AMTExplorer interface {
		SetupWsmanClient(device dto.Device, isRedirection, logAMTMessages bool)
		GetAMT8021xCredentialContext() (ieee8021x.Response, error)
		GetAMT8021xProfile() (ieee8021x.Response, error)
		GetAMTAlarmClockService() (alarmclock.Response, error)
		GetAMTAuditLog() (auditlog.Response, error)
		GetAMTAuthorizationService() (authorization.Response, error)
		GetAMTBootCapabilities() (boot.Response, error)
		GetAMTBootSettingData() (boot.Response, error)
		GetAMTEnvironmentDetectionSettingData() (environmentdetection.Response, error)
		GetAMTEthernetPortSettings() (ethernetport.Response, error)
		GetAMTGeneralSettings() (general.Response, error)
		GetAMTKerberosSettingData() (kerberos.Response, error)
		GetAMTManagementPresenceRemoteSAP() (managementpresence.Response, error)
		GetAMTMessageLog() (messagelog.Response, error)
		GetAMTMPSUsernamePassword() (mps.Response, error)
		GetAMTPublicKeyCertificate() (publickey.Response, error)
		GetAMTPublicKeyManagementService() (publickey.Response, error)
		GetAMTPublicPrivateKeyPair() (publicprivate.Response, error)
		GetAMTRedirectionService() (redirection.Response, error)
		GetAMTRemoteAccessPolicyAppliesToMPS() (remoteaccess.Response, error)
		GetAMTRemoteAccessPolicyRule() (remoteaccess.Response, error)
		GetAMTRemoteAccessService() (remoteaccess.Response, error)
		GetAMTSetupAndConfigurationService() (setupandconfiguration.Response, error)
		GetAMTTimeSynchronizationService() (timesynchronization.Response, error)
		GetAMTTLSCredentialContext() (tls.Response, error)
		GetAMTTLSProtocolEndpointCollection() (tls.Response, error)
		GetAMTTLSSettingData() (tls.Response, error)
		GetAMTUserInitiatedConnectionService() (userinitiatedconnection.Response, error)
		GetAMTWiFiPortConifgurationService() (wifiportconfiguration.Response, error)
		GetCIMBIOSElement() (bios.Response, error)
		GetCIMBootConfigSetting() (cimBoot.Response, error)
		GetCIMBootService() (cimBoot.Response, error)
		GetCIMBootSourceSetting() (cimBoot.Response, error)
		GetCIMCard() (card.Response, error)
		GetCIMChassis() (chassis.Response, error)
		GetCIMChip() (chip.Response, error)
		GetCIMComputerSystemPackage() (computer.Response, error)
		GetCIMConcreteDependency() (concrete.Response, error)
		GetCIMCredentialContext() (credential.Response, error)
		GetCIMIEEE8021xSettings() (cimIEEE8021x.Response, error)
		GetCIMKVMRedirectionSAP() (kvm.Response, error)
		GetCIMMediaAccessDevice() (mediaaccess.Response, error)
		GetCIMPhysicalMemory() (physical.Response, error)
		GetCIMPhysicalPackage() (physical.Response, error)
		GetCIMPowerManagementService() (power.Response, error)
		GetCIMProcessor() (processor.Response, error)
		GetCIMServiceAvailableToElement() (service.Response, error)
		GetCIMSoftwareIdentity() (software.Response, error)
		GetCIMSystemPackaging() (system.Response, error)
		GetCIMWiFiEndpointSettings() (wifi.Response, error)
		GetCIMWiFiPort() (wifi.Response, error)
		GetIPS8021xCredentialContext() (ipsIEEE8021x.Response, error)
		GetIPSAlarmClockOccurrence() (ipsAlarmClock.Response, error)
		GetIPSHostBasedSetupService() (hostbasedsetup.Response, error)
		GetIPSIEEE8021xSettings() (ipsIEEE8021x.Response, error)
		GetIPSOptInService() (optin.Response, error)
	}
	Management interface {
		SetupWsmanClient(device dto.Device, isRedirection, logMessages bool)
		DestroyWsmanClient(device dto.Device)
		GetAMTVersion() ([]software.SoftwareIdentity, error)
		GetSetupAndConfiguration() ([]setupandconfiguration.SetupAndConfigurationServiceResponse, error)
		GetFeatures() (dto.Features, error)
		SetFeatures(dto.Features) (dto.Features, error)
		GetAlarmOccurrences() ([]ipsAlarmClock.AlarmClockOccurrence, error)
		CreateAlarmOccurrences(name string, startTime time.Time, interval int, deleteOnCompletion bool) (alarmclock.AddAlarmOutput, error)
		DeleteAlarmOccurrences(instanceID string) error
		GetHardwareInfo() (interface{}, error)
		GetPowerState() ([]service.CIM_AssociatedPowerManagementService, error)
		GetPowerCapabilities() (boot.BootCapabilitiesResponse, error)
		GetGeneralSettings() (interface{}, error)
		CancelUserConsentRequest() (interface{}, error)
		GetUserConsentCode() (optin.StartOptIn_OUTPUT, error)
		SendConsentCode(code int) (interface{}, error)
		SendPowerAction(action int) (power.PowerActionResponse, error)
		GetBootData() (boot.BootCapabilitiesResponse, error)
		SetBootData(data boot.BootSettingDataRequest) (interface{}, error)
		SetBootConfigRole(role int) (interface{}, error)
		ChangeBootOrder(bootSource string) (cimBoot.ChangeBootOrder_OUTPUT, error)
		GetAuditLog(startIndex int) (auditlog.Response, error)
		GetEventLog() (messagelog.GetRecordsResponse, error)
		GetNetworkSettings() (interface{}, error)
		GetCertificates() (wsmanAPI.Certificates, error)
		GetCredentialRelationships() (credential.Items, error)
		GetConcreteDependencies() ([]concrete.ConcreteDependency, error)
	}
	Redirection interface {
		SetupWsmanClient(device dto.Device, isRedirection, logMessages bool) wsman.Messages
		RedirectConnect(ctx context.Context, deviceConnection *DeviceConnection) error
		RedirectClose(ctx context.Context, deviceConnection *DeviceConnection) error
		RedirectListen(ctx context.Context, deviceConnection *DeviceConnection) ([]byte, error)
		RedirectSend(ctx context.Context, deviceConnection *DeviceConnection, message []byte) error
	}
	Repository interface {
		GetCount(context.Context, string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Device, error)
		GetByID(ctx context.Context, guid, tenantID string) (*entity.Device, error)
		GetDistinctTags(ctx context.Context, tenantID string) ([]string, error)
		GetByTags(ctx context.Context, tags []string, method string, limit, offset int, tenantID string) ([]entity.Device, error)
		Delete(ctx context.Context, guid, tenantID string) (bool, error)
		Update(ctx context.Context, d *entity.Device) (bool, error)
		Insert(ctx context.Context, d *entity.Device) (string, error)
		GetByColumn(ctx context.Context, columnName, queryValue, tenantID string) ([]entity.Device, error)
	}
	Feature interface {
		// Repository/Database Calls
		GetCount(context.Context, string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]dto.Device, error)
		GetByID(ctx context.Context, guid, tenantID string) (*dto.Device, error)
		GetDistinctTags(ctx context.Context, tenantID string) ([]string, error)
		GetByTags(ctx context.Context, tags, method string, limit, offset int, tenantID string) ([]dto.Device, error)
		Delete(ctx context.Context, guid, tenantID string) error
		Update(ctx context.Context, d *dto.Device) (*dto.Device, error)
		Insert(ctx context.Context, d *dto.Device) (*dto.Device, error)
		GetByColumn(ctx context.Context, columnName, queryValue, tenantID string) ([]dto.Device, error)
		// Management Calls
		GetVersion(ctx context.Context, guid string) (map[string]interface{}, error)
		GetFeatures(ctx context.Context, guid string) (dto.Features, error)
		SetFeatures(ctx context.Context, guid string, features dto.Features) (dto.Features, error)
		GetAlarmOccurrences(ctx context.Context, guid string) ([]dto.AlarmClockOccurrence, error)
		CreateAlarmOccurrences(ctx context.Context, guid string, alarm dto.AlarmClockOccurrence) (dto.AddAlarmOutput, error)
		DeleteAlarmOccurrences(ctx context.Context, guid, instanceID string) error
		GetHardwareInfo(ctx context.Context, guid string) (interface{}, error)
		GetPowerState(ctx context.Context, guid string) (map[string]interface{}, error)
		GetPowerCapabilities(ctx context.Context, guid string) (map[string]interface{}, error)
		GetGeneralSettings(ctx context.Context, guid string) (interface{}, error)
		CancelUserConsent(ctx context.Context, guid string) (interface{}, error)
		GetUserConsentCode(ctx context.Context, guid string) (map[string]interface{}, error)
		SendConsentCode(ctx context.Context, code dto.UserConsent, guid string) (interface{}, error)
		SendPowerAction(ctx context.Context, guid string, action int) (power.PowerActionResponse, error)
		SetBootOptions(ctx context.Context, guid string, bootSetting dto.BootSetting) (power.PowerActionResponse, error)
		GetAuditLog(ctx context.Context, startIndex int, guid string) (dto.AuditLog, error)
		GetEventLog(ctx context.Context, guid string) ([]dto.EventLog, error)
		Redirect(ctx context.Context, conn *websocket.Conn, guid, mode string) error
		GetNetworkSettings(c context.Context, guid string) (interface{}, error)
		GetCertificates(c context.Context, guid string) (dto.SecuritySettings, error)
		GetExplorerSupportedCalls() []string
		ExecuteCall(ctx context.Context, guid, call, tenantID string) (*dto.Explorer, error)
	}
)
