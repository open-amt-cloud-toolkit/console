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
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	devices "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

func initAlarmsTest(t *testing.T) (*devices.UseCase, *MockWSMAN, *MockManagement, *MockRepository) {
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

func TestGetAlarmOccurrences(t *testing.T) {
	t.Parallel()

	dtoDevice := &dto.Device{
		GUID:     "device-guid-123",
		Tags:     nil,
		TenantID: "tenant-id-456",
	}

	device := &entity.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}

	tests := []test{
		{
			name:   "success",
			action: 0,
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(*dtoDevice, false, true).
					Return(hmm)
				hmm.EXPECT().
					GetAlarmOccurrences().
					Return([]alarmclock.AlarmClockOccurrence{}, nil)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: []dto.AlarmClockOccurrence{},
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
			res: []dto.AlarmClockOccurrence(nil),
			err: devices.ErrDatabase,
		},
		{
			name:   "GetAlarmOccurrences fails",
			action: 0,
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(gomock.Any(), false, true).
					Return(hmm)
				hmm.EXPECT().
					GetAlarmOccurrences().
					Return([]alarmclock.AlarmClockOccurrence{}, ErrGeneral)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: []dto.AlarmClockOccurrence(nil),
			err: ErrGeneral,
		},
		{
			name:   "GetAlarmOccurrences returns nil",
			action: 0,
			manMock: func(man *MockWSMAN, hmm *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(*dtoDevice, false, true).
					Return(hmm)
				hmm.EXPECT().
					GetAlarmOccurrences().
					Return(nil, nil)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: []dto.AlarmClockOccurrence{},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, wsmanMock, management, repo := initAlarmsTest(t)

			if tc.manMock != nil {
				tc.manMock(wsmanMock, management)
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
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}
	dtoDevice := &dto.Device{
		GUID:     "device-guid-123",
		Tags:     nil,
		TenantID: "tenant-id-456",
	}
	occ := dto.AlarmClockOccurrence{
		ElementName:        "test",
		InstanceID:         "test",
		StartTime:          time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Interval:           1,
		DeleteOnCompletion: true,
	}

	tests := []struct {
		name string

		action int

		manMock func(man *MockWSMAN, man2 *MockManagement)

		repoMock func(repo *MockRepository)

		res dto.AddAlarmOutput

		err error
	}{
		{
			name:   "success",
			action: 0,
			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(*dtoDevice, false, true).
					Return(man2)
				man2.EXPECT().
					CreateAlarmOccurrences(occ.InstanceID, occ.StartTime, 1, occ.DeleteOnCompletion).
					Return(amtAlarmClock.AddAlarmOutput{}, nil)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: dto.AddAlarmOutput{},
			err: nil,
		},
		{
			name:   "GetByID fails",
			action: 0,
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},
			res: dto.AddAlarmOutput{},
			err: devices.ErrDatabase,
		},
		{
			name:   "GetAlarmOccurrences fails",
			action: 0,
			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(*dtoDevice, false, true).
					Return(man2)
				man2.EXPECT().
					CreateAlarmOccurrences(occ.InstanceID, occ.StartTime, 1, occ.DeleteOnCompletion).
					Return(amtAlarmClock.AddAlarmOutput{}, ErrGeneral)
			},
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(device, nil)
			},
			res: dto.AddAlarmOutput{},
			err: devices.ErrAMT,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, wsmanMock, management, repo := initAlarmsTest(t)

			if tc.manMock != nil {
				tc.manMock(wsmanMock, management)
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
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}

	dtoDevice := &dto.Device{
		GUID:     "device-guid-123",
		Tags:     nil,
		TenantID: "tenant-id-456",
	}
	tests := []struct {
		name     string
		action   int
		manMock  func(man *MockWSMAN, man2 *MockManagement)
		repoMock func(repo *MockRepository)
		err      error
	}{
		{
			name:   "success",
			action: 0,
			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(*dtoDevice, false, true).
					Return(man2)
				man2.EXPECT().
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
			name:   "GetById fails",
			action: 0,
			repoMock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(context.Background(), device.GUID, "").
					Return(nil, ErrGeneral)
			},
			err: devices.ErrDatabase,
		},
		{
			name:   "GetAlarmOccurrences fails",
			action: 0,
			manMock: func(man *MockWSMAN, man2 *MockManagement) {
				man.EXPECT().
					SetupWsmanClient(*dtoDevice, false, true).
					Return(man2)
				man2.EXPECT().
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

			useCase, wsmanMock, management, repo := initAlarmsTest(t)

			if tc.manMock != nil {
				tc.manMock(wsmanMock, management)
			}

			tc.repoMock(repo)

			err := useCase.DeleteAlarmOccurrences(context.Background(), device.GUID, "")

			require.IsType(t, tc.err, err)
		})
	}
}
