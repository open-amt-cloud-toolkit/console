package ciraconfigs_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/ciraconfigs"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

var (
	errInternalServErr = errors.New("internal server error")
	errDB              = errors.New("database error")
	errGetByName       = fmt.Errorf("CIRAConfigsUseCase - GetByName - s.repo.GetByName: ciraconfig not found")
	errDelete          = fmt.Errorf("CIRAConfigsUseCase - Delete - s.repo.Delete: ciraconfig not found")
	errNotFound        = errors.New("ciraconfig not found")
)

type test struct {
	name       string
	top        int
	skip       int
	configName string
	tenantID   string
	input      entity.CIRAConfig
	mock       func(*MockRepository)
	res        interface{}
	err        error
}

func ciraconfigsTest(t *testing.T) (*ciraconfigs.UseCase, *MockRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockRepository(mockCtl)
	log := logger.New("error")
	useCase := ciraconfigs.New(repo, log)

	return useCase, repo
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
				repo.EXPECT().GetCount(context.Background(), "").Return(0, errInternalServErr)
			},
			res: 0,
			err: errInternalServErr,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := ciraconfigsTest(t)

			tc.mock(repo)

			res, err := useCase.GetCount(context.Background(), "")

			require.Equal(t, res, tc.res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	testCIRAConfigs := []entity.CIRAConfig{
		{
			ConfigName: "test-config-1",
			TenantID:   "tenant-id-456",
		},
		{
			ConfigName: "test-config-2",
			TenantID:   "tenant-id-456",
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
					Return(testCIRAConfigs, nil)
			},
			res: testCIRAConfigs,
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
					Return(nil, errDB)
			},
			res: []entity.CIRAConfig(nil),
			err: errDB,
		},
		{
			name:     "zero results",
			top:      10,
			skip:     20,
			tenantID: "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Get(context.Background(), 10, 20, "tenant-id-456").
					Return([]entity.CIRAConfig{}, nil)
			},
			res: []entity.CIRAConfig{},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := ciraconfigsTest(t)

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

func TestGetByName(t *testing.T) {
	t.Parallel()

	ciraconfig := entity.CIRAConfig{
		ConfigName: "test-config",
		TenantID:   "tenant-id-456",
		Version:    "1.0.0",
	}

	tests := []test{
		{
			name: "successful retrieval",
			input: entity.CIRAConfig{
				ConfigName: "test-config",
				TenantID:   "tenant-id-456",
			},
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByName(context.Background(), "test-config", "tenant-id-456").
					Return(ciraconfig, nil)
			},
			res: ciraconfig,
			err: nil,
		},
		{
			name: "ciraconfig not found",
			input: entity.CIRAConfig{
				ConfigName: "unknown-ciraconfig",
				TenantID:   "tenant-id-456",
			},
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByName(context.Background(), "unknown-ciraconfig", "tenant-id-456").
					Return(entity.CIRAConfig{}, errNotFound)
			},
			res: entity.CIRAConfig{},
			err: errGetByName,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := ciraconfigsTest(t)

			tc.mock(repo)

			res, err := useCase.GetByName(context.Background(), tc.input.ConfigName, tc.input.TenantID)

			require.Equal(t, tc.res, res)

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
			name:       "successful deletion",
			configName: "example-ciraconfig",
			tenantID:   "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Delete(context.Background(), "example-ciraconfig", "tenant-id-456").
					Return(true, nil)
			},
			res: true,
			err: nil,
		},
		{
			name:       "deletion fails - ciraconfig not found",
			configName: "nonexistent-ciraconfig",
			tenantID:   "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Delete(context.Background(), "nonexistent-ciraconfig", "tenant-id-456").
					Return(false, errNotFound)
			},
			res: false,
			err: errDelete,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := ciraconfigsTest(t)

			tc.mock(repo)

			result, err := useCase.Delete(context.Background(), tc.configName, tc.tenantID)

			require.Equal(t, tc.res, result)

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

	ciraconfig := &entity.CIRAConfig{
		ConfigName: "test-config",
		TenantID:   "tenant-id-456",
		Version:    "1.0.0",
	}

	tests := []test{
		{
			name: "successful update",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Update(context.Background(), ciraconfig).
					Return(true, nil)
			},
			res: true,
			err: nil,
		},
		{
			name: "update fails - database error",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Update(context.Background(), ciraconfig).
					Return(false, errInternalServErr)
			},
			res: false,
			err: errInternalServErr,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := ciraconfigsTest(t)

			tc.mock(repo)

			result, err := useCase.Update(context.Background(), ciraconfig)

			require.Equal(t, tc.res, result)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	ciraconfig := &entity.CIRAConfig{
		ConfigName: "test-config",
		TenantID:   "tenant-id-456",
		Version:    "1.0.0",
	}

	tests := []test{
		{
			name: "successful insertion",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Insert(context.Background(), ciraconfig).
					Return("unique-ciraconfig-id", nil)
			},
			res: "unique-ciraconfig-id",
			err: nil,
		},
		{
			name: "insertion fails - database error",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Insert(context.Background(), ciraconfig).
					Return("", errInternalServErr)
			},
			res: "",
			err: errInternalServErr,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := ciraconfigsTest(t)

			tc.mock(repo)

			id, err := useCase.Insert(context.Background(), ciraconfig)

			require.Equal(t, tc.res, id)

			if tc.err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
