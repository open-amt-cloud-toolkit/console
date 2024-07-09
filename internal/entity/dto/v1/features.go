package dto_v1

type Features struct {
	Redirection bool   `json:"redirection"`
	KVM         bool   `json:"KVM"`
	SOL         bool   `json:"SOL"`
	IDER        bool   `json:"IDER"`
	OptInState  int    `json:"optInState"`
	UserConsent string `json:"userConsent"`
}
