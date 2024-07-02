package devices_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	devices "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

func initNetworkTest(t *testing.T) (*devices.UseCase, *MockManagement, *MockRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockRepository(mockCtl)
	management := NewMockManagement(mockCtl)
	log := logger.New("error")
	u := devices.New(repo, management, NewMockRedirection(mockCtl), log)

	return u, management, repo
}

func TestGetNetworkSettings(t *testing.T) {
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
					GetNetworkSettings().
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
			name:   "GetNetworkSettings fails",
			action: 0,
			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return()
				man.EXPECT().
					GetNetworkSettings().
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

			useCase, management, repo := initNetworkTest(t)

			if tc.manMock != nil {
				tc.manMock(management)
			}

			tc.repoMock(repo)

			res, err := useCase.GetNetworkSettings(context.Background(), device.GUID)

			require.Equal(t, tc.res, res)
			require.IsType(t, tc.err, err)
		})
	}
}
