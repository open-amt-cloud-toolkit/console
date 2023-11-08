package devices

import (
	"html/template"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/ethernetport"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/general"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/setupandconfiguration"
	"go.etcd.io/bbolt"
)

type DeviceThing struct {
	db *bbolt.DB
	//parsed templates
	html *template.Template
}

type DeviceContent struct {
	Device                       Device
	GeneralSettings              general.GeneralSettings
	EthernetPort                 ethernetport.EthernetPort
	SetupAndConfigurationService setupandconfiguration.Setup
}
