package devices_test

import (
	"context"
	"errors"
	"testing"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/boot"
	cimBoot "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/boot"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/power"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/service"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/software"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	devices "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

var ErrGeneral = errors.New("general error")

type test struct {
	name     string
	manMock  func(*MockWSMAN, *MockManagement)
	repoMock func(*MockRepository)
	res      any
	err      error

	action int
}

func initPowerTest(t *testing.T) (*devices.UseCase, *MockWSMAN, *MockManagement, *MockRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockRepository(mockCtl)
	wsmanMock := NewMockWSMAN(mockCtl)
	wsmanMock.EXPECT().Worker().Return().AnyTimes()

	managementMock := NewMockManagement(mockCtl)
	log := logger.New("error")
	u := devices.New(repo, wsmanMock, NewMockRedirection(mockCtl), log)

	return u, wsmanMock, managementMock, repo
}

func TestSendPowerAction(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}

	powerActionRes := power.PowerActionResponse{
		ReturnValue: 0,
	}

	tests := []test{
		{
			name:   "success",
			action: 0,
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(hmm)
				hmm.EXPECT().
					SendPowerAction(0).
					Return(powerActionRes, nil)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: powerActionRes,
			err: nil,
		},
		{
			name:    "GetById fails",
			action:  0,
			manMock: func(_ *MockWSMAN, _ *MockManagement) {},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},
			res: power.PowerActionResponse{},
			err: devices.ErrDatabase,
		},
		{
			name:   "SendPowerAction fails",
			action: 0,
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(hmm)
				hmm.EXPECT().
					SendPowerAction(0).
					Return(power.PowerActionResponse{}, ErrGeneral)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: power.PowerActionResponse{},
			err: ErrGeneral,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, wsmanMock, management, repo := initPowerTest(t)

			tc.manMock(wsmanMock, management)
			tc.repoMock(repo)

			res, err := useCase.SendPowerAction(context.Background(), device.GUID, tc.action)

			require.Equal(t, tc.res, res)
			require.IsType(t, tc.err, err)
		})
	}
}

func TestGetPowerState(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}

	tests := []test{
		{
			name: "success",
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(hmm)
				hmm.EXPECT().
					GetPowerState().
					Return([]service.CIM_AssociatedPowerManagementService{{PowerState: 0}}, nil)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: map[string]interface{}{
				"powerstate": service.PowerState(0),
			},
			err: nil,
		},
		{
			name:    "GetById fails",
			manMock: func(_ *MockWSMAN, _ *MockManagement) {},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},
			res: map[string]interface{}(nil),
			err: devices.ErrDatabase,
		},
		{
			name: "GetPowerState fails",
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(hmm)
				hmm.EXPECT().
					GetPowerState().
					Return([]service.CIM_AssociatedPowerManagementService{}, ErrGeneral)
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

			useCase, wsmanMock, management, repo := initPowerTest(t)

			tc.manMock(wsmanMock, management)
			tc.repoMock(repo)

			res, err := useCase.GetPowerState(context.Background(), device.GUID)

			require.Equal(t, tc.res, res)
			require.IsType(t, tc.err, err)
		})
	}
}

func TestGetPowerCapabilities(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}

	tests := []test{
		{
			name: "success",
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(hmm)
				hmm.EXPECT().
					GetAMTVersion().
					Return([]software.SoftwareIdentity{}, nil)
				hmm.EXPECT().
					GetPowerCapabilities().
					Return(boot.BootCapabilitiesResponse{}, nil)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: map[string]interface{}{
				"Power cycle":              5,
				"Power down":               8,
				"Power on to IDE-R CDROM":  203,
				"Power on to IDE-R Floppy": 201,
				"Power on to PXE":          401,
				"Power up":                 2,
				"Reset":                    10,
				"Reset to IDE-R CDROM":     202,
				"Reset to IDE-R Floppy":    200,
				"Reset to PXE":             400,
			},
			err: nil,
		},
		{
			name:    "GetById fails",
			manMock: func(_ *MockWSMAN, _ *MockManagement) {},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},
			res: map[string]interface{}(nil),
			err: devices.ErrDatabase,
		},
		{
			name: "GetPowerCapabilities fails",
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(hmm)
				hmm.EXPECT().
					GetPowerCapabilities().
					Return(boot.BootCapabilitiesResponse{}, ErrGeneral)
				hmm.EXPECT().
					GetAMTVersion().
					Return(nil, nil)
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

			useCase, wsmanMock, management, repo := initPowerTest(t)
			tc.manMock(wsmanMock, management)
			tc.repoMock(repo)

			res, err := useCase.GetPowerCapabilities(context.Background(), device.GUID)

			require.Equal(t, tc.res, res)
			require.IsType(t, tc.err, err)
		})
	}
}

func TestSetBootOptions(t *testing.T) {
	t.Parallel()

	bootResponse := boot.BootSettingDataResponse{
		BIOSLastStatus:           []int{2, 0},
		BIOSPause:                false,
		BIOSSetup:                false,
		BootMediaIndex:           0,
		BootguardStatus:          127,
		ConfigurationDataReset:   false,
		ElementName:              "Intel(r) AMT Boot Configuration Settings",
		EnforceSecureBoot:        false,
		FirmwareVerbosity:        0,
		ForcedProgressEvents:     false,
		IDERBootDevice:           0,
		InstanceID:               "Intel(r) AMT:BootSettingData 0",
		LockKeyboard:             false,
		LockPowerButton:          false,
		LockResetButton:          false,
		LockSleepButton:          false,
		OptionsCleared:           true,
		OwningEntity:             "Intel(r) AMT",
		PlatformErase:            false,
		RPEEnabled:               false,
		RSEPassword:              "",
		ReflashBIOS:              false,
		SecureBootControlEnabled: false,
		SecureErase:              false,
		UEFIHTTPSBootEnabled:     false,
		UEFILocalPBABootEnabled:  false,
		UefiBootNumberOfParams:   0,
		UseIDER:                  false,
		UseSOL:                   false,
		UseSafeMode:              false,
		UserPasswordBypass:       false,
		WinREBootEnabled:         false,
	}

	device := &entity.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}

	bootSetting := dto.BootSetting{
		Action: 400,
		UseSOL: true,
	}

	powerActionRes := power.PowerActionResponse{ReturnValue: 5}

	tests := []test{
		{
			name: "success",
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(hmm)
				hmm.EXPECT().
					GetBootData().
					Return(bootResponse, nil)
				hmm.EXPECT().
					SetBootConfigRole(1).
					Return(powerActionRes, nil)
				hmm.EXPECT().
					ChangeBootOrder(string(cimBoot.PXE)).
					Return(cimBoot.ChangeBootOrder_OUTPUT{}, nil)
				hmm.EXPECT().
					SetBootData(gomock.Any()).
					Return(nil, nil)
				hmm.EXPECT().
					SendPowerAction(10).
					Return(powerActionRes, nil)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: powerActionRes,
			err: nil,
		},
		{
			name:    "GetById fails",
			manMock: func(_ *MockWSMAN, _ *MockManagement) {},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},
			res: power.PowerActionResponse{},
			err: devices.ErrDatabase,
		},
		{
			name: "GetBootData fails",
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(hmm)
				hmm.EXPECT().
					GetBootData().
					Return(boot.BootSettingDataResponse{}, ErrGeneral)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: power.PowerActionResponse{},
			err: ErrGeneral,
		},
		{
			name: "SetBootConfigRole fails",
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(hmm)
				hmm.EXPECT().
					GetBootData().
					Return(bootResponse, nil)
				hmm.EXPECT().
					SetBootConfigRole(1).
					Return(powerActionRes, ErrGeneral)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: power.PowerActionResponse{},
			err: ErrGeneral,
		},
		{
			name: "ChangeBootOrder fails",
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(hmm)
				hmm.EXPECT().
					GetBootData().
					Return(bootResponse, nil)
				hmm.EXPECT().
					SetBootConfigRole(1).
					Return(powerActionRes, nil)
				hmm.EXPECT().
					ChangeBootOrder(string(cimBoot.PXE)).
					Return(cimBoot.ChangeBootOrder_OUTPUT{}, ErrGeneral)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: power.PowerActionResponse{},
			err: ErrGeneral,
		},
		{
			name: "SetBootData fails",
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(hmm)
				hmm.EXPECT().
					GetBootData().
					Return(bootResponse, nil)
				hmm.EXPECT().
					SetBootConfigRole(1).
					Return(powerActionRes, nil)
				hmm.EXPECT().
					ChangeBootOrder(string(cimBoot.PXE)).
					Return(cimBoot.ChangeBootOrder_OUTPUT{}, nil)
				hmm.EXPECT().
					SetBootData(gomock.Any()).
					Return(nil, ErrGeneral)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: power.PowerActionResponse{},
			err: ErrGeneral,
		},
		{
			name: "GetPowerCapabilities fails",
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(hmm)
				hmm.EXPECT().
					GetBootData().
					Return(bootResponse, nil)
				hmm.EXPECT().
					SetBootConfigRole(1).
					Return(powerActionRes, nil)
				hmm.EXPECT().
					ChangeBootOrder(string(cimBoot.PXE)).
					Return(cimBoot.ChangeBootOrder_OUTPUT{}, nil)
				hmm.EXPECT().
					SetBootData(gomock.Any()).
					Return(nil, nil)
				hmm.EXPECT().
					SendPowerAction(10).
					Return(powerActionRes, ErrGeneral)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: power.PowerActionResponse{},
			err: ErrGeneral,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, wsmanMock, management, repo := initPowerTest(t)
			tc.manMock(wsmanMock, management)
			tc.repoMock(repo)

			res, err := useCase.SetBootOptions(context.Background(), device.GUID, bootSetting)

			require.Equal(t, tc.res, res)
			require.IsType(t, tc.err, err)
		})
	}
}
