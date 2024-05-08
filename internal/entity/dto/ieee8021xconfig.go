package dto

type IEEE8021xConfig struct {
	ProfileName            string `json:"profileName" binding:"required,max=32,alphanum" example:"My Profile"`
	AuthenticationProtocol int    `json:"authenticationProtocol" binding:"oneof=0 2 3 5 10" example:"1"`
	ServerName             string `json:"serverName,omitempty" example:"example.com"`
	Domain                 string `json:"domain,omitempty" example:"example.com"`
	Username               string `json:"username,omitempty" example:"my_username"`
	Password               string `json:"password,omitempty" example:"my_password"`
	RoamingIdentity        string `json:"roamingIdentity,omitempty" example:"my_roaming_identity"`
	ActiveInS0             bool   `json:"activeInS0,omitempty" example:"true"`
	PXETimeout             *int   `json:"pxeTimeout" binding:"required,number,gte=0,lte=86400" example:"60"`
	WiredInterface         bool   `json:"wiredInterface,omitempty" example:"false"`
	TenantID               string `json:"tenantId" example:"abc123"`
	Version                string `json:"version,omitempty" example:"1.0.0"`
}
