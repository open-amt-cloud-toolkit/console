package devices

import (
	"context"
)

func (uc *UseCase) GetNetworkSettings(c context.Context, guid string) (interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	response, err := device.GetNetworkSettings()
	if err != nil {
		return nil, err
	}

	return response, nil
}
