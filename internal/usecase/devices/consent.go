package devices

import (
	"context"
	"strconv"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
)

func (uc *UseCase) CancelUserConsent(c context.Context, guid string) (interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	response, err := device.CancelUserConsentRequest()
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

	device := uc.device.SetupWsmanClient(*item, false, true)

	code, err := device.GetUserConsentCode()
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

	device := uc.device.SetupWsmanClient(*item, false, true)

	consentCode, _ := strconv.Atoi(userConsent.ConsentCode)

	response, err := device.SendConsentCode(consentCode)
	if err != nil {
		return nil, err
	}

	return response, nil
}
