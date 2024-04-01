package entity

type ProfileWifiConfigs struct {
	Priority    int    `json:"priority,omitempty" example:"1"`
	ProfileName string `json:"profileName" example:"My Profile"`
	TenantID    string `json:"tenantId" example:"abc123"`
}
