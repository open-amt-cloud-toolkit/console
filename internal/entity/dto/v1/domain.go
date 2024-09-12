package dtov1

import "time"

type Domain struct {
	ProfileName                   string    `json:"profileName" binding:"required,alphanum" example:"My Profile"`
	DomainSuffix                  string    `json:"domainSuffix" binding:"required" example:"example.com"`
	ProvisioningCert              string    `json:"provisioningCert,omitempty" binding:"required" example:"-----BEGIN CERTIFICATE-----\n..."`
	ProvisioningCertStorageFormat string    `json:"provisioningCertStorageFormat" binding:"required,oneof=raw string" example:"string"`
	ProvisioningCertPassword      string    `json:"provisioningCertPassword,omitempty" binding:"required,lte=64" example:"my_password"`
	ExpirationDate                time.Time `json:"expirationDate,omitempty" example:"2022-01-01T00:00:00Z"`
	TenantID                      string    `json:"tenantId" example:"abc123"`
	Version                       string    `json:"version,omitempty" example:"1.0.0"`
}
