package entity

type IEEE8021xConfig struct {
	ProfileName            string
	AuthenticationProtocol int
	PXETimeout             *int
	WiredInterface         bool
	TenantID               string
	Version                string
}
