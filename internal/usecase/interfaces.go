// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"time"

	amtAlarmClock "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/alarmclock"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/boot"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/messagelog"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/setupandconfiguration"
	cimBoot "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/boot"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/power"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/software"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/alarmclock"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/optin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Translation -.
	Domain interface {
		GetCount(context.Context, string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Domain, error)
		GetDomainByDomainSuffix(ctx context.Context, domainSuffix, tenantID string) (*entity.Domain, error)
		GetByName(ctx context.Context, name, tenantID string) (*entity.Domain, error)
		Delete(ctx context.Context, name, tenantID string) (bool, error)
		Update(ctx context.Context, d *entity.Domain) (bool, error)
		Insert(ctx context.Context, d *entity.Domain) (string, error)
	}
	Device interface {
		GetCount(context.Context, string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Device, error)
		GetByID(ctx context.Context, guid, tenantID string) (entity.Device, error)
		GetDistinctTags(ctx context.Context, tenantID string) ([]string, error)
		GetByTags(ctx context.Context, tags []string, method string, limit, offset int, tenantID string) ([]entity.Device, error)
		Delete(ctx context.Context, guid, tenantID string) (bool, error)
		Update(ctx context.Context, d *entity.Device) (bool, error)
		Insert(ctx context.Context, d *entity.Device) (string, error)
	}
	Profile interface {
		GetCount(ctx context.Context, tenantID string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Profile, error)
		GetByName(ctx context.Context, profileName, tenantID string) (entity.Profile, error)
		Delete(ctx context.Context, profileName, tenantID string) (bool, error)
		Update(ctx context.Context, p *entity.Profile) (bool, error)
		Insert(ctx context.Context, p *entity.Profile) (string, error)
	}
	DeviceManagement interface {
		SetupWsmanClient(device entity.Device, logAMTMessages bool)
		GetAMTVersion() ([]software.SoftwareIdentity, error)
		GetSetupAndConfiguration() ([]setupandconfiguration.SetupAndConfigurationServiceResponse, error)
		GetFeatures() (interface{}, error)
		SetFeatures(dto.Features) (dto.Features, error)
		GetAlarmOccurrences() ([]alarmclock.AlarmClockOccurrence, error)
		CreateAlarmOccurrences(name string, startTime time.Time, interval int, deleteOnCompletion bool) (amtAlarmClock.AddAlarmOutput, error)
		DeleteAlarmOccurrences(instanceID string) error
		GetHardwareInfo() (interface{}, error)
		GetPowerState() (interface{}, error)
		GetPowerCapabilities() (boot.BootCapabilitiesResponse, error)
		GetGeneralSettings() (interface{}, error)
		CancelUserConsent() (interface{}, error)
		GetUserConsentCode() (optin.StartOptIn_OUTPUT, error)
		SendConsentCode(code int) (interface{}, error)
		SendPowerAction(action int) (power.PowerActionResponse, error)
		GetBootData() (boot.BootCapabilitiesResponse, error)
		SetBootData(data boot.BootSettingDataRequest) (interface{}, error)
		SetBootConfigRole(role int) (interface{}, error)
		ChangeBootOrder(bootSource string) (cimBoot.ChangeBootOrder_OUTPUT, error)
		GetAuditLog(startIndex int) (dto.AuditLog, error)
		GetEventLog() (messagelog.GetRecordsResponse, error)
	}
	IEEE8021xProfile interface {
		CheckProfileExists(ctx context.Context, profileName, tenantID string) (bool, error)
		GetCount(ctx context.Context, tenantID string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.IEEE8021xConfig, error)
		GetByName(ctx context.Context, profileName, tenantID string) (entity.IEEE8021xConfig, error)
		Delete(ctx context.Context, profileName, tenantID string) (bool, error)
		Update(ctx context.Context, p *entity.IEEE8021xConfig) (bool, error)
		Insert(ctx context.Context, p *entity.IEEE8021xConfig) (string, error)
	}
	CIRAConfig interface {
		GetCount(ctx context.Context, tenantID string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.CIRAConfig, error)
		GetByName(ctx context.Context, configName, tenantID string) (entity.CIRAConfig, error)
		Delete(ctx context.Context, profileName, tenantID string) (bool, error)
		Update(ctx context.Context, p *entity.CIRAConfig) (bool, error)
		Insert(ctx context.Context, p *entity.CIRAConfig) (string, error)
	}
	WirelessProfile interface {
		CheckProfileExists(ctx context.Context, profileName, tenantID string) (bool, error)
		GetCount(ctx context.Context, tenantID string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.WirelessConfig, error)
		GetByName(ctx context.Context, guid, tenantID string) (entity.WirelessConfig, error)
		Delete(ctx context.Context, profileName, tenantID string) (bool, error)
		Update(ctx context.Context, p *entity.WirelessConfig) (bool, error)
		Insert(ctx context.Context, p *entity.WirelessConfig) (string, error)
	}
)
