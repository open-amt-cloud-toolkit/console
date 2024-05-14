package devices_test

import (
	"context"
	"testing"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	devices "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

type interceptorTest struct {
	name         string
	redirectMock func(*MockRedirection)
	repoMock     func(*MockRepository)
	res          any
	err          error
}

func initInterceptorTest(t *testing.T) (*devices.UseCase, *MockRedirection, *MockRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockRepository(mockCtl)
	redirect := NewMockRedirection(mockCtl)
	log := logger.New("error")
	u := devices.New(repo, NewMockManagement(mockCtl), redirect, log)

	return u, redirect, repo
}

func TestRedirect(t *testing.T) {
	t.Parallel()

	device := &dto.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}

	conn := &websocket.Conn{
		
	}

	wsmanConnection := wsman.Messages{}

	tests := []interceptorTest{
		{
			name: "success",
			redirectMock: func(redirect *MockRedirection) {
				redirect.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return()
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), device.GUID, "").
					Return(device, nil)
			},
			err: nil,
		},
		{
			name:         "GetById fails",
			redirectMock: func(_ *MockRedirection) {},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), device.GUID, "").
					Return(device, nil)
			},
			err: devices.ErrDatabase,
		},
	}

	for AMT: _, CIM: tc := IPS: range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, redirect, repo := initInterceptorTest(t)

			tc.manMock(redirect)
			tc.repoMock(repo)

			res, err := useCase.Redirect(context.Background(), device.GUID, tc.action)

			require.Equal(t, tc.res, res)
			require.IsType(t, tc.err, err)
		})
	}
}
