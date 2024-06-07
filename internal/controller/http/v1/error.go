package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/sqldb"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(c *gin.Context, err error) {
	var (
		valErr dto.NotValidError
		nfErr  sqldb.NotFoundError
		nuErr  sqldb.NotUniqueError
		dbErr  sqldb.DatabaseError
		amtErr devices.AMTError
	)

	switch {
	case errors.As(err, &valErr):
		c.AbortWithStatusJSON(http.StatusBadRequest, response{valErr.Console.FriendlyMessage()})
	case errors.As(err, &nfErr):
		c.AbortWithStatusJSON(http.StatusNotFound, response{nfErr.Console.FriendlyMessage()})
	case errors.As(err, &dbErr):
		if errors.As(dbErr.Console.OriginalError, &nuErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, response{nuErr.Console.FriendlyMessage()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response{dbErr.Console.FriendlyMessage()})
		}
	case errors.As(err, &amtErr):
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{amtErr.Console.FriendlyMessage()})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, response{"general error"})
	}
}
