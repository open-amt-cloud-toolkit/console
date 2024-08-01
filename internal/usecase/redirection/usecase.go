package redirection

import "github.com/open-amt-cloud-toolkit/console/pkg/logger"

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
)

// UseCase -.
type UseCase struct {
	repo             Repository
	redirection      Feature
	redirConnections map[string]*DeviceConnection
	log              logger.Interface
}

// New -.
func New(r Repository, redirection Feature, log logger.Interface) *UseCase {
	uc := &UseCase{
		repo:             r,
		redirection:      redirection,
		redirConnections: make(map[string]*DeviceConnection),
		log:              log,
	}

	return uc
}
