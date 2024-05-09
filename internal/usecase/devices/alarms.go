package devices

import (
	"context"
	"strconv"

	amtAlarmClock "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/alarmclock"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/alarmclock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/utils"
)

func (uc *UseCase) GetAlarmOccurrences(c context.Context, guid string) ([]alarmclock.AlarmClockOccurrence, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil || item.GUID == "" {
		return nil, utils.ErrNotFound
	}

	uc.device.SetupWsmanClient(*item, false, true)

	alarms, err := uc.device.GetAlarmOccurrences()
	if err != nil {
		return nil, err
	}

	if alarms == nil {
		alarms = []alarmclock.AlarmClockOccurrence{}
	}

	return alarms, nil
}

func (uc *UseCase) CreateAlarmOccurrences(c context.Context, guid string, alarm dto.AlarmClockOccurrence) (amtAlarmClock.AddAlarmOutput, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil || item.GUID == "" {
		return amtAlarmClock.AddAlarmOutput{}, nil
	}

	uc.device.SetupWsmanClient(*item, false, true)

	interval, err := strconv.Atoi(alarm.Interval)
	if err != nil {
		return amtAlarmClock.AddAlarmOutput{}, err
	}

	alarmReference, err := uc.device.CreateAlarmOccurrences(alarm.InstanceID, alarm.StartTime, interval, alarm.DeleteOnCompletion)
	if err != nil {
		return amtAlarmClock.AddAlarmOutput{}, err
	}

	return alarmReference, nil
}

func (uc *UseCase) DeleteAlarmOccurrences(c context.Context, guid, instanceID string) error {
	item, err := uc.GetByID(c, guid, "")
	if err != nil || item.GUID == "" {
		return err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	err = uc.device.DeleteAlarmOccurrences(instanceID)
	if err != nil {
		return err
	}

	return nil
}
