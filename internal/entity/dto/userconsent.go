package dto

type UserConsent struct {
	ConsentCode string `json:"consentCode" binding:"required" example:"123456"`
}
