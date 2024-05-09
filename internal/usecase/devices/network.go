package devices

import (
	"context"

	"github.com/open-amt-cloud-toolkit/console/internal/usecase/utils"
)

func (uc *UseCase) GetNetworkSettings(c context.Context, guid string) (interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil || item.GUID == "" {
		return nil, utils.ErrNotFound
	}

	uc.device.SetupWsmanClient(*item, false, true)

	response, err := uc.device.GetNetworkSettings()
	if err != nil {
		return nil, err
	}

	return response, nil
}
