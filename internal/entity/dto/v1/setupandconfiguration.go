package dto

import (
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/setupandconfiguration"
)

type (
	SetupAndConfigurationServiceResponse struct {
		RequestedState                setupandconfiguration.RequestedState         `xml:"RequestedState,omitempty"`                // RequestedState is an integer enumeration that indicates the last requested or desired state for the element, irrespective of the mechanism through which it was requested.
		EnabledState                  setupandconfiguration.EnabledState           `xml:"EnabledState,omitempty"`                  // EnabledState is an integer enumeration that indicates the enabled and disabled states of an element.
		ElementName                   string                                       `xml:"ElementName,omitempty"`                   // A user-friendly name for the object. This property allows each instance to define a user-friendly name in addition to its key properties, identity data, and description information. Note that the Name property of ManagedSystemElement is also defined as a user-friendly name. But, it is often subclassed to be a Key. It is not reasonable that the same property can convey both identity and a user-friendly name, without inconsistencies. Where Name exists and is not a Key (such as for instances of LogicalDevice), the same information can be present in both the Name and ElementName properties. Note that if there is an associated instance of CIM_EnabledLogicalElementCapabilities, restrictions on this properties may exist as defined in ElementNameMask and MaxElementNameLen properties defined in that class.
		SystemCreationClassName       string                                       `xml:"SystemCreationClassName,omitempty"`       // The CreationClassName of the scoping System.
		SystemName                    string                                       `xml:"SystemName,omitempty"`                    // The Name of the scoping System.
		CreationClassName             string                                       `xml:"CreationClassName,omitempty"`             // CreationClassName indicates the name of the class or the subclass that is used in the creation of an instance. When used with the other key properties of this class, this property allows all instances of this class and its subclasses to be uniquely identified.
		Name                          string                                       `xml:"Name,omitempty"`                          // The Name property uniquely identifies the Service and provides an indication of the functionality that is managed. This functionality is described in more detail in the Description property of the object.
		ProvisioningMode              setupandconfiguration.ProvisioningModeValue  `xml:"ProvisioningMode,omitempty"`              // A Read-Only enumeration value that determines the behavior of Intel® AMT when it is deployed.
		ProvisioningState             setupandconfiguration.ProvisioningStateValue `xml:"ProvisioningState,omitempty"`             // An enumeration value that indicates the state of the Intel® AMT subsystem in the provisioning process"Pre" - the setup operation has not started."In" - the setup operation is in progress."Post" - Intel® AMT is configured.
		ZeroTouchConfigurationEnabled bool                                         `xml:"ZeroTouchConfigurationEnabled,omitempty"` // Indicates if Zero Touch Configuration (Remote Configuration) is enabled or disabled. This property affects only enterprise mode. It can be modified while in SMB mode
		ProvisioningServerOTP         string                                       `xml:"ProvisioningServerOTP,omitempty"`         // A optional binary data value containing 8-32 characters,that represents a one-time password (OTP), used to authenticate the Intel® AMT to the configuration server. This property can be retrieved only in IN Provisioning state, nevertheless, it is settable also in POST provisioning state.
		ConfigurationServerFQDN       string                                       `xml:"ConfigurationServerFQDN,omitempty"`       // The FQDN of the configuration server.
		PasswordModel                 setupandconfiguration.PasswordModelValue     `xml:"PasswordModel,omitempty"`                 // An enumeration value that determines the password model of Intel® AMT.
		DhcpDNSSuffix                 string                                       `xml:"DhcpDNSSuffix,omitempty"`                 // Domain name received from DHCP
		TrustedDNSSuffix              string                                       `xml:"TrustedDNSSuffix,omitempty"`              // Trusted domain name configured in MEBX
	}
)
