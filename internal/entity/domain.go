package entity

type Domain struct {
	ProfileName                   string `json:"profileName" example:"My Profile"`
	DomainSuffix                  string `json:"domainSuffix" example:"example.com"`
	ProvisioningCert              string `json:"provisioningCert" example:"-----BEGIN CERTIFICATE-----\n..."`
	ProvisioningCertStorageFormat string `json:"provisioningCertStorageFormat" example:"PKCS12"`
	ProvisioningCertPassword      string `json:"provisioningCertPassword" example:"my_password"`
	TenantID                      string `json:"tenantId" example:"abc123"`
	Version                       string `json:"version,omitempty" example:"1.0.0"`
}
