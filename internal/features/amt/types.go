package amt

import (
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/ethernetport"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/general"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/setupandconfiguration"
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
	GeneralSettings              general.GeneralSettings
	EthernetSettings             EthernetContent
	SetupAndConfigurationService setupandconfiguration.Setup
	Errors                       []error
}

type EthernetContent struct {
	Wired    ethernetport.EthernetPort
	Wireless ethernetport.EthernetPort
}
