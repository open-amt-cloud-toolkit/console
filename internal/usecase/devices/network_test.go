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
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
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
								LinkPolicy:             []ethernetport.LinkPolicy{14, 16},
								PhysicalConnectionType: 0,
								PhysicalNicMedium:      0,
							}, {
								LinkPolicy:              []ethernetport.LinkPolicy{14, 16},
								LinkPreference:          1,
								LinkControl:             1,
								WLANLinkProtectionLevel: 1,
								PhysicalConnectionType:  3,
								PhysicalNicMedium:       1,
							},
						},
						IPSIEEE8021xSettingsResult: ieee8021x.IEEE8021xSettingsResponse{},
						WiFiSettingsResult: []wifi.WiFiEndpointSettingsResponse{{
							SSID:                 "test-ssid",
							AuthenticationMethod: 6,
							EncryptionMethod:     3,
							Priority:             1,
							BSSType:              2,
						}},
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
			res: dto.NetworkSettings{
				Wired: dto.WiredNetworkInfo{
					IEEE8021x: dto.IEEE8021x{},
					NetworkInfo: dto.NetworkInfo{
						LinkPolicy:             []string{"Sx AC", "S0 DC"},
						PhysicalConnectionType: "Integrated LAN NIC",
						PhysicalNicMedium:      "SMBUS",
					},
				},
				Wireless: dto.WirelessNetworkInfo{
					WiFiNetworks: []dto.WiFiNetwork{{
						SSID:                 "test-ssid",
						AuthenticationMethod: "WPA2PSK",
						EncryptionMethod:     "TKIP",
						Priority:             1,
						BSSType:              "Independent",
					}},
					IEEE8021xSettings: []dto.IEEE8021xSettings{{}},
					NetworkInfo: dto.NetworkInfo{
						LinkPolicy:              []string{"Sx AC", "S0 DC"},
						LinkPreference:          "Management Engine",
						LinkControl:             "Management Engine",
						WLANLinkProtectionLevel: "None",
						PhysicalConnectionType:  "Wireless LAN",
						PhysicalNicMedium:       "PCIe",
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
			res: dto.NetworkSettings{},
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
			res: dto.NetworkSettings{},
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
