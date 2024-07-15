package devices_test

import (
	"context"
	"encoding/xml"
	"testing"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/auditlog"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/messagelog"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/setupandconfiguration"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/software"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	devices "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

func initInfoTest(t *testing.T) (*devices.UseCase, *MockWSMAN, *MockManagement, *MockRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)

	defer mockCtl.Finish()

	repo := NewMockRepository(mockCtl)

	wsmanMock := NewMockWSMAN(mockCtl)
	management := NewMockManagement(mockCtl)

	log := logger.New("error")

	u := devices.New(repo, wsmanMock, NewMockRedirection(mockCtl), log)

	return u, wsmanMock, management, repo
}

func TestGetVersion(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID: "device-guid-123",

		TenantID: "tenant-id-456",
	}

	softwares := []software.SoftwareIdentity{}

	responses := []setupandconfiguration.SetupAndConfigurationServiceResponse{}

	response := setupandconfiguration.SetupAndConfigurationServiceResponse{
		XMLName: xml.Name{Local: "AMT_SetupAndConfigurationService"},

		RequestedState: 1,

		EnabledState: 1,

		ElementName: "SampleElementName",

		SystemCreationClassName: "SampleSystemCreationClassName",

		SystemName: "SampleSystemName",

		CreationClassName: "SampleCreationClassName",

		Name: "SampleName",

		ProvisioningMode: 1,

		ProvisioningState: 1,

		ZeroTouchConfigurationEnabled: true,

		ProvisioningServerOTP: "SampleProvisioningServerOTP",

		ConfigurationServerFQDN: "SampleConfigurationServerFQDN",

		PasswordModel: 1,

		DhcpDNSSuffix: "SampleDhcpDNSSuffix",

		TrustedDNSSuffix: "SampleTrustedDNSSuffix",
	}

	responses = append(responses, response)

	tests := []test{
		{
			name: "success",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetAMTVersion().
					Return(softwares, nil)

				man2.EXPECT().
					GetSetupAndConfiguration().
					Return(responses, nil)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: map[string]interface{}{"AMT_SetupAndConfigurationService": map[string]interface{}{"response": setupandconfiguration.SetupAndConfigurationServiceResponse{XMLName: xml.Name{Space: "", Local: "AMT_SetupAndConfigurationService"}, RequestedState: 1, EnabledState: 1, ElementName: "SampleElementName", SystemCreationClassName: "SampleSystemCreationClassName", SystemName: "SampleSystemName", CreationClassName: "SampleCreationClassName", Name: "SampleName", ProvisioningMode: 1, ProvisioningState: 1, ZeroTouchConfigurationEnabled: true, ProvisioningServerOTP: "SampleProvisioningServerOTP", ConfigurationServerFQDN: "SampleConfigurationServerFQDN", PasswordModel: 1, DhcpDNSSuffix: "SampleDhcpDNSSuffix", TrustedDNSSuffix: "SampleTrustedDNSSuffix"}}, "CIM_SoftwareIdentity": map[string]interface{}{"responses": []software.SoftwareIdentity{}}},

			err: nil,
		},

		{
			name: "GetById fails",

			action: 0,

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
			name: "GetAMTVersion fails",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetAMTVersion().
					Return(softwares, ErrGeneral)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: map[string]interface{}(nil),

			err: ErrGeneral,
		},

		{
			name: "GetSetupAndConfiguration fails",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetAMTVersion().
					Return(softwares, nil)

				man2.EXPECT().
					GetSetupAndConfiguration().
					Return(responses, ErrGeneral)
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

			useCase, wsmanMock, management, repo := initInfoTest(t)

			tc.manMock(wsmanMock, management)

			tc.repoMock(repo)

			res, err := useCase.GetVersion(context.Background(), device.GUID)

			require.Equal(t, tc.res, res)

			require.IsType(t, tc.err, err)
		})
	}
}

func TestGetFeatures(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID: "device-guid-123",

		TenantID: "tenant-id-456",
	}

	tests := []test{
		{
			name: "success",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetFeatures().
					Return(dto.Features{}, nil)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: dto.Features{},

			err: nil,
		},

		{
			name: "GetById fails",

			action: 0,

			manMock: func(_ *MockWSMAN, _ *MockManagement) {},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},

			res: dto.Features{},

			err: devices.ErrDatabase,
		},

		{
			name: "GetFeatures fails",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetFeatures().
					Return(dto.Features{}, ErrGeneral)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: dto.Features{},

			err: ErrGeneral,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, wsmanMock, management, repo := initInfoTest(t)

			tc.manMock(wsmanMock, management)

			tc.repoMock(repo)

			res, err := useCase.GetFeatures(context.Background(), device.GUID)

			require.Equal(t, tc.res, res)

			require.IsType(t, tc.err, err)
		})
	}
}

func TestSetFeatures(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID: "device-guid-123",

		TenantID: "tenant-id-456",
	}

	featureSet := dto.Features{
		UserConsent: "kvm",

		EnableSOL: true,

		EnableIDER: true,

		EnableKVM: true,
	}

	tests := []test{
		{
			name: "success",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					SetFeatures(featureSet).
					Return(featureSet, nil)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: featureSet,

			err: nil,
		},

		{
			name: "GetById fails",

			action: 0,

			manMock: func(_ *MockWSMAN, _ *MockManagement) {},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},

			res: featureSet,

			err: devices.ErrDatabase,
		},

		{
			name: "GetFeatures fails",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					SetFeatures(featureSet).
					Return(featureSet, ErrGeneral)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: featureSet,

			err: ErrGeneral,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, wsmanMock, management, repo := initInfoTest(t)

			tc.manMock(wsmanMock, management)

			tc.repoMock(repo)

			res, err := useCase.SetFeatures(context.Background(), device.GUID, featureSet)

			require.Equal(t, tc.res, res)

			require.IsType(t, tc.err, err)
		})
	}
}

func TestGetHardwareInfo(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID: "device-guid-123",

		TenantID: "tenant-id-456",
	}

	tests := []test{
		{
			name: "success",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetHardwareInfo().
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

			manMock: func(_ *MockWSMAN, _ *MockManagement) {},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},

			res: nil,

			err: devices.ErrDatabase,
		},

		{
			name: "GetFeatures fails",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetHardwareInfo().
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

			useCase, wsmanMock, management, repo := initInfoTest(t)

			tc.manMock(wsmanMock, management)

			tc.repoMock(repo)

			res, err := useCase.GetHardwareInfo(context.Background(), device.GUID)

			require.Equal(t, tc.res, res)

			require.IsType(t, tc.err, err)
		})
	}
}

func TestGetAuditLog(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID: "device-guid-123",

		TenantID: "tenant-id-456",
	}

	tests := []test{
		{
			name: "success",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetAuditLog(1).
					Return(auditlog.Response{}, nil)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: dto.AuditLog{},

			err: nil,
		},

		{
			name: "GetById fails",

			action: 0,

			manMock: func(_ *MockWSMAN, _ *MockManagement) {},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},

			res: dto.AuditLog{TotalCount: 0, Records: nil},

			err: devices.ErrDatabase,
		},

		{
			name: "GetFeatures fails",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetAuditLog(1).
					Return(auditlog.Response{}, ErrGeneral)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: dto.AuditLog{TotalCount: 0, Records: nil},

			err: ErrGeneral,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, wsmanMock, management, repo := initInfoTest(t)

			tc.manMock(wsmanMock, management)

			tc.repoMock(repo)

			res, err := useCase.GetAuditLog(context.Background(), 1, device.GUID)

			require.Equal(t, tc.res, res)

			require.IsType(t, tc.err, err)
		})
	}
}

func TestGetEventLog(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID: "device-guid-123",

		TenantID: "tenant-id-456",
	}

	tests := []test{
		{
			name: "success",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetEventLog().
					Return(messagelog.GetRecordsResponse{}, nil)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: []dto.EventLog{},

			err: nil,
		},

		{
			name: "GetById fails",

			action: 0,

			manMock: func(_ *MockWSMAN, _ *MockManagement) {},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},

			res: []dto.EventLog(nil),

			err: devices.ErrDatabase,
		},

		{
			name: "GetFeatures fails",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetEventLog().
					Return(messagelog.GetRecordsResponse{}, ErrGeneral)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: []dto.EventLog(nil),

			err: ErrGeneral,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, wsmanMock, management, repo := initInfoTest(t)

			tc.manMock(wsmanMock, management)

			tc.repoMock(repo)

			res, err := useCase.GetEventLog(context.Background(), device.GUID)

			require.Equal(t, tc.res, res)

			require.IsType(t, tc.err, err)
		})
	}
}

func TestGetGeneralSettings(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID: "device-guid-123",

		TenantID: "tenant-id-456",
	}

	tests := []test{
		{
			name: "success",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetGeneralSettings().
					Return(gomock.Any(), nil)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: map[string]interface{}{"Body": gomock.Any()},

			err: nil,
		},

		{
			name: "GetById fails",

			action: 0,

			manMock: func(_ *MockWSMAN, _ *MockManagement) {},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},

			res: nil,

			err: devices.ErrDatabase,
		},

		{
			name: "GetFeatures fails",

			action: 0,

			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(man2)
				man2.EXPECT().
					GetGeneralSettings().
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

			useCase, wsmanMock, management, repo := initInfoTest(t)

			tc.manMock(wsmanMock, management)

			tc.repoMock(repo)

			res, err := useCase.GetGeneralSettings(context.Background(), device.GUID)

			require.Equal(t, tc.res, res)

			require.IsType(t, tc.err, err)
		})
	}
}
