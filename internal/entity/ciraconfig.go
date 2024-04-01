package entity

type CIRAConfig struct {
	ConfigName          string `json:"configName" example:"My CIRA Config"`
	MPSServerAddress    string `json:"mpsServerAddress" example:"https://example.com"`
	MpsPort             int    `json:"mpsPort" example:"443"`
	Username            string `json:"username" example:"my_username"`
	Password            string `json:"password,omitempty" example:"my_password"`
	CommonName          string `json:"commonName" example:"example.com"`
	ServerAddressFormat int    `json:"serverAddressFormat" example:"201"`
	AuthMethod          int    `json:"authMethod" example:"2"`
	MpsRootCertificate  string `json:"mpsRootCertificate" example:"-----BEGIN CERTIFICATE-----\n..."`
	ProxyDetails        string `json:"proxyDetails" example:"http://example.com"`
	TenantID            string `json:"tenantId" example:"abc123"`
	RegeneratePassword  bool   `json:"regeneratePassword,omitempty" example:"true"`
	Version             string `json:"version,omitempty" example:"1.0.0"`
}
