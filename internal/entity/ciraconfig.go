package entity

type CIRAConfig struct {
	ConfigName          string
	MPSAddress          string
	MPSPort             int
	Username            string
	Password            string
	CommonName          string
	ServerAddressFormat int
	AuthMethod          int
	MPSRootCertificate  string
	ProxyDetails        string
	TenantID            string
	RegeneratePassword  bool
	Version             string
}
