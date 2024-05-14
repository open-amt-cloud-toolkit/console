package devices

import (
	"context"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
)

func (uc *UseCase) CancelUserConsent(c context.Context, guid string) (interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	response, err := uc.device.CancelUserConsentRequest()
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (uc *UseCase) GetUserConsentCode(c context.Context, guid string) (map[string]interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	code, err := uc.device.GetUserConsentCode()
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"Body": code,
	}

	return response, nil
}

func (uc *UseCase) SendConsentCode(c context.Context, userConsent dto.UserConsent, guid string) (interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	response, err := uc.device.SendConsentCode(userConsent.ConsentCode)
	if err != nil {
		return nil, err
	}

	return response, nil
}
