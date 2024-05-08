package dto

import (
	"github.com/go-playground/validator/v10"
)

type Domain struct {
	ProfileName                   string `json:"profileName" binding:"required,alphanum" example:"My Profile"`
	DomainSuffix                  string `json:"domainSuffix" binding:"required" example:"example.com"`
	ProvisioningCert              string `json:"provisioningCert,omitempty" binding:"required" example:"-----BEGIN CERTIFICATE-----\n..."`
	ProvisioningCertStorageFormat string `json:"provisioningCertStorageFormat" binding:"required,storageformat" example:"PKCS12"`
	ProvisioningCertPassword      string `json:"provisioningCertPassword,omitempty" binding:"required,lte=64" example:"my_password"`
	TenantID                      string `json:"tenantId" example:"abc123"`
	Version                       string `json:"version,omitempty" example:"1.0.0"`
}

var StorageFormatValidation validator.Func = func(fl validator.FieldLevel) bool {
	provisioningCertStorageFormat, ok := fl.Field().Interface().(string)
	if ok {
		if provisioningCertStorageFormat != "raw" && provisioningCertStorageFormat != "string" {
			return false
		}
	}

	return true
}
