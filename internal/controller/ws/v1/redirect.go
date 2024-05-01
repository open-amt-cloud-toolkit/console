package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	Subprotocols:    []string{"direct"},
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
	EnableCompression: false,
}

type RedirectRoutes struct {
	d devices.Feature
	l logger.Interface
}

func RegisterRoutes(r *gin.Engine, l logger.Interface, t devices.Feature) {
	rr := &RedirectRoutes{
		t,
		l,
	}
	r.GET("/relay/webrelay.ashx", rr.websocketHandler)
}

func (r *RedirectRoutes) websocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.Error(c.Writer, "Could not open websocket connection", http.StatusInternalServerError)

		return
	}

	r.l.Info("Websocket connection opened")

	host := c.Query("host") // guid
	mode := c.Query("mode")

	err = r.d.Redirect(c, conn, host, mode)
	if err != nil {
		r.l.Error(err, "http - devices - v1 - redirect")
		errorResponse(c, http.StatusInternalServerError, "redirect failed")
	}
}
