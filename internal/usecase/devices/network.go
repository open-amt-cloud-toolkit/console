package devices

import (
	"context"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
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
		Wired: dto.WiredNetworkInfo{
			IEEE8021x: dto.IEEE8021x{
				Enabled:       int(response.IPSIEEE8021xSettingsResult.Enabled),
				AvailableInS0: response.IPSIEEE8021xSettingsResult.AvailableInS0,
				PxeTimeout:    response.IPSIEEE8021xSettingsResult.PxeTimeout,
			},
			NetworkInfo: dto.NetworkInfo{
				ElementName:                  response.EthernetPortSettingsResult[0].ElementName,
				InstanceID:                   response.EthernetPortSettingsResult[0].InstanceID,
				VLANTag:                      response.EthernetPortSettingsResult[0].VLANTag,
				SharedMAC:                    response.EthernetPortSettingsResult[0].SharedMAC,
				MACAddress:                   response.EthernetPortSettingsResult[0].MACAddress,
				LinkIsUp:                     response.EthernetPortSettingsResult[0].LinkIsUp,
				LinkPreference:               int(response.EthernetPortSettingsResult[0].LinkPreference),
				LinkControl:                  int(response.EthernetPortSettingsResult[0].LinkControl),
				SharedStaticIP:               response.EthernetPortSettingsResult[0].SharedStaticIp,
				SharedDynamicIP:              response.EthernetPortSettingsResult[0].SharedDynamicIP,
				IPSyncEnabled:                response.EthernetPortSettingsResult[0].IpSyncEnabled,
				DHCPEnabled:                  response.EthernetPortSettingsResult[0].DHCPEnabled,
				IPAddress:                    response.EthernetPortSettingsResult[0].IPAddress,
				SubnetMask:                   response.EthernetPortSettingsResult[0].SubnetMask,
				DefaultGateway:               response.EthernetPortSettingsResult[0].DefaultGateway,
				PrimaryDNS:                   response.EthernetPortSettingsResult[0].PrimaryDNS,
				SecondaryDNS:                 response.EthernetPortSettingsResult[0].SecondaryDNS,
				ConsoleTCPMaxRetransmissions: response.EthernetPortSettingsResult[0].ConsoleTcpMaxRetransmissions,
				WLANLinkProtectionLevel:      int(response.EthernetPortSettingsResult[0].WLANLinkProtectionLevel),
				PhysicalConnectionType:       int(response.EthernetPortSettingsResult[0].PhysicalConnectionType),
				PhysicalNicMedium:            int(response.EthernetPortSettingsResult[0].PhysicalNicMedium),
			},
		},
		Wireless: dto.WirelessNetworkInfo{
			NetworkInfo: dto.NetworkInfo{
				ElementName:                  response.EthernetPortSettingsResult[1].ElementName,
				InstanceID:                   response.EthernetPortSettingsResult[1].InstanceID,
				VLANTag:                      response.EthernetPortSettingsResult[1].VLANTag,
				SharedMAC:                    response.EthernetPortSettingsResult[1].SharedMAC,
				MACAddress:                   response.EthernetPortSettingsResult[1].MACAddress,
				LinkIsUp:                     response.EthernetPortSettingsResult[1].LinkIsUp,
				LinkPreference:               int(response.EthernetPortSettingsResult[1].LinkPreference),
				LinkControl:                  int(response.EthernetPortSettingsResult[1].LinkControl),
				SharedStaticIP:               response.EthernetPortSettingsResult[1].SharedStaticIp,
				SharedDynamicIP:              response.EthernetPortSettingsResult[1].SharedDynamicIP,
				IPSyncEnabled:                response.EthernetPortSettingsResult[1].IpSyncEnabled,
				DHCPEnabled:                  response.EthernetPortSettingsResult[1].DHCPEnabled,
				IPAddress:                    response.EthernetPortSettingsResult[1].IPAddress,
				SubnetMask:                   response.EthernetPortSettingsResult[1].SubnetMask,
				DefaultGateway:               response.EthernetPortSettingsResult[1].DefaultGateway,
				PrimaryDNS:                   response.EthernetPortSettingsResult[1].PrimaryDNS,
				SecondaryDNS:                 response.EthernetPortSettingsResult[1].SecondaryDNS,
				ConsoleTCPMaxRetransmissions: response.EthernetPortSettingsResult[1].ConsoleTcpMaxRetransmissions,
				WLANLinkProtectionLevel:      int(response.EthernetPortSettingsResult[1].WLANLinkProtectionLevel),
				PhysicalConnectionType:       int(response.EthernetPortSettingsResult[1].PhysicalConnectionType),
				PhysicalNicMedium:            int(response.EthernetPortSettingsResult[1].PhysicalNicMedium),
			},
		},
	}

	convertLinkPolicy(response, &ns)

	for _, v := range response.WiFiSettingsResult {
		ns.Wireless.WiFiNetworks = append(ns.Wireless.WiFiNetworks, dto.WiFiNetwork{
			SSID:                 v.SSID,
			AuthenticationMethod: int(v.AuthenticationMethod),
			EncryptionMethod:     int(v.EncryptionMethod),
			Priority:             v.Priority,
			BSSType:              int(v.BSSType),
		})
	}

	for i := range response.CIMIEEE8021xSettingsResult.IEEE8021xSettingsItems {
		v := &response.CIMIEEE8021xSettingsResult.IEEE8021xSettingsItems[i]
		ns.Wireless.IEEE8021xSettings = append(ns.Wireless.IEEE8021xSettings, dto.IEEE8021xSettings{
			AuthenticationProtocol:          v.AuthenticationProtocol,
			RoamingIdentity:                 v.RoamingIdentity,
			ServerCertificateName:           v.ServerCertificateName,
			ServerCertificateNameComparison: v.ServerCertificateNameComparison,
			Username:                        v.Username,
			Password:                        v.Password,
			Domain:                          v.Domain,
			ProtectedAccessCredential:       v.ProtectedAccessCredential,
		})
	}

	return ns, nil
}

func convertLinkPolicy(response wsman.NetworkResults, ns *dto.NetworkSettings) {
	for _, v := range response.EthernetPortSettingsResult[0].LinkPolicy {
		ns.Wired.NetworkInfo.LinkPolicy = append(ns.Wired.NetworkInfo.LinkPolicy, int(v))
	}

	for _, v := range response.EthernetPortSettingsResult[1].LinkPolicy {
		ns.Wireless.NetworkInfo.LinkPolicy = append(ns.Wireless.NetworkInfo.LinkPolicy, int(v))
	}
}
