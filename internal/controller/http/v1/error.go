package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(c *gin.Context, err error) {
	var (
		valErr validator.ValidationErrors
		nfErr  consoleerrors.NotFoundError
		nuErr  consoleerrors.NotUniqueError
		dbErr  consoleerrors.DatabaseError
		amtErr consoleerrors.AMTError
	)

	switch {
	case errors.As(err, &valErr):
		c.AbortWithStatusJSON(http.StatusBadRequest, response{err.Error()})
	case errors.As(err, &nfErr):
		c.AbortWithStatusJSON(http.StatusNotFound, response{err.Error()})
	case errors.As(err, &nuErr):
		c.AbortWithStatusJSON(http.StatusBadRequest, response{err.Error()})
	case errors.As(err, &dbErr):
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{err.Error()})
	case errors.As(err, &amtErr):
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{err.Error()})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{"general error"})
	}
}
