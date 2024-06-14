package dto

import "time"

type AlarmClockOccurrence struct {
	ElementName        string    `json:"elementName" binding:"required" example:"test"`
	InstanceID         string    `json:"instanceID" binding:"" example:"test"`
	StartTime          time.Time `json:"startTime" binding:"required" example:"2024-01-01T00:00:00Z"`
	Interval           string    `json:"interval" binding:"required" example:"1"`
	DeleteOnCompletion bool      `json:"deleteOnCompletion" binding:"required" example:"true"`
}
type DeleteAlarmOccurrenceRequest struct {
	InstanceID *string `json:"instanceID" binding:"" example:"test"`
}
