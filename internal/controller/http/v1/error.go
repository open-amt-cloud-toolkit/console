package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/domains"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/sqldb"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(c *gin.Context, err error) {
	var (
		valErr          validator.ValidationErrors
		nfErr           sqldb.NotFoundError
		nuErr           sqldb.NotUniqueError
		dbErr           sqldb.DatabaseError
		amtErr          devices.AMTError
		certExpErr      domains.CertExpirationError
		certPasswordErr domains.CertPasswordError
	)

	switch {
	case errors.As(err, &valErr):
		c.AbortWithStatusJSON(http.StatusBadRequest, response{err.Error()})
	case errors.As(err, &certExpErr):
		c.AbortWithStatusJSON(http.StatusBadRequest, response{certExpErr.Console.FriendlyMessage()})
	case errors.As(err, &certPasswordErr):
		c.AbortWithStatusJSON(http.StatusBadRequest, response{certPasswordErr.Console.FriendlyMessage()})
	case errors.As(err, &nfErr):
		c.AbortWithStatusJSON(http.StatusNotFound, response{nfErr.Console.FriendlyMessage()})
	case errors.As(err, &nuErr):
		c.AbortWithStatusJSON(http.StatusBadRequest, response{nuErr.Console.FriendlyMessage()})
	case errors.As(err, &dbErr):
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{dbErr.Console.FriendlyMessage()})
	case errors.As(err, &amtErr):
		if strings.Contains(amtErr.Console.Error(), "400 Bad Request") {
			c.AbortWithStatusJSON(http.StatusBadRequest, response{amtErr.Console.FriendlyMessage()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response{amtErr.Console.FriendlyMessage()})
		}
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{"general error"})
	}
}
