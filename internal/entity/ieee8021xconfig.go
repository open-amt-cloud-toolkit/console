package entity

type IEEE8021xConfig struct {
	ProfileName            string
	AuthenticationProtocol int
	ServerName             string
	Domain                 string
	Username               string
	Password               string
	RoamingIdentity        string
	ActiveInS0             bool
	PXETimeout             *int
	WiredInterface         bool
	TenantID               string
	Version                string
}
