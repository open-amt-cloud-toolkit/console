package dto

import "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/auditlog"

type AuditLog struct {
	TotalCount int                       `json:"totalCnt" binding:"required" example:"0"`
	Records    []auditlog.AuditLogRecord `json:"records" binding:"required"`
}
