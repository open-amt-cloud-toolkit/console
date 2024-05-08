package entity

type CIRAConfig struct {
	ConfigName          string
	MPSServerAddress    string
	MpsPort             int
	Username            string
	Password            string
	CommonName          string
	ServerAddressFormat int
	AuthMethod          int
	MpsRootCertificate  string
	ProxyDetails        string
	TenantID            string
	RegeneratePassword  bool
	Version             string
}
