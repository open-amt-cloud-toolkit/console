package dto

import "encoding/xml"

type (
	GetUserConsentMessage struct {
		Body UserConsentMessage `json:"Body" binding:"required"`
	}

	UserConsentMessage struct {
		Name        xml.Name `json:"XMLName" binding:"required"`
		ReturnValue int      `json:"ReturnValue" binding:"required"`
	}

	UserConsentCode struct {
		ConsentCode string `json:"consentCode" binding:"required" example:"123456"`
	}
)
