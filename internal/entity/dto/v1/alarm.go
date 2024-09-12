package dtov1

import "time"

type (
	AlarmClockOccurrence struct {
		ElementName        string    `json:"elementName" binding:"required" example:"test"`
		InstanceID         string    `json:"instanceID" binding:"" example:"test"`
		StartTime          time.Time `json:"startTime" binding:"required" example:"2024-01-01T00:00:00Z"`
		Interval           int       `json:"interval" binding:"number" example:"1"`
		DeleteOnCompletion bool      `json:"deleteOnCompletion" binding:"required" example:"true"`
	}
	DeleteAlarmOccurrenceRequest struct {
		InstanceID *string `json:"instanceID" binding:"" example:"test"`
	}

	AddAlarmOutput struct {
		ReturnValue int // Return code. 0 indicates success
	}
)
