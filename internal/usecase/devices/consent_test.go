package devices_test

import (
	"context"
	"encoding/xml"
	"testing"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/optin"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	"github.com/open-amt-cloud-toolkit/console/internal/mocks"
	devices "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

func initConsentTest(t *testing.T) (*devices.UseCase, *mocks.MockWSMAN, *mocks.MockManagement, *mocks.MockDeviceManagementRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)

	defer mockCtl.Finish()

	repo := mocks.NewMockDeviceManagementRepository(mockCtl)

	wsmanMock := mocks.NewMockWSMAN(mockCtl)
	wsmanMock.EXPECT().Worker().Return().AnyTimes()

	management := mocks.NewMockManagement(mockCtl)

	log := logger.New("error")

	u := devices.New(repo, wsmanMock, mocks.NewMockRedirection(mockCtl), log, mocks.MockCrypto{})

	return u, wsmanMock, management, repo
}

func TestCancelUserConsent(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID: "device-guid-123",

		TenantID: "tenant-id-456",
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
					CancelUserConsentRequest().
					Return(gomock.Any(), nil)
			},
			repoMock: func(repo *mocks.MockDeviceManagementRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: gomock.Any(),
			err: nil,
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
			res: nil,
			err: devices.ErrGeneral,
		},

		{
			name:   "CancelUserConsentRequest fails",
			action: 0,
			manMock: func(man *mocks.MockWSMAN, man2 *mocks.MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					CancelUserConsentRequest().
					Return(nil, ErrGeneral)
			},
			repoMock: func(repo *mocks.MockDeviceManagementRepository) {
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

			useCase, wsmanMock, management, repo := initConsentTest(t)

			tc.manMock(wsmanMock, management)

			tc.repoMock(repo)

			res, err := useCase.CancelUserConsent(context.Background(), device.GUID)

			require.Equal(t, tc.res, res)

			require.IsType(t, tc.err, err)
		})
	}
}

func TestGetUserConsentCode(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}

	code := optin.StartOptIn_OUTPUT{
		XMLName: xml.Name{
			Local: "StartOptIn_OUTPUT",
		},
		ReturnValue: 10,
	}

	response := map[string]interface{}{
		"Body": code,
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
					GetUserConsentCode().
					Return(code, nil)
			},
			repoMock: func(repo *mocks.MockDeviceManagementRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: response,
			err: nil,
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
			res: map[string]interface{}(nil),
			err: devices.ErrGeneral,
		},

		{
			name:   "GetUserConsentCode fails",
			action: 0,
			manMock: func(man *mocks.MockWSMAN, man2 *mocks.MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetUserConsentCode().
					Return(code, ErrGeneral)
			},
			repoMock: func(repo *mocks.MockDeviceManagementRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: map[string]interface{}(nil),
			err: ErrGeneral,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, wsmanMock, management, repo := initConsentTest(t)

			tc.manMock(wsmanMock, management)

			tc.repoMock(repo)

			res, err := useCase.GetUserConsentCode(context.Background(), device.GUID)

			require.Equal(t, tc.res, res)

			require.IsType(t, tc.err, err)
		})
	}
}

func TestSendConsentCode(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}

	consent := dto.UserConsent{
		ConsentCode: "123456",
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
					SendConsentCode(123456).
					Return(optin.SendOptInCode_OUTPUT{}, nil)
			},
			repoMock: func(repo *mocks.MockDeviceManagementRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: optin.SendOptInCode_OUTPUT{},
			err: nil,
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
			res: nil,
			err: devices.ErrGeneral,
		},
		{
			name:   "SendConsentCode fails",
			action: 0,
			manMock: func(man *mocks.MockWSMAN, man2 *mocks.MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					SendConsentCode(123456).
					Return(optin.SendOptInCode_OUTPUT{}, ErrGeneral)
			},
			repoMock: func(repo *mocks.MockDeviceManagementRepository) {
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

			useCase, wsmanMock, management, repo := initConsentTest(t)

			tc.manMock(wsmanMock, management)

			tc.repoMock(repo)

			res, err := useCase.SendConsentCode(context.Background(), consent, device.GUID)

			require.Equal(t, tc.res, res)

			require.IsType(t, tc.err, err)
		})
	}
}
