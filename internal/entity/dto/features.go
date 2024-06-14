package dto

type Features struct {
	UserConsent string `json:"userConsent" example:"kvm"`
	EnableSOL   bool   `json:"enableSOL" example:"true"`
	EnableIDER  bool   `json:"enableIDER" example:"true"`
	EnableKVM   bool   `json:"enableKVM" example:"true"`
}
