package devices_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type test struct {
	name     string
	guid     string
	tenantID string
	top      int
	skip     int
	mock     func(*MockRepository)
	res      interface{}
	err      error
}

func devicesTest(t *testing.T) (*devices.UseCase, *MockRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockRepository(mockCtl)
	management := NewMockManagement(mockCtl)
	log := logger.New("error")
	u := devices.New(repo, management, NewMockRedirection(mockCtl), log)

	return u, repo
}

func TestGetCount(t *testing.T) {
	t.Parallel()

	tests := []test{
		{
			name: "empty result",
			mock: func(repo *MockRepository) {
				repo.EXPECT().GetCount(context.Background(), "").Return(0, nil)
			},
			res: 0,
			err: nil,
		},
		{
			name: "result with error",
			mock: func(repo *MockRepository) {
				repo.EXPECT().GetCount(context.Background(), "").Return(0, devices.ErrDatabase)
			},
			res: 0,
			err: devices.ErrDatabase,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := devicesTest(t)

			tc.mock(repo)

			res, err := useCase.GetCount(context.Background(), tc.tenantID)

			require.Equal(t, tc.res, res)
			require.IsType(t, tc.err, err)
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	testDevices := []entity.Device{
		{
			GUID:     "guid-123",
			TenantID: "tenant-id-456",
		},
		{
			GUID:     "guid-456",
			TenantID: "tenant-id-456",
		},
	}

	testDeviceDTOs := []dto.Device{
		{
			GUID:     "guid-123",
			TenantID: "tenant-id-456",
			Tags:     []string{""},
		},
		{
			GUID:     "guid-456",
			TenantID: "tenant-id-456",
			Tags:     []string{""},
		},
	}

	tests := []test{
		{
			name:     "successful retrieval",
			top:      10,
			skip:     0,
			tenantID: "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Get(context.Background(), 10, 0, "tenant-id-456").
					Return(testDevices, nil)
			},
			res: testDeviceDTOs,
			err: nil,
		},
		{
			name:     "database error",
			top:      5,
			skip:     0,
			tenantID: "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Get(context.Background(), 5, 0, "tenant-id-456").
					Return(nil, devices.ErrDatabase)
			},
			res: []dto.Device(nil),
			err: devices.ErrDatabase,
		},
		{
			name:     "zero results",
			top:      10,
			skip:     20,
			tenantID: "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Get(context.Background(), 10, 20, "tenant-id-456").
					Return([]entity.Device{}, nil)
			},
			res: []dto.Device{},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := devicesTest(t)

			tc.mock(repo)

			results, err := useCase.Get(context.Background(), tc.top, tc.skip, tc.tenantID)

			require.Equal(t, tc.res, results)

			if tc.err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}
	deviceDTO := &dto.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
		Tags:     []string{""},
	}

	tests := []test{
		{
			name:     "successful retrieval",
			guid:     "device-guid-123",
			tenantID: "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), "device-guid-123", "tenant-id-456").
					Return(device, nil)
			},
			res: deviceDTO,
			err: nil,
		},
		{
			name:     "device not found",
			guid:     "device-guid-unknown",
			tenantID: "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), "device-guid-unknown", "tenant-id-456").
					Return(nil, nil)
			},
			res: nil,
			err: devices.ErrNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := devicesTest(t)

			tc.mock(repo)

			got, err := useCase.GetByID(context.Background(), tc.guid, tc.tenantID)

			if tc.err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, got)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []test{
		{
			name:     "successful deletion",
			guid:     "guid-123",
			tenantID: "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Delete(context.Background(), "guid-123", "tenant-id-456").
					Return(true, nil)
			},
			err: nil,
		},
		{
			name:     "deletion fails - device not found",
			guid:     "guid-456",
			tenantID: "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Delete(context.Background(), "guid-456", "tenant-id-456").
					Return(false, nil)
			},
			err: devices.ErrNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo := devicesTest(t)
			tc.mock(repo)

			err := useCase.Delete(context.Background(), tc.guid, tc.tenantID)

			if tc.err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	device := &entity.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
	}

	deviceDTO := &dto.Device{
		GUID:     "device-guid-123",
		TenantID: "tenant-id-456",
		Tags:     []string{""},
	}

	tests := []test{
		{
			name: "successful update",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Update(context.Background(), device).
					Return(true, nil)
				repo.EXPECT().
					GetByID(gomock.Any(), "device-guid-123", "tenant-id-456").
					Return(device, nil)
			},
			res: deviceDTO,
			err: nil,
		},
		{
			name: "update fails - database error",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Update(context.Background(), device).
					Return(false, devices.ErrDatabase)
			},
			res: (*dto.Device)(nil),
			err: devices.ErrDatabase,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo := devicesTest(t)
			tc.mock(repo)

			result, err := useCase.Update(context.Background(), deviceDTO)

			require.Equal(t, tc.res, result)
			require.IsType(t, tc.err, err)
		})
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	tests := []test{
		{
			name: "successful insertion",
			mock: func(repo *MockRepository) {
				device := &entity.Device{
					GUID:     "device-guid-123",
					TenantID: "tenant-id-456",
				}

				repo.EXPECT().
					Insert(context.Background(), device).
					Return("unique-device-id", nil)
				repo.EXPECT().
					GetByID(gomock.Any(), device.GUID, "tenant-id-456").
					Return(device, nil)
			},
			res: nil, // little bit different in that the expectation is handled in the loop
			err: nil,
		},
		{
			name: "insertion fails - database error",
			mock: func(repo *MockRepository) {
				device := &entity.Device{
					GUID:     "device-guid-123",
					TenantID: "tenant-id-456",
				}

				repo.EXPECT().
					Insert(context.Background(), device).
					Return("", devices.ErrDatabase)
			},
			res: (*dto.Device)(nil),
			err: devices.ErrDatabase,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo := devicesTest(t)
			tc.mock(repo)

			deviceDTO := &dto.Device{
				GUID:     "device-guid-123",
				TenantID: "tenant-id-456",
				Tags:     []string{""},
			}

			insertedDevice, err := useCase.Insert(context.Background(), deviceDTO)

			if tc.err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.err.Error())
				require.Equal(t, tc.res, insertedDevice)
			} else {
				require.NoError(t, err)
				require.Equal(t, deviceDTO.TenantID, insertedDevice.TenantID)
				require.NotEmpty(t, deviceDTO.GUID)
			}
		})
	}
}
