package devices_test

import (
	"testing"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	dtov1 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	devices "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
)

func initRedirectionTest(t *testing.T) (*devices.Redirector, *MockRedirection, *MockRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockRepository(mockCtl)
	redirect := NewMockRedirection(mockCtl)
	u := &devices.Redirector{}

	return u, redirect, repo
}

type redTest struct {
	name    string
	redMock func(*MockRedirection)
	res     any
}

func TestSetupWsmanClient(t *testing.T) {
	t.Parallel()

	device := &dtov1.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}

	tests := []redTest{
		{
			name: "success",
			redMock: func(redirect *MockRedirection) {
				redirect.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(wsman.Messages{})
			},
			res: wsman.Messages{},
		},
		{
			name: "fail",
			redMock: func(redirect *MockRedirection) {
				redirect.EXPECT().
					SetupWsmanClient(gomock.Any(), true, true).
					Return(wsman.Messages{})
			},
			res: wsman.Messages{},
		},
	}

	for _, tc := range tests {
		tc := tc // Necessary for proper parallel execution
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			redirector, redirect, _ := initRedirectionTest(t)

			tc.redMock(redirect)

			res := redirector.SetupWsmanClient(*device, true, true)

			require.IsType(t, tc.res, res)
		})
	}
}

func TestNewRedirector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "success",
		},
	}

	for _, tc := range tests {
		tc := tc // Necessary for proper parallel execution
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Call the function under test
			redirector := devices.NewRedirector()

			// Assert that the returned redirector is not nil
			require.NotNil(t, redirector)
		})
	}
}
