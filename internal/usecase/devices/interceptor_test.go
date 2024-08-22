//nolint:gci // ignore import order
package devices_test

import (
	"context"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	devices "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman"
)

func TestRedirect(t *testing.T) {
	t.Parallel()

	mockConn := &websocket.Conn{}
	guid := "device-guid-123"
	mode := "default"

	tests := []struct {
		name        string
		setup       func(*MockRedirection, *MockRepository, *MockWSMAN, *sync.WaitGroup)
		expectedErr error
	}{
		{
			name: "successful redirection",
			setup: func(mockRedir *MockRedirection, mockRepo *MockRepository, mockWSMAN *MockWSMAN, wg *sync.WaitGroup) {
				mockWSMAN.EXPECT().Worker().Do(func() {
					defer wg.Done()
				}).Times(1)
				mockRepo.EXPECT().GetByID(gomock.Any(), guid, "").Return(&entity.Device{
					GUID:     guid,
					Username: "user",
					Password: "pass",
				}, nil)
				mockRedir.EXPECT().SetupWsmanClient(gomock.Any(), true, true).Return(wsman.Messages{})
				mockRedir.EXPECT().RedirectConnect(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "GetByID fail redirection",
			setup: func(_ *MockRedirection, mockRepo *MockRepository, mockWSMAN *MockWSMAN, wg *sync.WaitGroup) {
				mockWSMAN.EXPECT().Worker().Do(func() {
					defer wg.Done()
				}).Times(1)
				mockRepo.EXPECT().GetByID(gomock.Any(), guid, "").Return(nil, ErrGeneral)
			},
			expectedErr: ErrGeneral,
		},
		{
			name: "RedirectConnect fail redirection",
			setup: func(mockRedir *MockRedirection, mockRepo *MockRepository, mockWSMAN *MockWSMAN, wg *sync.WaitGroup) {
				mockWSMAN.EXPECT().Worker().Do(func() {
					defer wg.Done()
				}).Times(1)
				mockRepo.EXPECT().GetByID(gomock.Any(), guid, "").Return(&entity.Device{
					GUID:     guid,
					Username: "user",
					Password: "pass",
				}, nil)
				mockRedir.EXPECT().SetupWsmanClient(gomock.Any(), true, true).Return(wsman.Messages{})
				mockRedir.EXPECT().RedirectConnect(gomock.Any(), gomock.Any()).Return(ErrGeneral)
			},
			expectedErr: ErrGeneral,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRedirection := NewMockRedirection(ctrl)
			mockRepo := NewMockRepository(ctrl)
			mockWSMAN := NewMockWSMAN(ctrl)

			var wg sync.WaitGroup

			wg.Add(1)

			tc.setup(mockRedirection, mockRepo, mockWSMAN, &wg)

			uc := devices.New(mockRepo, mockWSMAN, mockRedirection, logger.New("test"))

			wg.Wait()

			err := uc.Redirect(context.Background(), mockConn, guid, mode)

			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
