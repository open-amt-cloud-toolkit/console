package devices

import (
	"context"
	"strconv"
	"time"

	amtAlarmClock "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/alarmclock"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/alarmclock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
)

func (uc *UseCase) GetAlarmOccurrences(c context.Context, guid string) ([]dto.AlarmClockOccurrence, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	alarms, err := device.GetAlarmOccurrences()
	if err != nil {
		return nil, err
	}

	if alarms == nil {
		alarms = []alarmclock.AlarmClockOccurrence{}
	}

	// iterate over the data and convert each entity to dto
	d1 := make([]dto.AlarmClockOccurrence, len(alarms))

	for i := range alarms {
		tmpEntity := alarms[i] // create a new variable to avoid memory aliasing
		d1[i] = *uc.alarmOccurenceEntityToDTO(&tmpEntity)
	}

	return d1, nil
}

func (uc *UseCase) CreateAlarmOccurrences(c context.Context, guid string, alarm dto.AlarmClockOccurrence) (dto.AddAlarmOutput, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return dto.AddAlarmOutput{}, err
	}

	alarm.InstanceID = alarm.ElementName

	device := uc.device.SetupWsmanClient(*item, false, true)

	alarmReference, err := device.CreateAlarmOccurrences(alarm.InstanceID, alarm.StartTime, alarm.Interval, alarm.DeleteOnCompletion)
	if err != nil {
		return dto.AddAlarmOutput{}, ErrAMT.Wrap("CreateAlarmOccurrences", "device.CreateAlarmOccurrences", err)
	}

	d1 := *uc.addAlarmOutputEntityToDTO(&alarmReference)

	return d1, nil
}

func (uc *UseCase) DeleteAlarmOccurrences(c context.Context, guid, instanceID string) error {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	err = device.DeleteAlarmOccurrences(instanceID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) addAlarmOutputEntityToDTO(d *amtAlarmClock.AddAlarmOutput) *dto.AddAlarmOutput {
	d1 := &dto.AddAlarmOutput{
		ReturnValue: int(d.ReturnValue),
	}

	return d1
}

func (uc *UseCase) alarmOccurenceEntityToDTO(d *alarmclock.AlarmClockOccurrence) *dto.AlarmClockOccurrence {
	startTime, _ := time.Parse(time.RFC3339, d.StartTime)
	interval, _ := strconv.Atoi(d.Interval)
	d1 := &dto.AlarmClockOccurrence{
		ElementName:        d.ElementName,
		InstanceID:         d.InstanceID,
		StartTime:          startTime,
		Interval:           interval,
		DeleteOnCompletion: d.DeleteOnCompletion,
	}

	return d1
}
