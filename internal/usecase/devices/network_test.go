package devices_test

import (
	"context"
	"testing"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/ethernetport"
	cimieee8021x "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/ieee8021x"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/wifi"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/ieee8021x"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	dtov1 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	devices "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

func initNetworkTest(t *testing.T) (*devices.UseCase, *MockWSMAN, *MockManagement, *MockRepository) {
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
			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetNetworkSettings().
					Return(wsman.NetworkResults{
						EthernetPortSettingsResult: []ethernetport.SettingsResponse{
							{
								LinkPolicy: []ethernetport.LinkPolicy{1, 2},
							}, {
								LinkPolicy: []ethernetport.LinkPolicy{1, 2},
							},
						},
						IPSIEEE8021xSettingsResult: ieee8021x.IEEE8021xSettingsResponse{},
						WiFiSettingsResult:         []wifi.WiFiEndpointSettingsResponse{{}},
						CIMIEEE8021xSettingsResult: cimieee8021x.PullResponse{
							IEEE8021xSettingsItems: []cimieee8021x.IEEE8021xSettingsResponse{{}},
						},
					}, nil)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: dtov1.NetworkSettings{
				Wired: dtov1.WiredNetworkInfo{
					IEEE8021x: dtov1.IEEE8021x{},
					NetworkInfo: dtov1.NetworkInfo{
						LinkPolicy: []int{1, 2},
					},
				},
				Wireless: dtov1.WirelessNetworkInfo{
					WiFiNetworks:      []dtov1.WiFiNetwork{{}},
					IEEE8021xSettings: []dtov1.IEEE8021xSettings{{}},
					NetworkInfo: dtov1.NetworkInfo{
						LinkPolicy: []int{1, 2},
					},
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
			res: dtov1.NetworkSettings{},
			err: devices.ErrDatabase,
		},
		{
			name:   "GetNetworkSettings fails",
			action: 0,
			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetNetworkSettings().
					Return(wsman.NetworkResults{}, ErrGeneral)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: dtov1.NetworkSettings{},
			err: ErrGeneral,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, wsmanMock, management, repo := initNetworkTest(t)

			if tc.manMock != nil {
				tc.manMock(wsmanMock, management)
			}

			tc.repoMock(repo)

			res, err := useCase.GetNetworkSettings(context.Background(), device.GUID)

			require.Equal(t, tc.res, res)
			require.IsType(t, tc.err, err)
		})
	}
}
