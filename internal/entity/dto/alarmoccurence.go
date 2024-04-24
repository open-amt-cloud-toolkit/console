package dto

import "time"

type AlarmClockOccurrence struct {
	InstanceID         string    `json:"instanceID" binding:"required" example:"test"`
	StartTime          time.Time `json:"startTime" binding:"required" example:"2024-01-01T00:00:00Z"`
	Interval           int       `json:"interval" binding:"required" example:"1"`
	DeleteOnCompletion bool      `json:"deleteOnCompletion" binding:"required" example:"true"`
}
