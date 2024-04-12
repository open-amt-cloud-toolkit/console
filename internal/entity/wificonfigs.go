package entity

type ProfileWifiConfigs struct {
	Priority    int    `json:"priority,omitempty" example:"1"`
	ProfileName string `json:"profileName" example:"My Profile"`
	TenantID    string `json:"tenantId" example:"abc123"`
}

type WirelessConfig struct {
	ProfileName            string           `json:"profileName,omitempty" example:"My Profile"`
	AuthenticationMethod   int              `json:"authenticationMethod" example:"1"`
	EncryptionMethod       int              `json:"encryptionMethod" example:"2"`
	SSID                   string           `json:"ssid" example:"abc"`
	PSKValue               int              `json:"pskValue" example:"3"`
	PSKPassphrase          string           `json:"pskPassphrase" example:"abc"`
	LinkPolicy             []string         `json:"linkPolicy"`
	TenantID               string           `json:"tenantId" example:"abc123"`
	IEEE8021xProfileName   string           `json:"ieee8021xProfileName,omitempty" example:"My Profile"`
	IEEE8021xProfileObject *IEEE8021xConfig `json:"ieee8021xProfileObject,omitempty"`
	Version                string           `json:"version"`
}
