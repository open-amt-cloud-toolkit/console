package dto

import "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/auditlog"

type AuditLog struct {
	TotalCount int                       `json:"totalCnt" binding:"required" example:"0"`
	Records    []auditlog.AuditLogRecord `json:"records" binding:"required"`
}

// type AuditLogRecord struct {
// 	AuditAppID     int       `json:"AuditAppId" binding:"required" example:"0"`
// 	EventID        int       `json:"EventId" binding:"required" example:"0"`
// 	InitiatorType  int       `json:"InitiatorType" binding:"required" example:"0"`
// 	AuditApp       string    `json:"AuditApp" binding:"required" example:"MikeIsAmazing"`
// 	Event          string    `json:"Event" binding:"required" example:"MikeIsAmazing"`
// 	Initiator      string    `json:"Initiator" binding:"required" example:"MikeIsAmazing"`
// 	Time           time.Time `json:"Time" binding:"required" example:"MikeIsAmazing"`
// 	MCLocationType int       `json:"MCLocationType" binding:"required" example:"0"`
// 	NetAddress     string    `json:"NetAddress" binding:"required" example:"MikeIsAmazing"`
// 	Ex             string    `json:"Ex" binding:"required" example:"MikeIsAmazing"`
// 	ExStr          string    `json:"ExStr" binding:"required" example:"MikeIsAmazing"`
// }
