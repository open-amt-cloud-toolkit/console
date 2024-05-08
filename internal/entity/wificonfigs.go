package entity

type ProfileWifiConfigs struct {
	Priority    int
	ProfileName string
	TenantID    string
}

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
}
