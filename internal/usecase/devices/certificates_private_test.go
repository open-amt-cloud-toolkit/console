package devices

import (
	"encoding/xml"
	"testing"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/publickey"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/publicprivate"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/concrete"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/credential"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/models"
	"github.com/stretchr/testify/require"

	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
)

type certificateTest struct {
	name string
	res  SecuritySettings
	err  error

	profileType string
}

var getResponse wsman.Certificates = wsman.Certificates{
	ConcreteDependencyResponse: concrete.PullResponse{
		XMLName: xml.Name{
			Space: "http://schemas.xmlsoap.org/ws/2004/09/enumeration",
			Local: "PullResponse",
		},
		Items: []concrete.ConcreteDependency{
			{
				Antecedent: models.AssociationReference{
					Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
					ReferenceParameters: models.ReferenceParametersNoNamespace{
						XMLName: xml.Name{
							Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
							Local: "ReferenceParameters",
						},
						ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_PublicKeyCertificate",
						SelectorSet: models.SelectorNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
								Local: "SelectorSet",
							},
							Selectors: []models.SelectorResponse{
								{
									XMLName: xml.Name{
										Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
										Local: "Selector",
									},
									Name: "InstanceID",
									Text: "Intel(r) AMT Certificate: Handle: 0",
								},
							},
						},
					},
				},
				Dependent: models.AssociationReference{
					Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
					ReferenceParameters: models.ReferenceParametersNoNamespace{
						XMLName: xml.Name{
							Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
							Local: "ReferenceParameters",
						},
						ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_PublicPrivateKeyPair",
						SelectorSet: models.SelectorNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
								Local: "SelectorSet",
							},
							Selectors: []models.SelectorResponse{
								{
									XMLName: xml.Name{
										Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
										Local: "Selector",
									},
									Name: "InstanceID",
									Text: "Intel(r) AMT Key: Handle: 0",
								},
							},
						},
					},
				},
			},
			{
				Antecedent: models.AssociationReference{
					Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
					ReferenceParameters: models.ReferenceParametersNoNamespace{
						XMLName: xml.Name{
							Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
							Local: "ReferenceParameters",
						},
						ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_PublicKeyCertificate",
						SelectorSet: models.SelectorNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
								Local: "SelectorSet",
							},
							Selectors: []models.SelectorResponse{
								{
									XMLName: xml.Name{
										Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
										Local: "Selector",
									},
									Name: "InstanceID",
									Text: "Intel(r) AMT Certificate: Handle: 1",
								},
							},
						},
					},
				},
				Dependent: models.AssociationReference{
					Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
					ReferenceParameters: models.ReferenceParametersNoNamespace{
						XMLName: xml.Name{
							Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
							Local: "ReferenceParameters",
						},
						ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_PublicPrivateKeyPair",
						SelectorSet: models.SelectorNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
								Local: "SelectorSet",
							},
							Selectors: []models.SelectorResponse{
								{
									XMLName: xml.Name{
										Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
										Local: "Selector",
									},
									Name: "InstanceID",
									Text: "Intel(r) AMT Key: Handle: 1",
								},
							},
						},
					},
				},
			},
			{
				Antecedent: models.AssociationReference{
					Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
					ReferenceParameters: models.ReferenceParametersNoNamespace{
						XMLName: xml.Name{
							Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
							Local: "ReferenceParameters",
						},
						ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_PublicKeyCertificate",
						SelectorSet: models.SelectorNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
								Local: "SelectorSet",
							},
							Selectors: []models.SelectorResponse{
								{
									XMLName: xml.Name{
										Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
										Local: "Selector",
									},
									Name: "InstanceID",
									Text: "Intel(r) AMT Certificate: Handle: 3",
								},
							},
						},
					},
				},
				Dependent: models.AssociationReference{
					Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
					ReferenceParameters: models.ReferenceParametersNoNamespace{
						XMLName: xml.Name{
							Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
							Local: "ReferenceParameters",
						},
						ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_PublicPrivateKeyPair",
						SelectorSet: models.SelectorNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
								Local: "SelectorSet",
							},
							Selectors: []models.SelectorResponse{
								{
									XMLName: xml.Name{
										Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
										Local: "Selector",
									},
									Name: "InstanceID",
									Text: "Intel(r) AMT Key: Handle: 2",
								},
							},
						},
					},
				},
			},
			{
				Antecedent: models.AssociationReference{
					Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
					ReferenceParameters: models.ReferenceParametersNoNamespace{
						XMLName: xml.Name{
							Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
							Local: "ReferenceParameters",
						},
						ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_PublicKeyCertificate",
						SelectorSet: models.SelectorNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
								Local: "SelectorSet",
							},
							Selectors: []models.SelectorResponse{
								{
									XMLName: xml.Name{
										Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
										Local: "Selector",
									},
									Name: "InstanceID",
									Text: "Intel(r) AMT Certificate: Handle: 4",
								},
							},
						},
					},
				},
				Dependent: models.AssociationReference{
					Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
					ReferenceParameters: models.ReferenceParametersNoNamespace{
						XMLName: xml.Name{
							Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
							Local: "ReferenceParameters",
						},
						ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_PublicPrivateKeyPair",
						SelectorSet: models.SelectorNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
								Local: "SelectorSet",
							},
							Selectors: []models.SelectorResponse{
								{
									XMLName: xml.Name{
										Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
										Local: "Selector",
									},
									Name: "InstanceID",
									Text: "Intel(r) AMT Key: Handle: 3",
								},
							},
						},
					},
				},
			},
		},
	},
	PublicKeyCertificateResponse: publickey.RefinedPullResponse{
		KeyManagementItems: []publickey.RefinedKeyManagementResponse{},
		PublicKeyCertificateItems: []publickey.RefinedPublicKeyCertificateResponse{
			{
				ElementName:            "Intel(r) AMT Certificate",
				InstanceID:             "Intel(r) AMT Certificate: Handle: 2",
				X509Certificate:        "TestCertRoot",
				TrustedRootCertificate: true,
				Issuer:                 "TestIssuer",
				Subject:                "TestSubject",
				ReadOnlyCertificate:    false,
				PublicKeyHandle:        "",
				AssociatedProfiles:     nil,
			},
			{
				ElementName:            "Intel(r) AMT Certificate",
				InstanceID:             "Intel(r) AMT Certificate: Handle: 0",
				X509Certificate:        "TestCert0",
				TrustedRootCertificate: false,
				Issuer:                 "TestIssuer",
				Subject:                "TestSubject2",
				ReadOnlyCertificate:    false,
				PublicKeyHandle:        "",
				AssociatedProfiles:     nil,
			},
			{
				ElementName:            "Intel(r) AMT Certificate",
				InstanceID:             "Intel(r) AMT Certificate: Handle: 1",
				X509Certificate:        "TestCert1",
				TrustedRootCertificate: false,
				Issuer:                 "TestIssuer",
				Subject:                "TestSubject2",
				ReadOnlyCertificate:    false,
				PublicKeyHandle:        "",
				AssociatedProfiles:     nil,
			},
			{
				ElementName:            "Intel(r) AMT Certificate",
				InstanceID:             "Intel(r) AMT Certificate: Handle: 3",
				X509Certificate:        "TestCert3",
				TrustedRootCertificate: false,
				Issuer:                 "TestIssuer",
				Subject:                "TestSubject2",
				ReadOnlyCertificate:    false,
				PublicKeyHandle:        "",
				AssociatedProfiles:     nil,
			},
		},
	},
	PublicPrivateKeyPairResponse: publicprivate.RefinedPullResponse{
		PublicPrivateKeyPairItems: []publicprivate.RefinedPublicPrivateKeyPair{
			{
				ElementName:       "Intel(r) AMT Key",
				InstanceID:        "Intel(r) AMT Key: Handle: 0",
				DERKey:            "Key0",
				CertificateHandle: "",
			},
			{
				ElementName:       "Intel(r) AMT Key",
				InstanceID:        "Intel(r) AMT Key: Handle: 1",
				DERKey:            "Key1",
				CertificateHandle: "",
			},
			{
				ElementName:       "Intel(r) AMT Key",
				InstanceID:        "Intel(r) AMT Key: Handle: 2",
				DERKey:            "Key2",
				CertificateHandle: "",
			},
		},
	},
	CIMCredentialContextResponse: credential.PullResponse{
		XMLName: xml.Name{
			Space: "http://schemas.xmlsoap.org/ws/2004/09/enumeration",
			Local: "PullResponse",
		},
		Items: credential.Items{
			CredentialContext: []credential.CredentialContext{
				{
					ElementInContext: models.AssociationReference{
						Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
						ReferenceParameters: models.ReferenceParametersNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
								Local: "ReferenceParameters",
							},
							ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_PublicKeyCertificate",
							SelectorSet: models.SelectorNoNamespace{
								XMLName: xml.Name{
									Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
									Local: "SelectorSet",
								},
								Selectors: []models.SelectorResponse{
									{
										XMLName: xml.Name{
											Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
											Local: "Selector",
										},
										Name: "InstanceID",
										Text: "Intel(r) AMT Certificate: Handle: 2",
									},
								},
							},
						},
					},
					ElementProvidingContext: models.AssociationReference{
						Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
						ReferenceParameters: models.ReferenceParametersNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
								Local: "ReferenceParameters",
							},
							ResourceURI: "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_IEEE8021xSettings",
							SelectorSet: models.SelectorNoNamespace{
								XMLName: xml.Name{
									Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
									Local: "SelectorSet",
								},
								Selectors: []models.SelectorResponse{
									{
										XMLName: xml.Name{
											Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
											Local: "Selector",
										},
										Name: "InstanceID",
										Text: "Intel(r) AMT:IEEE 802.1x Settings TestWifi8021xTLS",
									},
								},
							},
						},
					},
				},
				{
					ElementInContext: models.AssociationReference{
						Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
						ReferenceParameters: models.ReferenceParametersNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
								Local: "ReferenceParameters",
							},
							ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_PublicKeyCertificate",
							SelectorSet: models.SelectorNoNamespace{
								XMLName: xml.Name{
									Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
									Local: "SelectorSet",
								},
								Selectors: []models.SelectorResponse{
									{
										XMLName: xml.Name{
											Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
											Local: "Selector",
										},
										Name: "InstanceID",
										Text: "Intel(r) AMT Certificate: Handle: 3",
									},
								},
							},
						},
					},
					ElementProvidingContext: models.AssociationReference{
						Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
						ReferenceParameters: models.ReferenceParametersNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
								Local: "ReferenceParameters",
							},
							ResourceURI: "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_IEEE8021xSettings",
							SelectorSet: models.SelectorNoNamespace{
								XMLName: xml.Name{
									Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
									Local: "SelectorSet",
								},
								Selectors: []models.SelectorResponse{
									{
										XMLName: xml.Name{
											Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
											Local: "Selector",
										},
										Name: "InstanceID",
										Text: "Intel(r) AMT:IEEE 802.1x Settings TestWifi8021xTLS",
									},
								},
							},
						},
					},
				},
				{
					ElementInContext: models.AssociationReference{
						Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
						ReferenceParameters: models.ReferenceParametersNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
								Local: "ReferenceParameters",
							},
							ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_PublicKeyCertificate",
							SelectorSet: models.SelectorNoNamespace{
								XMLName: xml.Name{
									Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
									Local: "SelectorSet",
								},
								Selectors: []models.SelectorResponse{
									{
										XMLName: xml.Name{
											Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
											Local: "Selector",
										},
										Name: "InstanceID",
										Text: "Intel(r) AMT Certificate: Handle: 2",
									},
								},
							},
						},
					},
					ElementProvidingContext: models.AssociationReference{
						Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
						ReferenceParameters: models.ReferenceParametersNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
								Local: "ReferenceParameters",
							},
							ResourceURI: "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_IEEE8021xSettings",
							SelectorSet: models.SelectorNoNamespace{
								XMLName: xml.Name{
									Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
									Local: "SelectorSet",
								},
								Selectors: []models.SelectorResponse{
									{
										XMLName: xml.Name{
											Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
											Local: "Selector",
										},
										Name: "InstanceID",
										Text: "Intel(r) AMT:IEEE 802.1x Settings TestWifi8021xTLS2",
									},
								},
							},
						},
					},
				},
			},
			CredentialContextTLS: []credential.CredentialContext{
				{
					ElementInContext: models.AssociationReference{
						Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
						ReferenceParameters: models.ReferenceParametersNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
								Local: "ReferenceParameters",
							},
							ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_PublicKeyCertificate",
							SelectorSet: models.SelectorNoNamespace{
								XMLName: xml.Name{
									Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
									Local: "SelectorSet",
								},
								Selectors: []models.SelectorResponse{
									{
										XMLName: xml.Name{
											Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
											Local: "Selector",
										},
										Name: "InstanceID",
										Text: "Intel(r) AMT Certificate: Handle: 0",
									},
								},
							},
						},
					},
					ElementProvidingContext: models.AssociationReference{
						Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
						ReferenceParameters: models.ReferenceParametersNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
								Local: "ReferenceParameters",
							},
							ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_TLSProtocolEndpointCollection",
							SelectorSet: models.SelectorNoNamespace{
								XMLName: xml.Name{
									Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
									Local: "SelectorSet",
								},
								Selectors: []models.SelectorResponse{
									{
										XMLName: xml.Name{
											Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
											Local: "Selector",
										},
										Name: "ElementName",
										Text: "TLSProtocolEndpoint Instances Collection",
									},
								},
							},
						},
					},
				},
			},
			CredentialContext8021x: []credential.CredentialContext{
				{
					ElementInContext: models.AssociationReference{
						Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
						ReferenceParameters: models.ReferenceParametersNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
								Local: "ReferenceParameters",
							},
							ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_PublicKeyCertificate",
							SelectorSet: models.SelectorNoNamespace{
								XMLName: xml.Name{
									Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
									Local: "SelectorSet",
								},
								Selectors: []models.SelectorResponse{
									{
										XMLName: xml.Name{
											Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
											Local: "Selector",
										},
										Name: "InstanceID",
										Text: "Intel(r) AMT Certificate: Handle: 1",
									},
								},
							},
						},
					},
					ElementProvidingContext: models.AssociationReference{
						Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
						ReferenceParameters: models.ReferenceParametersNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
								Local: "ReferenceParameters",
							},
							ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_TLSProtocolEndpointCollection",
							SelectorSet: models.SelectorNoNamespace{
								XMLName: xml.Name{
									Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
									Local: "SelectorSet",
								},
								Selectors: []models.SelectorResponse{
									{
										XMLName: xml.Name{
											Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
											Local: "Selector",
										},
										Name: "InstanceID",
										Text: "Intel(r) AMT: 8021X Settings",
									},
								},
							},
						},
					},
				},
				{
					ElementInContext: models.AssociationReference{
						Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
						ReferenceParameters: models.ReferenceParametersNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
								Local: "ReferenceParameters",
							},
							ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_PublicKeyCertificate",
							SelectorSet: models.SelectorNoNamespace{
								XMLName: xml.Name{
									Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
									Local: "SelectorSet",
								},
								Selectors: []models.SelectorResponse{
									{
										XMLName: xml.Name{
											Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
											Local: "Selector",
										},
										Name: "InstanceID",
										Text: "Intel(r) AMT Certificate: Handle: 2",
									},
								},
							},
						},
					},
					ElementProvidingContext: models.AssociationReference{
						Address: "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
						ReferenceParameters: models.ReferenceParametersNoNamespace{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/ws/2004/08/addressing",
								Local: "ReferenceParameters",
							},
							ResourceURI: "http://intel.com/wbem/wscim/1/amt-schema/1/AMT_TLSProtocolEndpointCollection",
							SelectorSet: models.SelectorNoNamespace{
								XMLName: xml.Name{
									Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
									Local: "SelectorSet",
								},
								Selectors: []models.SelectorResponse{
									{
										XMLName: xml.Name{
											Space: "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd",
											Local: "Selector",
										},
										Name: "InstanceID",
										Text: "Intel(r) AMT: 8021X Settings",
									},
								},
							},
						},
					},
				},
			},
		},
		EndOfSequence: xml.Name{
			Space: "http://schemas.xmlsoap.org/ws/2004/09/enumeration",
			Local: "EndOfSequence",
		},
	},
}

func TestProcessCertificates(t *testing.T) {
	t.Parallel()

	securitySettings := SecuritySettings{
		Certificates: getResponse.PublicKeyCertificateResponse,
		Keys:         getResponse.PublicPrivateKeyPairResponse,
	}

	tests := []certificateTest{
		{
			name: "success",
			res: SecuritySettings{
				ProfileAssociation: []ProfileAssociation{
					{
						Type:      "Wireless",
						ProfileID: "TestWifi8021xTLS",
						RootCertificate: publickey.RefinedPublicKeyCertificateResponse{
							ElementName:            "Intel(r) AMT Certificate",
							InstanceID:             "Intel(r) AMT Certificate: Handle: 2",
							X509Certificate:        "TestCertRoot",
							TrustedRootCertificate: true,
							ReadOnlyCertificate:    false,
							PublicKeyHandle:        "",
							AssociatedProfiles:     nil,
						},
						ClientCertificate: publickey.RefinedPublicKeyCertificateResponse{
							ElementName:            "Intel(r) AMT Certificate",
							InstanceID:             "Intel(r) AMT Certificate: Handle: 3",
							X509Certificate:        "TestCert3",
							TrustedRootCertificate: false,
							ReadOnlyCertificate:    false,
							PublicKeyHandle:        "",
							AssociatedProfiles:     nil,
						},
						Key: publicprivate.RefinedPublicPrivateKeyPair{
							ElementName:       "Intel(r) AMT Key",
							InstanceID:        "Intel(r) AMT Key: Handle: 2",
							DERKey:            "Key2",
							CertificateHandle: "",
						},
					},
					{
						Type:      "Wireless",
						ProfileID: "TestWifi8021xTLS2",
						RootCertificate: publickey.RefinedPublicKeyCertificateResponse{
							ElementName:            "Intel(r) AMT Certificate",
							InstanceID:             "Intel(r) AMT Certificate: Handle: 2",
							X509Certificate:        "TestCertRoot",
							TrustedRootCertificate: true,
							ReadOnlyCertificate:    false,
							PublicKeyHandle:        "",
							AssociatedProfiles: []string{
								"TestWifi8021xTLS",
							},
						},
						ClientCertificate: publickey.RefinedPublicKeyCertificateResponse{},
						Key:               publicprivate.RefinedPublicPrivateKeyPair{},
					},
				},
				Certificates: publickey.RefinedPullResponse{
					KeyManagementItems: nil,
					PublicKeyCertificateItems: []publickey.RefinedPublicKeyCertificateResponse{
						{
							ElementName:            "Intel(r) AMT Certificate",
							InstanceID:             "Intel(r) AMT Certificate: Handle: 2",
							X509Certificate:        "TestRootCert",
							TrustedRootCertificate: true,
							ReadOnlyCertificate:    false,
							PublicKeyHandle:        "",
							AssociatedProfiles: []string{
								"exampleWifi8021xTLS",
								"exampleWifi8021xTLS2",
							},
						},
						{
							ElementName:            "Intel(r) AMT Certificate",
							InstanceID:             "Intel(r) AMT Certificate: Handle: 0",
							X509Certificate:        "TestCert0",
							TrustedRootCertificate: false,
							ReadOnlyCertificate:    false,
							PublicKeyHandle:        "Intel(r) AMT Key: Handle: 0",
							AssociatedProfiles:     nil,
						},
						{
							ElementName:            "Intel(r) AMT Certificate",
							InstanceID:             "Intel(r) AMT Certificate: Handle: 1",
							X509Certificate:        "TestCert1",
							TrustedRootCertificate: false,
							ReadOnlyCertificate:    false,
							PublicKeyHandle:        "",
							AssociatedProfiles:     nil,
						},
						{
							ElementName:            "Intel(r) AMT Certificate",
							InstanceID:             "Intel(r) AMT Certificate: Handle: 3",
							X509Certificate:        "TestCert3",
							TrustedRootCertificate: false,
							ReadOnlyCertificate:    false,
							PublicKeyHandle:        "",
							AssociatedProfiles: []string{
								"TestWifi8021xTLS",
							},
						},
					},
				},
				Keys: publicprivate.RefinedPullResponse{
					PublicPrivateKeyPairItems: []publicprivate.RefinedPublicPrivateKeyPair{
						{
							ElementName:       "Intel(r) AMT Key",
							InstanceID:        "Intel(r) AMT Key: Handle: 0",
							DERKey:            "Key0",
							CertificateHandle: "",
						},
						{
							ElementName:       "Intel(r) AMT Key",
							InstanceID:        "Intel(r) AMT Key: Handle: 1",
							DERKey:            "Key1",
							CertificateHandle: "",
						},
						{
							ElementName:       "Intel(r) AMT Key",
							InstanceID:        "Intel(r) AMT Key: Handle: 2",
							DERKey:            "Key2",
							CertificateHandle: "Intel(r) AMT Certificate: Handle: 3",
						},
					},
				},
			},
			err:         nil,
			profileType: "Wireless",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			processCertificates(getResponse.CIMCredentialContextResponse.Items.CredentialContext, getResponse, tc.profileType, &securitySettings)

			require.Equal(t, tc.res.ProfileAssociation[0].Type, securitySettings.ProfileAssociation[0].Type)
			require.Equal(t, tc.res.ProfileAssociation[0].ClientCertificate.(publickey.RefinedPublicKeyCertificateResponse).InstanceID, securitySettings.ProfileAssociation[0].ClientCertificate.(publickey.RefinedPublicKeyCertificateResponse).InstanceID)
			require.Equal(t, tc.res.ProfileAssociation[1].RootCertificate.(publickey.RefinedPublicKeyCertificateResponse).AssociatedProfiles, securitySettings.ProfileAssociation[1].RootCertificate.(publickey.RefinedPublicKeyCertificateResponse).AssociatedProfiles)
			require.Equal(t, tc.res.Certificates.(publickey.RefinedPullResponse).PublicKeyCertificateItems[3].AssociatedProfiles, securitySettings.Certificates.(publickey.RefinedPullResponse).PublicKeyCertificateItems[3].AssociatedProfiles)
			require.Equal(t, tc.res.Keys.(publicprivate.RefinedPullResponse).PublicPrivateKeyPairItems[2].CertificateHandle, securitySettings.Keys.(publicprivate.RefinedPullResponse).PublicPrivateKeyPairItems[2].CertificateHandle)
		})
	}
}
