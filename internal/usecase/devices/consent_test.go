package devices_test

import (
	"context"
	"encoding/xml"
	"testing"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/optin"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	devices "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

func initConsentTest(t *testing.T) (*devices.UseCase, *MockManagement, *MockRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)

	defer mockCtl.Finish()

	repo := NewMockRepository(mockCtl)

	management := NewMockManagement(mockCtl)

	amt := NewMockAMTExplorer(mockCtl)

	log := logger.New("error")

	u := devices.New(repo, management, NewMockRedirection(mockCtl), amt, log)

	return u, management, repo
}

func TestCancelUserConsent(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID: "device-guid-123",

		TenantID: "tenant-id-456",
	}

	tests := []test{
		{
			name: "success",

			action: 0,

			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return()

				man.EXPECT().
					CancelUserConsentRequest().
					Return(gomock.Any(), nil)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: gomock.Any(),

			err: nil,
		},

		{
			name: "GetById fails",

			action: 0,

			manMock: func(_ *MockManagement) {},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},

			res: nil,

			err: devices.ErrDatabase,
		},

		{
			name: "CancelUserConsentRequest fails",

			action: 0,

			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return()

				man.EXPECT().
					CancelUserConsentRequest().
					Return(nil, ErrGeneral)
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

			useCase, management, repo := initConsentTest(t)

			tc.manMock(management)

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
		GUID: "device-guid-123",

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
			name: "success",

			action: 0,

			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return()

				man.EXPECT().
					GetUserConsentCode().
					Return(code, nil)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: response,

			err: nil,
		},

		{
			name: "GetById fails",

			action: 0,

			manMock: func(_ *MockManagement) {},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},

			res: map[string]interface{}(nil),

			err: devices.ErrDatabase,
		},

		{
			name: "GetUserConsentCode fails",

			action: 0,

			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return()

				man.EXPECT().
					GetUserConsentCode().
					Return(code, ErrGeneral)
			},

			repoMock: func(repo *MockRepository) {
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

			useCase, management, repo := initConsentTest(t)

			tc.manMock(management)

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
		GUID: "device-guid-123",

		TenantID: "tenant-id-456",
	}

	consent := dto.UserConsent{
		ConsentCode: "123456",
	}

	tests := []test{
		{
			name: "success",

			action: 0,

			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return()

				man.EXPECT().
					SendConsentCode(123456).
					Return(optin.SendOptInCode_OUTPUT{}, nil)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: optin.SendOptInCode_OUTPUT{},

			err: nil,
		},

		{
			name: "GetById fails",

			action: 0,

			manMock: func(_ *MockManagement) {},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},

			res: nil,

			err: devices.ErrDatabase,
		},

		{
			name: "SendConsentCode fails",

			action: 0,

			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return()

				man.EXPECT().
					SendConsentCode(123456).
					Return(optin.SendOptInCode_OUTPUT{}, ErrGeneral)
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

			useCase, management, repo := initConsentTest(t)

			tc.manMock(management)

			tc.repoMock(repo)

			res, err := useCase.SendConsentCode(context.Background(), consent, device.GUID)

			require.Equal(t, tc.res, res)

			require.IsType(t, tc.err, err)
		})
	}
}
