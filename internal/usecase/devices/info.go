package devices

import (
	"context"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/messagelog"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
)

func (uc *UseCase) GetVersion(c context.Context, guid string) (map[string]interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	version, err := uc.device.GetAMTVersion()
	if err != nil {
		return nil, err
	}

	data, err := uc.device.GetSetupAndConfiguration()
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"CIM_SoftwareIdentity": map[string]interface{}{
			"responses": version,
		},
		"AMT_SetupAndConfigurationService": map[string]interface{}{
			"response": data[0],
		},
	}

	return response, nil
}

func (uc *UseCase) GetFeatures(c context.Context, guid string) (interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	features, err := uc.device.GetFeatures()
	if err != nil {
		return nil, err
	}

	return features, nil
}

func (uc *UseCase) SetFeatures(c context.Context, guid string, features dto.Features) (dto.Features, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return features, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	features, err = uc.device.SetFeatures(features)
	if err != nil {
		return features, err
	}

	return features, nil
}

func (uc *UseCase) GetHardwareInfo(c context.Context, guid string) (interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	hwInfo, err := uc.device.GetHardwareInfo()
	if err != nil {
		return nil, err
	}

	return hwInfo, nil
}

func (uc *UseCase) GetAuditLog(c context.Context, startIndex int, guid string) (dto.AuditLog, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return dto.AuditLog{}, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	response, err := uc.device.GetAuditLog(startIndex)
	if err != nil {
		return dto.AuditLog{}, err
	}

	auditLogResponse := dto.AuditLog{}
	auditLogResponse.TotalCount = response.Body.ReadRecordsResponse.TotalRecordCount
	auditLogResponse.Records = response.Body.DecodedRecordsResponse

	return auditLogResponse, nil
}

func (uc *UseCase) GetEventLog(c context.Context, guid string) (messagelog.GetRecordsResponse, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return messagelog.GetRecordsResponse{}, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	eventLogs, err := uc.device.GetEventLog()
	if err != nil {
		return messagelog.GetRecordsResponse{}, err
	}

	return eventLogs, nil
}

func (uc *UseCase) GetGeneralSettings(c context.Context, guid string) (interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	generalSettings, err := uc.device.GetGeneralSettings()
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"Body": generalSettings,
	}

	return response, nil
}
