package wificonfigs_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/wificonfigs"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

var (
	errInternalServerErr = errors.New("internal server error")
	errDB                = errors.New("database error")
	errNotFound          = errors.New("wirelessconfig not found")
	errDelete            = fmt.Errorf("WificonfigsUseCase - Delete - s.repo.Delete: wirelessconfig not found")
	errGetByName         = fmt.Errorf("WificonfigsUseCase - GetByName - s.repo.GetByName: wirelessconfig not found")
)

type test struct {
	name           string
	top            int
	skip           int
	input          entity.WirelessConfig
	profileName    string
	tenantID       string
	mock           func(*MockRepository, ...interface{})
	expectedResult interface{}
	err            error
}

func wificonfigsTest(t *testing.T) (*wificonfigs.UseCase, *MockRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockRepository(mockCtl)
	log := logger.New("error")
	useCase := wificonfigs.New(repo, log)

	return useCase, repo
}

func TestCheckProfileExists(t *testing.T) {
	t.Parallel()

	tests := []test{
		{
			name:        "empty result",
			profileName: "example-wirelessconfig",
			tenantID:    "tenant-id-456",
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().CheckProfileExists(context.Background(), args[0], args[1]).Return(false, nil)
			},
			expectedResult: false,
			err:            nil,
		},
		{
			name:        "result with error",
			profileName: "nonexistent-wirelessconfig",
			tenantID:    "tenant-id-456",
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().CheckProfileExists(context.Background(), args[0], args[1]).Return(false, errInternalServerErr)
			},
			expectedResult: false,
			err:            errInternalServerErr,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := wificonfigsTest(t)

			tc.mock(repo, tc.profileName, tc.tenantID)

			res, err := useCase.CheckProfileExists(context.Background(), tc.profileName, tc.tenantID)

			require.Equal(t, tc.expectedResult, res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestGetCount(t *testing.T) {
	t.Parallel()

	tests := []test{
		{
			name: "empty result",
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().GetCount(context.Background(), "").Return(args[0], args[1])
			},
			expectedResult: 0,
			err:            nil,
		},
		{
			name: "result with error",
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().GetCount(context.Background(), "").Return(args[0], args[1])
			},
			expectedResult: 0,
			err:            errInternalServerErr,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo := wificonfigsTest(t)
			tc.mock(repo, tc.expectedResult, tc.err)

			res, err := useCase.GetCount(context.Background(), "")

			require.Equal(t, res, tc.expectedResult)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	testWifiConfigs := []entity.WirelessConfig{
		{
			ProfileName: "test-wirelessconfig-1",
			TenantID:    "tenant-id-456",
		},
		{
			ProfileName: "test-wirelessconfig-2",
			TenantID:    "tenant-id-456",
		},
	}

	tests := []test{
		{
			name:     "successful retrieval",
			top:      10,
			skip:     0,
			tenantID: "tenant-id-456",
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().
					Get(context.Background(), args[0], args[1], args[2]).
					Return(testWifiConfigs, nil)
			},
			expectedResult: testWifiConfigs,
			err:            nil,
		},
		{
			name:     "database error",
			top:      5,
			skip:     0,
			tenantID: "tenant-id-456",
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().
					Get(context.Background(), args[0], args[1], args[2]).
					Return(nil, errDB)
			},
			expectedResult: []entity.WirelessConfig(nil),
			err:            errDB,
		},
		{
			name:     "zero results",
			top:      10,
			skip:     20,
			tenantID: "tenant-id-456",
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().
					Get(context.Background(), args[0], args[1], args[2]).
					Return([]entity.WirelessConfig{}, nil)
			},
			expectedResult: []entity.WirelessConfig{},
			err:            nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo := wificonfigsTest(t)

			tc.mock(repo, tc.top, tc.skip, tc.tenantID)

			results, err := useCase.Get(context.Background(), tc.top, tc.skip, tc.tenantID)

			require.Equal(t, tc.expectedResult, results)

			if tc.err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetByName(t *testing.T) {
	t.Parallel()

	wirelessConfig := entity.WirelessConfig{
		ProfileName: "test-WirelessConfig",
		TenantID:    "tenant-id-456",
		Version:     "123",
	}

	tests := []test{
		{
			name: "successful retrieval",
			input: entity.WirelessConfig{
				ProfileName: "test-wirelessConfig",
				TenantID:    "tenant-id-456",
			},
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().
					GetByName(context.Background(), args[0], args[1]).
					Return(wirelessConfig, nil)
			},
			expectedResult: wirelessConfig,
			err:            nil,
		},
		{
			name: "WirelessConfig not found",
			input: entity.WirelessConfig{
				ProfileName: "unknown-WirelessConfig",
				TenantID:    "tenant-id-456",
			},
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().
					GetByName(context.Background(), args[0], args[1]).
					Return(entity.WirelessConfig{}, errNotFound)
			},
			expectedResult: entity.WirelessConfig{},
			err:            errGetByName,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo := wificonfigsTest(t)

			tc.mock(repo, tc.input.ProfileName, tc.input.TenantID)

			res, err := useCase.GetByName(context.Background(), tc.input.ProfileName, tc.input.TenantID)

			require.Equal(t, tc.expectedResult, res)

			if tc.err != nil {
				require.Contains(t, err.Error(), tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []test{
		{
			name:        "successful deletion",
			profileName: "example-wirelessconfig",
			tenantID:    "tenant-id-456",
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().
					Delete(context.Background(), args[0], args[1]).
					Return(true, nil)
			},
			expectedResult: true,
			err:            nil,
		},
		{
			name:        "deletion fails - wirelessconfig not found",
			profileName: "nonexistent-wirelessconfig",
			tenantID:    "tenant-id-456",
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().
					Delete(context.Background(), args[0], args[1]).
					Return(false, errNotFound)
			},
			expectedResult: false,
			err:            errDelete,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo := wificonfigsTest(t)

			tc.mock(repo, tc.profileName, tc.tenantID)

			result, err := useCase.Delete(context.Background(), tc.profileName, tc.tenantID)

			require.Equal(t, tc.expectedResult, result)

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

	tests := []test{
		{
			name: "successful update",
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().
					Update(context.Background(), args[0]).
					Return(true, nil)
			},
			expectedResult: true,
			err:            nil,
		},
		{
			name: "update fails - database error",
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().
					Update(context.Background(), args[0]).
					Return(false, errInternalServerErr)
			},
			expectedResult: false,
			err:            errInternalServerErr,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			wirelessconfig := &entity.WirelessConfig{
				ProfileName: "example-wirelessconfig",
				TenantID:    "tenant-id-456",
				Version:     "123",
			}

			useCase, repo := wificonfigsTest(t)

			tc.mock(repo, wirelessconfig)

			result, err := useCase.Update(context.Background(), wirelessconfig)

			require.Equal(t, tc.expectedResult, result)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	tests := []test{
		{
			name: "successful insertion",
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().
					Insert(context.Background(), args[0]).
					Return("unique-wirelessconfig-id", nil)
			},
			expectedResult: "unique-wirelessconfig-id",
			err:            nil,
		},
		{
			name: "insertion fails - database error",
			mock: func(repo *MockRepository, args ...interface{}) {
				repo.EXPECT().
					Insert(context.Background(), args[0]).
					Return("", errInternalServerErr)
			},
			expectedResult: "",
			err:            errInternalServerErr,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			wirelessconfig := &entity.WirelessConfig{
				ProfileName: "new-wirelessconfig",
				TenantID:    "tenant-id-789",
				Version:     "123",
			}

			useCase, repo := wificonfigsTest(t)

			tc.mock(repo, wirelessconfig)

			id, err := useCase.Insert(context.Background(), wirelessconfig)

			require.Equal(t, tc.expectedResult, id)

			if tc.err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
