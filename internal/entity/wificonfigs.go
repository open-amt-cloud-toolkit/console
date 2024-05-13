package entity

type WirelessConfig struct {
	ProfileName          string
	AuthenticationMethod int
	EncryptionMethod     int
	SSID                 string
	PSKValue             int
	PSKPassphrase        string
	LinkPolicy           *string
	TenantID             string
	IEEE8021xProfileName *string
	Version              string
	//	columns to populate from join query IEEE8021xProfileName
	AuthenticationProtocol *int
	ServerName             *string
	Domain                 *string
	Username               *string
	Password               *string
	RoamingIdentity        *string
	ActiveInS0             *bool
	PXETimeout             *int
	WiredInterface         *bool
}
