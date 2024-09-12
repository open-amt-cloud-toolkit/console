package dtov1

type UserConsent struct {
	ConsentCode string `json:"consentCode" binding:"required" example:"123456"`
}
