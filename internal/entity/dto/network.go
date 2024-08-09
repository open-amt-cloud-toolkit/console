package dto

// NetworkSettings defines the network settings for a device.
type NetworkSettings struct {
	Wired    NetworkInfo `json:"wired"`
	Wireless NetworkInfo `json:"wireless"`
}

// NetworkResults defines the network results for a device.
type NetworkInfo struct {
	ElementName                  string `json:"elementName"`                  // The user-friendly name for this instance of SettingData. In addition, the user-friendly name can be used as an index property for a search or query. (Note: The name does not have to be unique within a namespace.)
	InstanceID                   string `json:"instanceID"`                   // Within the scope of the instantiating Namespace, InstanceID opaquely and uniquely identifies an instance of this class.
	VLANTag                      int    `json:"vlanTag"`                      // Indicates whether VLAN is in use and what is the VLAN tag when used.
	SharedMAC                    bool   `json:"sharedMAC"`                    // Indicates whether Intel® AMT shares it's MAC address with the host system.
	MACAddress                   string `json:"macAddress"`                   // The MAC address used by Intel® AMT in a string format. For Example: 01-02-3f-b0-99-99. (This property can only be read and can't be changed.)
	LinkIsUp                     bool   `json:"linkIsUp"`                     // Indicates whether the network link is up
	LinkPolicy                   []int  `json:"linkPolicy"`                   // Enumeration values for link policy restrictions for better power consumption. If Intel® AMT will not be able to determine the exact power state, the more restrictive closest configuration applies.
	LinkPreference               int    `json:"linkPreference"`               // Determines whether the link is preferred to be owned by ME or host
	LinkControl                  int    `json:"linkControl"`                  // Determines whether the link is owned by ME or host.  Additional Notes: This property is read-only.
	SharedStaticIp               bool   `json:"sharedStaticIP"`               // Indicates whether the static host IP is shared with ME.
	SharedDynamicIP              bool   `json:"sharedDynamicIP"`              // Indicates whether the dynamic host IP is shared with ME. This property is read only.
	IpSyncEnabled                bool   `json:"ipSyncEnabled"`                // Indicates whether the IP synchronization between host and ME is enabled.
	DHCPEnabled                  bool   `json:"dhcpEnabled"`                  // Indicates whether DHCP is in use. Additional Notes: 'DHCPEnabled' is a required field for the Put command.
	IPAddress                    string `json:"ipAddress"`                    // String representation of IP address. Get operation - reports the acquired IP address (whether in static or DHCP mode). Put operation - sets the IP address (in static mode only).
	SubnetMask                   string `json:"subnetMask"`                   // Subnet mask in a string format.For example: 255.255.0.0
	DefaultGateway               string `json:"defaultGateway"`               // Default Gateway in a string format. For example: 10.12.232.1
	PrimaryDNS                   string `json:"primaryDNS"`                   // Primary DNS in a string format. For example: 10.12.232.1
	SecondaryDNS                 string `json:"secondaryDNS"`                 // Secondary DNS in a string format. For example: 10.12.232.1
	ConsoleTcpMaxRetransmissions int    `json:"consoleTCPMaxRetransmissions"` // Indicates the number of retransmissions host TCP SW tries ifno ack is accepted
	WLANLinkProtectionLevel      int    `json:"wlanLinkProtectionLevel"`      // Defines the level of the link protection feature activation. Read only property.
	PhysicalConnectionType       int    `json:"physicalConnectionType"`       // Indicates the physical connection type of this network interface. Note: Applicable in Intel AMT 15.0 and later.
	PhysicalNicMedium            int    `json:"physicalNICMedium"`            // Indicates which medium is currently used by Intel® AMT to communicate with the NIC. Note: Applicable in Intel AMT 15.0 and later.
}
