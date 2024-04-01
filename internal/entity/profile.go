package entity

type Profile struct {
	ProfileName                string               `json:"profileName,omitempty" example:"My Profile"`
	AMTPassword                string               `json:"amtPassword,omitempty" example:"my_password"`
	GenerateRandomPassword     bool                 `json:"generateRandomPassword,omitempty" example:"true"`
	CIRAConfigName             string               `json:"ciraConfigName,omitempty" example:"My CIRA Config"`
	Activation                 string               `json:"activation" example:"activate"`
	MEBXPassword               string               `json:"mebxPassword,omitempty" example:"my_password"`
	GenerateRandomMEBxPassword bool                 `json:"generateRandomMEBxPassword,omitempty" example:"true"`
	CIRAConfigObject           *CIRAConfig          `json:"ciraConfigObject,omitempty"`
	Tags                       []string             `json:"tags,omitempty" example:"tag1,tag2"`
	DhcpEnabled                bool                 `json:"dhcpEnabled,omitempty" example:"true"`
	IPSyncEnabled              bool                 `json:"ipSyncEnabled,omitempty" example:"true"`
	LocalWifiSyncEnabled       bool                 `json:"localWifiSyncEnabled,omitempty" example:"true"`
	WifiConfigs                []ProfileWifiConfigs `json:"wifiConfigs,omitempty"`
	TenantID                   string               `json:"tenantId" example:"abc123"`
	TLSMode                    int                  `json:"tlsMode,omitempty" example:"1"`
	TLSCerts                   *TLSCerts            `json:"tlsCerts,omitempty"`
	TLSSigningAuthority        string               `json:"tlsSigningAuthority,omitempty" example:"SelfSigned"`
	UserConsent                string               `json:"userConsent,omitempty" example:"All"`
	IDEREnabled                bool                 `json:"iderEnabled,omitempty" example:"true"`
	KVMEnabled                 bool                 `json:"kvmEnabled,omitempty" example:"true"`
	SOLEnabled                 bool                 `json:"solEnabled,omitempty" example:"true"`
	Ieee8021xProfileName       string               `json:"ieee8021xProfileName,omitempty" example:"My Profile"`
	Ieee8021xProfileObject     *IEEE8021xConfig     `json:"ieee8021xProfileObject,omitempty"`
	Version                    string               `json:"version,omitempty" example:"1.0.0"`
}

const (
	TLSModeNone int = iota
	TLSModeServerOnly
	TLSModeServerAllowNonTLS
	TLSModeMutualOnly
	TLSModeMutualAllowNonTLS
)

const (
	TLSSigningAuthoritySelfSigned  string = "SelfSigned"
	TLSSigningAuthorityMicrosoftCA string = "MicrosoftCA"
)

const (
	UserConsentNone    string = "None"
	UserConsentAll     string = "All"
	UserConsentKVMOnly string = "KVM"
)
