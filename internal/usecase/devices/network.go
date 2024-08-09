package devices

import (
	"context"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
)

func (uc *UseCase) GetNetworkSettings(c context.Context, guid string) (dto.NetworkSettings, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return dto.NetworkSettings{}, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	response, err := device.GetNetworkSettings()
	if err != nil {
		return dto.NetworkSettings{}, err
	}

	ns := dto.NetworkSettings{
		Wired: dto.NetworkInfo{
			ElementName: response.EthernetPortSettingsResult[0].ElementName,
			InstanceID:  response.EthernetPortSettingsResult[0].InstanceID,
			VLANTag:     response.EthernetPortSettingsResult[0].VLANTag,
			SharedMAC:   response.EthernetPortSettingsResult[0].SharedMAC,
			MACAddress:  response.EthernetPortSettingsResult[0].MACAddress,
			LinkIsUp:    response.EthernetPortSettingsResult[0].LinkIsUp,
			//LinkPolicy:                   response.EthernetPortSettingsResult[0].LinkPolicy,
			//LinkPreference:               response.EthernetPortSettingsResult[0].LinkPreference,
			//LinkControl:                  response.EthernetPortSettingsResult[0].LinkControl,
			SharedStaticIp:               response.EthernetPortSettingsResult[0].SharedStaticIp,
			SharedDynamicIP:              response.EthernetPortSettingsResult[0].SharedDynamicIP,
			IpSyncEnabled:                response.EthernetPortSettingsResult[0].IpSyncEnabled,
			DHCPEnabled:                  response.EthernetPortSettingsResult[0].DHCPEnabled,
			IPAddress:                    response.EthernetPortSettingsResult[0].IPAddress,
			SubnetMask:                   response.EthernetPortSettingsResult[0].SubnetMask,
			DefaultGateway:               response.EthernetPortSettingsResult[0].DefaultGateway,
			PrimaryDNS:                   response.EthernetPortSettingsResult[0].PrimaryDNS,
			SecondaryDNS:                 response.EthernetPortSettingsResult[0].SecondaryDNS,
			ConsoleTcpMaxRetransmissions: response.EthernetPortSettingsResult[0].ConsoleTcpMaxRetransmissions,
			WLANLinkProtectionLevel:      int(response.EthernetPortSettingsResult[0].WLANLinkProtectionLevel),
			PhysicalConnectionType:       int(response.EthernetPortSettingsResult[0].PhysicalConnectionType),
			PhysicalNicMedium:            int(response.EthernetPortSettingsResult[0].PhysicalNicMedium),
		},
		Wireless: dto.NetworkInfo{
			ElementName: response.EthernetPortSettingsResult[1].ElementName,
			InstanceID:  response.EthernetPortSettingsResult[1].InstanceID,
			VLANTag:     response.EthernetPortSettingsResult[1].VLANTag,
			SharedMAC:   response.EthernetPortSettingsResult[1].SharedMAC,
			MACAddress:  response.EthernetPortSettingsResult[1].MACAddress,
			LinkIsUp:    response.EthernetPortSettingsResult[1].LinkIsUp,
			//LinkPolicy:                   response.EthernetPortSettingsResult[1].LinkPolicy,
			//LinkPreference:               response.EthernetPortSettingsResult[1].LinkPreference,
			//LinkControl:                  response.EthernetPortSettingsResult[1].LinkControl,
			SharedStaticIp:               response.EthernetPortSettingsResult[1].SharedStaticIp,
			SharedDynamicIP:              response.EthernetPortSettingsResult[1].SharedDynamicIP,
			IpSyncEnabled:                response.EthernetPortSettingsResult[1].IpSyncEnabled,
			DHCPEnabled:                  response.EthernetPortSettingsResult[1].DHCPEnabled,
			IPAddress:                    response.EthernetPortSettingsResult[1].IPAddress,
			SubnetMask:                   response.EthernetPortSettingsResult[1].SubnetMask,
			DefaultGateway:               response.EthernetPortSettingsResult[1].DefaultGateway,
			PrimaryDNS:                   response.EthernetPortSettingsResult[1].PrimaryDNS,
			SecondaryDNS:                 response.EthernetPortSettingsResult[1].SecondaryDNS,
			ConsoleTcpMaxRetransmissions: response.EthernetPortSettingsResult[1].ConsoleTcpMaxRetransmissions,
			WLANLinkProtectionLevel:      int(response.EthernetPortSettingsResult[1].WLANLinkProtectionLevel),
			PhysicalConnectionType:       int(response.EthernetPortSettingsResult[1].PhysicalConnectionType),
			PhysicalNicMedium:            int(response.EthernetPortSettingsResult[1].PhysicalNicMedium),
		},
	}

	return ns, nil
}
