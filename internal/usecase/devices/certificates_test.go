package devices_test

import (
	"context"
	"encoding/xml"
	"testing"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/credential"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/models"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	dtov1 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	devices "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	wsman "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

func initCertificateTest(t *testing.T) (*devices.UseCase, *MockWSMAN, *MockManagement, *MockRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockRepository(mockCtl)
	wsmanMock := NewMockWSMAN(mockCtl)
	wsmanMock.EXPECT().Worker().Return().AnyTimes()

	management := NewMockManagement(mockCtl)
	log := logger.New("error")
	u := devices.New(repo, wsmanMock, NewMockRedirection(mockCtl), log)

	return u, wsmanMock, management, repo
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
			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetCertificates().
					Return(wsman.Certificates{}, nil)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: dtov1.SecuritySettings{
				ProfileAssociation: []dtov1.ProfileAssociation(nil),
				CertificateResponse: dtov1.CertificatePullResponse{
					KeyManagementItems: []dtov1.RefinedKeyManagementResponse{},
					Certificates:       []dtov1.RefinedCertificate{},
				},
				KeyResponse: dtov1.KeyPullResponse{
					Keys: []dtov1.Key{},
				},
			},
			err: nil,
		},
		{
			name:   "success with CIMCredentialContext",
			action: 0,
			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
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
			res: dtov1.SecuritySettings{
				ProfileAssociation: []dtov1.ProfileAssociation{
					{
						Type:              "TLS",
						ProfileID:         "TLSProtocolEndpoint Instances Collection",
						RootCertificate:   nil,
						ClientCertificate: nil,
						Key:               nil,
					},
				},
				CertificateResponse: dtov1.CertificatePullResponse{
					KeyManagementItems: []dtov1.RefinedKeyManagementResponse{},
					Certificates:       []dtov1.RefinedCertificate{},
				},
				KeyResponse: dtov1.KeyPullResponse{
					Keys: []dtov1.Key{},
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
			res: dtov1.SecuritySettings{},
			err: devices.ErrDatabase,
		},
		{
			name:   "GetCertificates fails",
			action: 0,
			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetCertificates().
					Return(wsman.Certificates{}, ErrGeneral)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: dtov1.SecuritySettings{},
			err: ErrGeneral,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, wsmanMock, management, repo := initCertificateTest(t)

			if tc.manMock != nil {
				tc.manMock(wsmanMock, management)
			}

			tc.repoMock(repo)

			res, err := useCase.GetCertificates(context.Background(), device.GUID)

			require.Equal(t, tc.res, res)
			require.IsType(t, tc.err, err)
		})
	}
}
