package devices_test

import (
	"context"
	"testing"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/amterror"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/redirection"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/kvm"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/optin"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	dtov2 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v2"
	"github.com/open-amt-cloud-toolkit/console/internal/mocks"
	devices "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
)

const DestinationUnreachable string = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><a:Envelope xmlns:g=\"http://schemas.dmtf.org/wbem/wsman/1/cimbinding.xsd\" xmlns:f=\"http://schemas.xmlsoap.org/ws/2004/08/eventing\" xmlns:e=\"http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd\" xmlns:d=\"http://schemas.xmlsoap.org/ws/2004/09/transfer\" xmlns:c=\"http://schemas.xmlsoap.org/ws/2004/09/enumeration\" xmlns:b=\"http://schemas.xmlsoap.org/ws/2004/08/addressing\" xmlns:a=\"http://www.w3.org/2003/05/soap-envelope\" xmlns:h=\"http://schemas.xmlsoap.org/ws/2005/02/trust\" xmlns:i=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\"><a:Header><b:To>http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</b:To><b:RelatesTo>0</b:RelatesTo><b:Action a:mustUnderstand=\"true\">http://schemas.xmlsoap.org/ws/2004/08/addressing/fault</b:Action><b:MessageID>uuid:00000000-8086-8086-8086-000000000061</b:MessageID></a:Header><a:Body><a:Fault><a:Code><a:Value>a:Sender</a:Value><a:Subcode><a:Value>b:DestinationUnreachable</a:Value></a:Subcode></a:Code><a:Reason><a:Text xml:lang=\"en-US\">No route can be determined to reach the destination role defined by the WSAddressing To.</a:Text></a:Reason><a:Detail></a:Detail></a:Fault></a:Body></a:Envelope>"

func TestGetFeatures(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}

	featureSet := dto.Features{
		UserConsent:  "kvm",
		EnableSOL:    true,
		EnableIDER:   true,
		EnableKVM:    true,
		Redirection:  true,
		KVMAvailable: true,
		OptInState:   1,
	}

	featureSetNoKVM := dto.Features{
		UserConsent:  "kvm",
		EnableSOL:    true,
		EnableIDER:   true,
		EnableKVM:    false,
		Redirection:  true,
		KVMAvailable: false,
		OptInState:   1,
	}

	featureSetV2 := dtov2.Features{
		UserConsent:  "kvm",
		EnableSOL:    true,
		EnableIDER:   true,
		EnableKVM:    true,
		Redirection:  true,
		KVMAvailable: true,
		OptInState:   1,
	}

	featureSetV2NoKVM := dtov2.Features{
		UserConsent:  "kvm",
		EnableSOL:    true,
		EnableIDER:   true,
		EnableKVM:    false,
		Redirection:  true,
		KVMAvailable: false,
		OptInState:   1,
	}

	tests := []test{
		{
			name:   "success",
			action: 0,
			manMock: func(man *mocks.MockWSMAN, man2 *mocks.MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetAMTRedirectionService().
					Return(redirection.Response{
						Body: redirection.Body{
							GetAndPutResponse: redirection.RedirectionResponse{
								EnabledState:    32771,
								ListenerEnabled: true,
							},
						},
					}, nil)
				man2.EXPECT().
					GetIPSOptInService().
					Return(optin.Response{
						Body: optin.Body{
							GetAndPutResponse: optin.OptInServiceResponse{
								OptInRequired: 1,
								OptInState:    1,
							},
						},
					}, nil)
				man2.EXPECT().
					GetKVMRedirection().
					Return(kvm.Response{
						Body: kvm.Body{
							GetResponse: kvm.KVMRedirectionSAP{
								EnabledState: kvm.EnabledState(redirection.Enabled),
							},
						},
					}, nil)
			},
			repoMock: func(repo *mocks.MockDeviceManagementRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res:   featureSet,
			resV2: featureSetV2,
			err:   nil,
		},
		{
			name:    "GetById fails",
			action:  0,
			manMock: func(_ *mocks.MockWSMAN, _ *mocks.MockManagement) {},
			repoMock: func(repo *mocks.MockDeviceManagementRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},
			res:   dto.Features{},
			resV2: dtov2.Features{},
			err:   devices.ErrGeneral,
		},
		{
			name:   "GetFeatures fails",
			action: 0,
			manMock: func(man *mocks.MockWSMAN, man2 *mocks.MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetAMTRedirectionService().
					Return(redirection.Response{}, ErrGeneral)
				man2.EXPECT().
					GetIPSOptInService().
					Return(optin.Response{}, ErrGeneral)
				man2.EXPECT().
					GetKVMRedirection().
					Return(kvm.Response{}, ErrGeneral)
			},
			repoMock: func(repo *mocks.MockDeviceManagementRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res:   dto.Features{},
			resV2: dtov2.Features{},
			err:   ErrGeneral,
		},
		{
			name:   "GetFeatures on ISM",
			action: 0,
			manMock: func(man *mocks.MockWSMAN, man2 *mocks.MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetAMTRedirectionService().
					Return(redirection.Response{
						Body: redirection.Body{
							GetAndPutResponse: redirection.RedirectionResponse{
								EnabledState:    32771,
								ListenerEnabled: true,
							},
						},
					}, nil)
				man2.EXPECT().
					GetIPSOptInService().
					Return(optin.Response{
						Body: optin.Body{
							GetAndPutResponse: optin.OptInServiceResponse{
								OptInState:    1,
								OptInRequired: 1,
							},
						},
					}, nil)
				man2.EXPECT().
					GetKVMRedirection().
					Return(kvm.Response{}, amterror.DecodeAMTErrorString(DestinationUnreachable))
			},
			repoMock: func(repo *mocks.MockDeviceManagementRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res:   featureSetNoKVM,
			resV2: featureSetV2NoKVM,
			err:   nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, wsmanMock, management, repo := initInfoTest(t)

			tc.manMock(wsmanMock, management)

			tc.repoMock(repo)

			v1, v2, err := useCase.GetFeatures(context.Background(), device.GUID)

			require.Equal(t, tc.res, v1)

			require.Equal(t, tc.resV2, v2)

			require.IsType(t, tc.err, err)
		})
	}
}

func TestSetFeatures(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}

	featureSet := dto.Features{
		UserConsent: "kvm",
		EnableSOL:   true,
		EnableIDER:  true,
		EnableKVM:   true,
		Redirection: true,
	}

	featureSetV2 := dtov2.Features{
		UserConsent:  "kvm",
		EnableSOL:    true,
		EnableIDER:   true,
		EnableKVM:    true,
		Redirection:  true,
		KVMAvailable: true,
	}

	tests := []test{
		{
			name:   "success",
			action: 0,
			manMock: func(man *mocks.MockWSMAN, man2 *mocks.MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					RequestAMTRedirectionServiceStateChange(featureSet.EnableSOL, featureSet.EnableIDER).
					Return(redirection.EnableIDERAndSOL, 1, nil)
				man2.EXPECT().
					SetKVMRedirection(true).
					Return(1, nil)
				man2.EXPECT().
					GetAMTRedirectionService().
					Return(redirection.Response{
						Body: redirection.Body{
							GetAndPutResponse: redirection.RedirectionResponse{
								EnabledState:    32771,
								ListenerEnabled: true,
							},
						},
					}, nil)
				man2.EXPECT().
					SetAMTRedirectionService(redirection.RedirectionRequest{
						EnabledState:    redirection.EnabledState(redirection.EnableIDERAndSOL),
						ListenerEnabled: true,
					}).
					Return(redirection.Response{
						Body: redirection.Body{
							GetAndPutResponse: redirection.RedirectionResponse{
								EnabledState:    32771,
								ListenerEnabled: true,
							},
						},
					}, nil)
				man2.EXPECT().
					GetIPSOptInService().
					Return(optin.Response{
						Body: optin.Body{
							GetAndPutResponse: optin.OptInServiceResponse{
								OptInRequired: 1,
								OptInState:    0,
							},
						},
					}, nil)
				man2.EXPECT().
					SetIPSOptInService(optin.OptInServiceRequest{
						OptInRequired: 1,
						OptInState:    0,
					}).
					Return(nil)
			},
			repoMock: func(repo *mocks.MockDeviceManagementRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res:   featureSet,
			resV2: featureSetV2,
			err:   nil,
		},
		{
			name:    "GetById fails",
			action:  0,
			manMock: func(_ *mocks.MockWSMAN, _ *mocks.MockManagement) {},
			repoMock: func(repo *mocks.MockDeviceManagementRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},
			res:   dto.Features{},
			resV2: dtov2.Features{},
			err:   devices.ErrGeneral,
		},
		{
			name:   "GetFeatures fails",
			action: 0,
			manMock: func(man *mocks.MockWSMAN, man2 *mocks.MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					RequestAMTRedirectionServiceStateChange(featureSet.EnableSOL, featureSet.EnableIDER).
					Return(redirection.RequestedState(0), 0, ErrGeneral)
				man2.EXPECT().
					SetKVMRedirection(true).
					Return(0, ErrGeneral)
				man2.EXPECT().
					GetAMTRedirectionService().
					Return(redirection.Response{}, ErrGeneral)
				man2.EXPECT().
					GetIPSOptInService().
					Return(optin.Response{}, ErrGeneral)
				man2.EXPECT().
					SetIPSOptInService(optin.OptInServiceRequest{}).
					Return(ErrGeneral)
			},
			repoMock: func(repo *mocks.MockDeviceManagementRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res:   dto.Features{},
			resV2: dtov2.Features{},
			err:   ErrGeneral,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, wsmanMock, management, repo := initInfoTest(t)

			tc.manMock(wsmanMock, management)

			tc.repoMock(repo)

			v1, v2, err := useCase.SetFeatures(context.Background(), device.GUID, featureSet)

			require.Equal(t, tc.res, v1)

			require.Equal(t, tc.resV2, v2)

			require.IsType(t, tc.err, err)
		})
	}
}
