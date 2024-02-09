package amt

import (
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/ethernetport"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/general"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/setupandconfiguration"
)

type AMTConnectionParameters struct {
	Target            string
	Username          string
	Password          string
	UseDigest         bool
	UseTLS            bool
	SelfSignedAllowed bool
}

type AMTSpecific struct {
	UUID                         string
	GeneralSettings              general.GeneralSettingsResponse
	EthernetSettings             EthernetContent
	SetupAndConfigurationService setupandconfiguration.SetupAndConfigurationServiceResponse
	Errors                       []error
}

type EthernetContent struct {
	Wired    ethernetport.SettingsResponse
	Wireless ethernetport.SettingsResponse
}
