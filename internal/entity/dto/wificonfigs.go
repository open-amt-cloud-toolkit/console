package dto

type WirelessConfigCountResponse struct {
	Count int              `json:"totalCount"`
	Data  []WirelessConfig `json:"data"`
}

type WirelessConfig struct {
	ProfileName            string           `json:"profileName,omitempty" example:"My Profile"`
	AuthenticationMethod   int              `json:"authenticationMethod" binding:"oneof=4 5 6 7" example:"1"`
	EncryptionMethod       int              `json:"encryptionMethod" binding:"oneof=3 4" example:"2"`
	SSID                   string           `json:"ssid" binding:"max=32" example:"abc"`
	PSKValue               int              `json:"pskValue" example:"3"`
	PSKPassphrase          string           `json:"pskPassphrase,omitempty" binding:"omitempty,min=8,max=32" example:"abc"`
	LinkPolicy             []int            `json:"linkPolicy"`
	TenantID               string           `json:"tenantId" example:"abc123"`
	IEEE8021xProfileName   *string          `json:"ieee8021xProfileName,omitempty" example:"My Profile"`
	IEEE8021xProfileObject *IEEE8021xConfig `json:"ieee8021xProfileObject,omitempty"`
	Version                string           `json:"version"`
}

type ProfileWiFiConfigs struct {
	Priority    int    `json:"priority,omitempty" example:"1"`
	ProfileName string `json:"profileName" example:"My Profile"`
	TenantID    string `json:"tenantId" example:"abc123"`
}
