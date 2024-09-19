package amtexplorer

import (
	"context"

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
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
)

type (
	AMTExplorer interface {
		// SetupWsmanClient(device dto.Device, isRedirection, logAMTMessages bool) *wsmanAPI.ConnectionEntry
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
		GetAMTWiFiPortConfigurationService() (wifiportconfiguration.Response, error)
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
	Feature interface {
		GetExplorerSupportedCalls() []string
		ExecuteCall(ctx context.Context, guid, call, tenantID string) (*dto.Explorer, error)
	}
	Repository interface {
		GetByID(ctx context.Context, guid, tenantID string) (*entity.Device, error)
	}
	WSMAN interface {
		SetupWsmanClient(device dto.Device, logMessages bool) AMTExplorer
		DestroyWsmanClient(device dto.Device)
	}
)
