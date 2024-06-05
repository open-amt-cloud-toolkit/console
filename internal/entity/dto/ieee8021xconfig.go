package dto

import "github.com/go-playground/validator/v10"

type IEEE8021xConfig struct {
	ProfileName            string `json:"profileName" binding:"required,max=32,alphanum" example:"My Profile"`
	AuthenticationProtocol int    `json:"authenticationProtocol" binding:"matchAuthProtocol" example:"1"`
	PXETimeout             *int   `json:"pxeTimeout" binding:"required,number,gte=0,lte=86400" example:"60"`
	WiredInterface         bool   `json:"wiredInterface,omitempty" example:"false"`
	TenantID               string `json:"tenantId" example:"abc123"`
	Version                string `json:"version,omitempty" example:"1.0.0"`
}

func AuthProtocolValidator(fl validator.FieldLevel) bool {
	config := fl.Parent().Interface().(IEEE8021xConfig)
	authProtocol := config.AuthenticationProtocol

	if config.WiredInterface {
		return validator.New().Var(authProtocol, "oneof=0 2 3 5 10") == nil
	}

	return validator.New().Var(authProtocol, "oneof=0 2") == nil
}
