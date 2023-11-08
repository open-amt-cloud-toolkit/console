package devices

import (
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/ethernetport"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/general"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/wsman/amt/setupandconfiguration"
)

var provisioningModeLookup = map[int]string{
	1: "Admin Control Mode",
	4: "Client Control Mode",
}

var provisioningStateLookup = map[int]string{
	0: "Pre-Provisioning",
	1: "In Provisioning",
	2: "Post Provisioning",
}

func CreateWsmanConnection(d Device) wsman.Messages {
	cp := wsman.ClientParameters{
		Target:            d.Address,
		Username:          d.Username,
		Password:          d.Password,
		UseDigest:         true,
		UseTLS:            d.UseTLS,
		SelfSignedAllowed: d.SelfSignedAllowed,
	}
	wsman := wsman.NewMessages(cp)
	return wsman
}

func GetGeneralSettings(wsman wsman.Messages) (gs general.GeneralSettings, err error) {
	response, err := wsman.AMT.GeneralSettings.Get()
	if err != nil {
		return
	}
	gs = response.Body.AMTGeneralSettings
	return
}

func GetEthernetSettings(wsman wsman.Messages) (ep ethernetport.EthernetPort, err error) {
	var selector ethernetport.Selector
	selector.Name = "InstanceID"
	selector.Value = "Intel(r) AMT Ethernet Port Settings 0"
	response, err := wsman.AMT.EthernetPortSettings.Get(selector)
	if err != nil {
		return
	}
	ep = response.Body.EthernetPort
	return
}

func GetSetupAndConfigurationService(wsman wsman.Messages) (sc setupandconfiguration.Setup, err error) {
	response, err := wsman.AMT.SetupAndConfigurationService.Get()
	if err != nil {
		return
	}
	sc = response.Body.Setup
	return
}
