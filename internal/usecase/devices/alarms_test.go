package devices_test

import (
	"context"
	"testing"
	"time"

	amtAlarmClock "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/alarmclock"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/alarmclock"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	devices "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

func initAlarmsTest(t *testing.T) (*devices.UseCase, *MockManagement, *MockRepository) {
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

func TestGetAlarmOccurrences(t *testing.T) {
	t.Parallel()

	dtoDevice := &dto.Device{
		GUID:     "device-guid-123",
		Tags:     []string{""},
		TenantID: "tenant-id-456",
	}

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
					SetupWsmanClient(*dtoDevice, false, true).
					Return()

				man.EXPECT().
					GetAlarmOccurrences().
					Return([]alarmclock.AlarmClockOccurrence{}, nil)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: []alarmclock.AlarmClockOccurrence{},

			err: nil,
		},

		{
			name: "GetById fails",

			action: 0,

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},

			res: []alarmclock.AlarmClockOccurrence(nil),

			err: devices.ErrDatabase,
		},

		{
			name: "GetAlarmOccurrences fails",

			action: 0,

			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return()

				man.EXPECT().
					GetAlarmOccurrences().
					Return([]alarmclock.AlarmClockOccurrence{}, ErrGeneral)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: []alarmclock.AlarmClockOccurrence(nil),

			err: ErrGeneral,
		},

		{
			name: "GetAlarmOccurrences returns nil",

			action: 0,

			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(*dtoDevice, false, true).
					Return()

				man.EXPECT().
					GetAlarmOccurrences().
					Return(nil, nil)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: []alarmclock.AlarmClockOccurrence{},

			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, management, repo := initAlarmsTest(t)

			if tc.manMock != nil {
				tc.manMock(management)
			}

			tc.repoMock(repo)

			res, err := useCase.GetAlarmOccurrences(context.Background(), device.GUID)

			require.Equal(t, tc.res, res)

			require.IsType(t, tc.err, err)
		})
	}
}

func TestCreateAlarmOccurrences(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID: "device-guid-123",

		TenantID: "tenant-id-456",
	}
	dtoDevice := &dto.Device{
		GUID:     "device-guid-123",
		Tags:     []string{""},
		TenantID: "tenant-id-456",
	}
	occ := dto.AlarmClockOccurrence{
		ElementName: "test",

		InstanceID: "test",

		StartTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),

		Interval: 1,

		DeleteOnCompletion: true,
	}

	tests := []struct {
		name string

		action int

		manMock func(man *MockManagement)

		repoMock func(repo *MockRepository)

		res amtAlarmClock.AddAlarmOutput

		err error
	}{
		{
			name: "success",

			action: 0,

			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(*dtoDevice, false, true).
					Return()

				man.EXPECT().
					CreateAlarmOccurrences(occ.InstanceID, occ.StartTime, 1, occ.DeleteOnCompletion).
					Return(amtAlarmClock.AddAlarmOutput{}, nil)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: amtAlarmClock.AddAlarmOutput{},

			err: nil,
		},

		{
			name: "GetByID fails",

			action: 0,

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},

			res: amtAlarmClock.AddAlarmOutput{},

			err: devices.ErrDatabase,
		},

		{
			name: "GetAlarmOccurrences fails",

			action: 0,

			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(*dtoDevice, false, true).
					Return()

				man.EXPECT().
					CreateAlarmOccurrences(occ.InstanceID, occ.StartTime, 1, occ.DeleteOnCompletion).
					Return(amtAlarmClock.AddAlarmOutput{}, ErrGeneral)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			res: amtAlarmClock.AddAlarmOutput{},

			err: devices.ErrAMT,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, management, repo := initAlarmsTest(t)

			if tc.manMock != nil {
				tc.manMock(management)
			}

			tc.repoMock(repo)

			res, err := useCase.CreateAlarmOccurrences(context.Background(), device.GUID, occ)

			require.Equal(t, tc.res, res)

			require.IsType(t, tc.err, err)
		})
	}
}

func TestDeleteAlarmOccurrences(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID: "device-guid-123",

		TenantID: "tenant-id-456",
	}

	dtoDevice := &dto.Device{
		GUID:     "device-guid-123",
		Tags:     []string{""},
		TenantID: "tenant-id-456",
	}
	tests := []struct {
		name string

		action int

		manMock func(man *MockManagement)

		repoMock func(repo *MockRepository)

		err error
	}{
		{
			name: "success",

			action: 0,

			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(*dtoDevice, false, true).
					Return()

				man.EXPECT().
					DeleteAlarmOccurrences("").
					Return(nil)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			err: nil,
		},

		{
			name: "GetById fails",

			action: 0,

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},

			err: devices.ErrDatabase,
		},

		{
			name: "GetAlarmOccurrences fails",

			action: 0,

			manMock: func(man *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(*dtoDevice, false, true).
					Return()

				man.EXPECT().
					DeleteAlarmOccurrences("").
					Return(ErrGeneral)
			},

			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},

			err: ErrGeneral,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, management, repo := initAlarmsTest(t)

			if tc.manMock != nil {
				tc.manMock(management)
			}

			tc.repoMock(repo)

			err := useCase.DeleteAlarmOccurrences(context.Background(), device.GUID, "")

			require.IsType(t, tc.err, err)
		})
	}
}
