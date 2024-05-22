package devices

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman"
	amtAlarmClock "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/alarmclock"
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
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/alarmclock"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/optin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
)

type (
	Management interface {
		SetupWsmanClient(device dto.Device, isRedirection, logAMTMessages bool)
		GetAMTVersion() ([]software.SoftwareIdentity, error)
		GetSetupAndConfiguration() ([]setupandconfiguration.SetupAndConfigurationServiceResponse, error)
		GetFeatures() (interface{}, error)
		SetFeatures(dto.Features) (dto.Features, error)
		GetAlarmOccurrences() ([]alarmclock.AlarmClockOccurrence, error)
		CreateAlarmOccurrences(name string, startTime time.Time, interval int, deleteOnCompletion bool) (amtAlarmClock.AddAlarmOutput, error)
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
		GetCredentialRelationships() (credential.Items, error)
		GetConcreteDependencies() ([]concrete.ConcreteDependency, error)
	}
	Redirection interface {
		SetupWsmanClient(device dto.Device, isRedirection, logAMTMessages bool) wsman.Messages
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
		// Management Calls
		GetVersion(ctx context.Context, guid string) (map[string]interface{}, error)
		GetFeatures(ctx context.Context, guid string) (interface{}, error)
		SetFeatures(ctx context.Context, guid string, features dto.Features) (dto.Features, error)
		GetAlarmOccurrences(ctx context.Context, guid string) ([]alarmclock.AlarmClockOccurrence, error)
		CreateAlarmOccurrences(ctx context.Context, guid string, alarm dto.AlarmClockOccurrence) (amtAlarmClock.AddAlarmOutput, error)
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
	}
)
