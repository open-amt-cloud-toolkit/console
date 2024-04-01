package dto

type UserConsent struct {
	ConsentCode int `json:"consentCode" binding:"required" example:"123456"`
}
