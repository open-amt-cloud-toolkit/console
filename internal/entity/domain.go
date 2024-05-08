package entity

type Domain struct {
	ProfileName                   string
	DomainSuffix                  string
	ProvisioningCert              string
	ProvisioningCertStorageFormat string
	ProvisioningCertPassword      string
	TenantID                      string
	Version                       string
}
