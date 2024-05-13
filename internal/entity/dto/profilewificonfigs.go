package dto

type ProfileWiFiConfigs struct {
	Priority            int    `json:"priority,omitempty" example:"1"`
	WirelessProfileName string `json:"profileName" example:"My Profile"`
	ProfileName         string `json:"profileProfileName" example:"My Wireless Profile"`
	TenantID            string `json:"tenantId" example:"abc123"`
}
