package amtexplorer

import (
	"sync"
	"time"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/client"

	dtov1 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	wsmanAPI "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

var (
	connections   = make(map[string]*wsmanAPI.ConnectionEntry)
	connectionsMu sync.Mutex
	expireAfter   = 5 * time.Minute // Set the expiration duration as needed
)

type GoWSMANMessages struct {
	log logger.Interface
}

func NewGoWSMANMessages(log logger.Interface) *GoWSMANMessages {
	return &GoWSMANMessages{
		log: log,
	}
}

func (g GoWSMANMessages) DestroyWsmanClient(device dtov1.Device) {
	if entry, ok := connections[device.GUID]; ok {
		entry.Timer.Stop()
		removeConnection(device.GUID)
	}
}

func (g GoWSMANMessages) SetupWsmanClient(device dtov1.Device, logAMTMessages bool) AMTExplorer {
	clientParams := client.Parameters{
		Target:            device.Hostname,
		Username:          device.Username,
		Password:          device.Password,
		UseDigest:         true,
		UseTLS:            device.UseTLS,
		SelfSignedAllowed: device.AllowSelfSigned,
		LogAMTMessages:    logAMTMessages,
		IsRedirection:     false,
	}

	if device.CertHash != "" {
		clientParams.PinnedCert = device.CertHash
	}

	connectionsMu.Lock()
	defer connectionsMu.Unlock()

	if entry, ok := connections[device.GUID]; ok {
		entry.Timer.Stop() // Stop the previous timer
		entry.Timer = time.AfterFunc(expireAfter, func() {
			removeConnection(device.GUID)
		})
	} else {
		wsmanMsgs := wsman.NewMessages(clientParams)
		timer := time.AfterFunc(expireAfter, func() {
			removeConnection(device.GUID)
		})
		connections[device.GUID] = &wsmanAPI.ConnectionEntry{
			WsmanMessages: wsmanMsgs,
			Timer:         timer,
		}
	}

	return connections[device.GUID]
}

func removeConnection(guid string) {
	connectionsMu.Lock()
	defer connectionsMu.Unlock()

	delete(connections, guid)
}
