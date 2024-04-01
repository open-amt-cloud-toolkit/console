package entity

type IEEE8021xConfig struct {
	ProfileName            string `json:"profileName" example:"My Profile"`
	AuthenticationProtocol int    `json:"authenticationProtocol" example:"1"`
	ServerName             string `json:"serverName,omitempty" example:"example.com"`
	Domain                 string `json:"domain,omitempty" example:"example.com"`
	Username               string `json:"username,omitempty" example:"my_username"`
	Password               string `json:"password,omitempty" example:"my_password"`
	RoamingIdentity        string `json:"roamingIdentity,omitempty" example:"my_roaming_identity"`
	ActiveInS0             bool   `json:"activeInS0,omitempty" example:"true"`
	PxeTimeout             int    `json:"pxeTimeout" example:"60"`
	WiredInterface         bool   `json:"wiredInterface,omitempty" example:"false"`
	TenantID               string `json:"tenantId" example:"abc123"`
	Version                string `json:"version,omitempty" example:"1.0.0"`
}
