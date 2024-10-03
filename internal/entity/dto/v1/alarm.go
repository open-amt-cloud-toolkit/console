package dto

import "time"

type (
	AlarmClockOccurrence struct {
		ElementName        string    `json:"ElementName" binding:"required" example:"test"`
		InstanceID         string    `json:"InstanceID" binding:"" example:"test"`
		StartTime          time.Time `json:"StartTime" binding:"required" example:"2024-01-01T00:00:00Z"`
		Interval           int       `json:"Interval" default:"0" example:"1"`
		IntervalInMinutes  int       `json:"IntervalInMinutes" example:"1"`
		DeleteOnCompletion bool      `json:"DeleteOnCompletion" binding:"" example:"true"`
	}
	DeleteAlarmOccurrenceRequest struct {
		Name string `json:"Name" binding:"required" example:"test"`
	}

	AddAlarmOutput struct {
		ReturnValue int // Return code. 0 indicates success
	}
)
