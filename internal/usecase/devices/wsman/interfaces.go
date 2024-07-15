package wsman

import (
	"time"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/alarmclock"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/auditlog"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/boot"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/messagelog"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/setupandconfiguration"
	cimBoot "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/boot"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/concrete"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/credential"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/power"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/service"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/software"
	ipsAlarmClock "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/alarmclock"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/optin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
)

type Management interface {
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
	GetCertificates() (Certificates, error)
	GetCredentialRelationships() (credential.Items, error)
	GetConcreteDependencies() ([]concrete.ConcreteDependency, error)
}
