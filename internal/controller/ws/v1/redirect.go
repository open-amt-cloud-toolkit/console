package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type RedirectRoutes struct {
	d devices.Feature
	l logger.Interface
	u Upgrader
}

func RegisterRoutes(r *gin.Engine, l logger.Interface, t devices.Feature, u Upgrader) {
	rr := &RedirectRoutes{
		t,
		l,
		u,
	}
	r.GET("/relay/webrelay.ashx", rr.websocketHandler)
}

func (r *RedirectRoutes) websocketHandler(c *gin.Context) {
	conn, err := r.u.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	r.l.Info("Websocket connection opened")

	err = r.d.Redirect(c, conn, c.Query("host"), c.Query("mode"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redirect failed"})
		return
	}

	c.Status(http.StatusSwitchingProtocols)
}
