package devices

import (
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/ethernetport"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/general"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/setupandconfiguration"
)

type DeviceContent struct {
	Device                       Device
	UUID                         string
	GeneralSettings              general.GeneralSettings
	EthernetPort                 []ethernetport.EthernetPort
	SetupAndConfigurationService setupandconfiguration.Setup
}
