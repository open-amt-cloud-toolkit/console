package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/sqldb"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(c *gin.Context, err error) {
	var (
		valErr validator.ValidationErrors
		nfErr  sqldb.NotFoundError
		nuErr  sqldb.NotUniqueError
		dbErr  sqldb.DatabaseError
		amtErr devices.AMTError
	)

	switch {
	case errors.As(err, &valErr):
		c.AbortWithStatusJSON(http.StatusBadRequest, response{err.Error()})
	case errors.As(err, &nfErr):
		c.AbortWithStatusJSON(http.StatusNotFound, response{nfErr.Console.FriendlyMessage()})
	case errors.As(err, &nuErr):
		c.AbortWithStatusJSON(http.StatusBadRequest, response{nuErr.Console.FriendlyMessage()})
	case errors.As(err, &dbErr):
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{dbErr.Console.FriendlyMessage()})
	case errors.As(err, &amtErr):
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{amtErr.Console.FriendlyMessage()})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{"general error"})
	}
}
