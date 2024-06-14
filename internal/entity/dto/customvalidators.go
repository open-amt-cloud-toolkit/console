package dto

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func SetupCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("matchFormat", AddressFormatValidator)
		v.RegisterValidation("matchAuthProtocol", AuthProtocolValidator)
	}
}
