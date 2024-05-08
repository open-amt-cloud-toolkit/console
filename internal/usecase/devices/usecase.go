package devices

import (
	"strings"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

const (
	RedirectionCommandsStartRedirectionSession      = 16
	RedirectionCommandsStartRedirectionSessionReply = 17
	RedirectionCommandsEndRedirectionSession        = 18
	RedirectionCommandsAuthenticateSession          = 19
	RedirectionCommandsAuthenticateSessionReply     = 20
	StartRedirectionSessionReplyStatusSuccess       = 0
	StartRedirectionSessionReplyStatusUnknown       = 1
	StartRedirectionSessionReplyStatusBusy          = 2
	StartRedirectionSessionReplyStatusUnsupported   = 3
	StartRedirectionSessionReplyStatusError         = 0xFF
	AuthenticationTypeQuery                         = 0
	AuthenticationTypeUserPass                      = 1
	AuthenticationTypeKerberos                      = 2
	AuthenticationTypeBadDigest                     = 3
	AuthenticationTypeDigest                        = 4
	AuthenticationStatusSuccess                     = 0
	AuthenticationStatusFail                        = 1
	AuthenticationStatusNotSupported                = 2

	// MinAMTVersion - minimum AMT version required for certain features in power capabilities.
	MinAMTVersion = 9
)

// UseCase -.
type UseCase struct {
	repo             Repository
	device           Management
	redirection      Redirection
	redirConnections map[string]*DeviceConnection
	log              logger.Interface
}

// New -.
func New(r Repository, d Management, redirection Redirection, log logger.Interface) *UseCase {
	return &UseCase{
		repo:             r,
		device:           d,
		redirection:      redirection,
		redirConnections: make(map[string]*DeviceConnection),
		log:              log,
	}
}

// convert dto.Device to entity.Device.
func (uc *UseCase) dtoToEntity(d *dto.Device) *entity.Device {
	// convert []string to comma separated string
	if d.Tags == nil {
		d.Tags = []string{}
	}

	tags := strings.Join(d.Tags, ", ")

	d1 := &entity.Device{
		ConnectionStatus: d.ConnectionStatus,
		MPSInstance:      d.MPSInstance,
		Hostname:         d.Hostname,
		GUID:             d.GUID,
		MPSUsername:      d.MPSUsername,
		Tags:             tags,
		TenantID:         d.TenantID,
		FriendlyName:     d.FriendlyName,
		DNSSuffix:        d.DNSSuffix,
		LastConnected:    d.LastConnected,
		LastSeen:         d.LastSeen,
		LastDisconnected: d.LastDisconnected,
		// DeviceInfo:       d.DeviceInfo,
		Username:        d.Username,
		Password:        d.Password,
		UseTLS:          d.UseTLS,
		AllowSelfSigned: d.AllowSelfSigned,
	}

	return d1
}

// convert entity.Device to dto.Device.
func (uc *UseCase) entityToDTO(d *entity.Device) *dto.Device {
	// convert comma separated string to []string
	tags := strings.Split(d.Tags, ",")

	d1 := &dto.Device{
		ConnectionStatus: d.ConnectionStatus,
		MPSInstance:      d.MPSInstance,
		Hostname:         d.Hostname,
		GUID:             d.GUID,
		MPSUsername:      d.MPSUsername,
		Tags:             tags,
		TenantID:         d.TenantID,
		FriendlyName:     d.FriendlyName,
		DNSSuffix:        d.DNSSuffix,
		LastConnected:    d.LastConnected,
		LastSeen:         d.LastSeen,
		LastDisconnected: d.LastDisconnected,
		// DeviceInfo:       d.DeviceInfo,
		Username:        d.Username,
		Password:        d.Password,
		UseTLS:          d.UseTLS,
		AllowSelfSigned: d.AllowSelfSigned,
	}

	return d1
}
