package entity

type Domain struct {
	ProfileName                   string
	DomainSuffix                  string
	ProvisioningCert              string
	ProvisioningCertStorageFormat string
	ProvisioningCertPassword      string
	ExpirationDate                string
	TenantID                      string
	Version                       string
}
