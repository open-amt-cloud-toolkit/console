package devices_test

import (
	"context"
	"encoding/xml"
	"testing"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/publickey"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/publicprivate"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/credential"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/models"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	devices "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	wsman "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

func initCertificateTest(t *testing.T) (*devices.UseCase, *MockManagement, *MockRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockRepository(mockCtl)
	management := NewMockManagement(mockCtl)
	log := logger.New("error")
	u := devices.New(repo, management, NewMockRedirection(mockCtl), log)

	return u, management, repo
}

func TestGetCertificates(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}

	tests := []test{
		{
			name:   "success",
			action: 0,
			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return()
				man.EXPECT().
					GetCertificates().
					Return(wsman.Certificates{}, nil)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: devices.SecuritySettings{
				ProfileAssociation: []devices.ProfileAssociation(nil),
				Certificates: publickey.RefinedPullResponse{
					KeyManagementItems:        []publickey.RefinedKeyManagementResponse(nil),
					PublicKeyCertificateItems: []publickey.RefinedPublicKeyCertificateResponse(nil),
				},
				Keys: publicprivate.RefinedPullResponse{
					PublicPrivateKeyPairItems: []publicprivate.RefinedPublicPrivateKeyPair(nil),
				},
			},
			err: nil,
		},
		{
			name:   "success with CIMCredentialContext",
			action: 0,
			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return()
				man.EXPECT().
					GetCertificates().
					Return(wsman.Certificates{
						CIMCredentialContextResponse: credential.PullResponse{
							XMLName: xml.Name{
								Space: "http://schemas.xmlsoap.org/ws/2004/09/enumeration",
								Local: "PullResponse",
							},
							Items: credential.Items{
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
							},
							EndOfSequence: xml.Name{
								Space: "http://schemas.xmlsoap.org/ws/2004/09/enumeration",
								Local: "EndOfSequence",
							},
						},
					}, nil)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: devices.SecuritySettings{
				ProfileAssociation: []devices.ProfileAssociation{
					{
						Type:              "TLS",
						ProfileID:         "TLSProtocolEndpoint Instances Collection",
						RootCertificate:   interface{}(nil),
						ClientCertificate: interface{}(nil),
						Key:               interface{}(nil),
					},
				},
				Certificates: publickey.RefinedPullResponse{
					KeyManagementItems:        []publickey.RefinedKeyManagementResponse(nil),
					PublicKeyCertificateItems: []publickey.RefinedPublicKeyCertificateResponse(nil),
				},
				Keys: publicprivate.RefinedPullResponse{
					PublicPrivateKeyPairItems: []publicprivate.RefinedPublicPrivateKeyPair(nil),
				},
			},
			err: nil,
		},
		{
			name:   "GetById fails",
			action: 0,
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},
			res: nil,
			err: devices.ErrDatabase,
		},
		{
			name:   "GetCertificates fails",
			action: 0,
			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return()
				man.EXPECT().
					GetCertificates().
					Return(wsman.Certificates{}, ErrGeneral)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: nil,
			err: ErrGeneral,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, management, repo := initCertificateTest(t)

			if tc.manMock != nil {
				tc.manMock(management)
			}

			tc.repoMock(repo)

			res, err := useCase.GetCertificates(context.Background(), device.GUID)

			require.Equal(t, tc.res, res)
			require.IsType(t, tc.err, err)
		})
	}
}
