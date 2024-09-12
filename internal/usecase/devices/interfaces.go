package devices

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/power"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	dtov1 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	dtov2 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v2"
	wsmanAPI "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
)

type (
	WSMAN interface {
		SetupWsmanClient(device dtov1.Device, isRedirection, logMessages bool) wsmanAPI.Management
		DestroyWsmanClient(device dtov1.Device)
		Worker()
	}

	Redirection interface {
		SetupWsmanClient(device dtov1.Device, isRedirection, logMessages bool) wsman.Messages
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
		Get(ctx context.Context, top, skip int, tenantID string) ([]dtov1.Device, error)
		GetByID(ctx context.Context, guid, tenantID string) (*dtov1.Device, error)
		GetDistinctTags(ctx context.Context, tenantID string) ([]string, error)
		GetByTags(ctx context.Context, tags, method string, limit, offset int, tenantID string) ([]dtov1.Device, error)
		Delete(ctx context.Context, guid, tenantID string) error
		Update(ctx context.Context, d *dtov1.Device) (*dtov1.Device, error)
		Insert(ctx context.Context, d *dtov1.Device) (*dtov1.Device, error)
		GetByColumn(ctx context.Context, columnName, queryValue, tenantID string) ([]dtov1.Device, error)
		// Management Calls
		GetVersion(ctx context.Context, guid string) (dtov1.Version, dtov2.Version, error)
		GetFeatures(ctx context.Context, guid string) (dtov1.Features, error)
		SetFeatures(ctx context.Context, guid string, features dtov1.Features) (dtov1.Features, error)
		GetAlarmOccurrences(ctx context.Context, guid string) ([]dtov1.AlarmClockOccurrence, error)
		CreateAlarmOccurrences(ctx context.Context, guid string, alarm dtov1.AlarmClockOccurrence) (dtov1.AddAlarmOutput, error)
		DeleteAlarmOccurrences(ctx context.Context, guid, instanceID string) error
		GetHardwareInfo(ctx context.Context, guid string) (interface{}, error)
		GetPowerState(ctx context.Context, guid string) (map[string]interface{}, error)
		GetPowerCapabilities(ctx context.Context, guid string) (map[string]interface{}, error)
		GetGeneralSettings(ctx context.Context, guid string) (interface{}, error)
		CancelUserConsent(ctx context.Context, guid string) (interface{}, error)
		GetUserConsentCode(ctx context.Context, guid string) (map[string]interface{}, error)
		SendConsentCode(ctx context.Context, code dtov1.UserConsent, guid string) (interface{}, error)
		SendPowerAction(ctx context.Context, guid string, action int) (power.PowerActionResponse, error)
		SetBootOptions(ctx context.Context, guid string, bootSetting dtov1.BootSetting) (power.PowerActionResponse, error)
		GetAuditLog(ctx context.Context, startIndex int, guid string) (dtov1.AuditLog, error)
		GetEventLog(ctx context.Context, guid string) ([]dtov1.EventLog, error)
		Redirect(ctx context.Context, conn *websocket.Conn, guid, mode string) error
		GetNetworkSettings(c context.Context, guid string) (dtov1.NetworkSettings, error)
		GetCertificates(c context.Context, guid string) (dtov1.SecuritySettings, error)
		GetTLSSettingData(c context.Context, guid string) ([]dtov1.SettingDataResponse, error)
		GetDiskInfo(c context.Context, guid string) (interface{}, error)
		GetDeviceCertificate(c context.Context, guid string) (dtov1.Certificate, error)
	}
)
