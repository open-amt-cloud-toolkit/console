package devices

import (
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
