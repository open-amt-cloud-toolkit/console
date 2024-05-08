package entity

type Profile struct {
	ProfileName                string
	AMTPassword                string
	CreationDate               string
	CreatedBy                  string
	GenerateRandomPassword     bool
	CIRAConfigName             *string
	Activation                 string
	MEBXPassword               string
	GenerateRandomMEBxPassword bool
	Tags                       string
	DHCPEnabled                bool
	IPSyncEnabled              bool
	LocalWiFiSyncEnabled       bool
	TenantID                   string
	TLSMode                    int
	TLSSigningAuthority        string
	UserConsent                string
	IDEREnabled                bool
	KVMEnabled                 bool
	SOLEnabled                 bool
	IEEE8021xProfileName       *string
	Version                    string
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
