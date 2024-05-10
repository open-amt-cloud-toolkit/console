package profiles_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/profiles"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type test struct {
	name        string
	top         int
	skip        int
	tenantID    string
	profileName string
	input       entity.Profile
	mock        func(*MockRepository)
	res         interface{}
	err         error
}

func profilesTest(t *testing.T) (*profiles.UseCase, *MockRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockRepository(mockCtl)
	log := logger.New("error")
	useCase := profiles.New(repo, log)

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
				repo.EXPECT().GetCount(context.Background(), "").Return(0, profiles.ErrDatabase)
			},
			res: 0,
			err: profiles.ErrDatabase,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo := profilesTest(t)

			tc.mock(repo)

			res, err := useCase.GetCount(context.Background(), "")

			require.Equal(t, tc.res, res)
			require.IsType(t, tc.err, err)
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	testProfiles := []entity.Profile{
		{
			ProfileName: "test-profile-1",
			TenantID:    "tenant-id-456",
		},
		{
			ProfileName: "test-profile-2",
			TenantID:    "tenant-id-456",
		},
	}

	testProfileDTOs := []dto.Profile{
		{
			ProfileName: "test-profile-1",
			TenantID:    "tenant-id-456",
			Tags:        []string{""},
		},
		{
			ProfileName: "test-profile-2",
			TenantID:    "tenant-id-456",
			Tags:        []string{""},
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
					Return(testProfiles, nil)
			},
			res: testProfileDTOs,
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
					Return(nil, profiles.ErrDatabase)
			},
			res: []dto.Profile(nil),
			err: profiles.ErrDatabase,
		},
		{
			name:     "zero results",
			top:      10,
			skip:     20,
			tenantID: "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Get(context.Background(), 10, 20, "tenant-id-456").
					Return([]entity.Profile{}, nil)
			},
			res: []dto.Profile{},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo := profilesTest(t)

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

	profile := &entity.Profile{
		ProfileName: "test-profile",
		TenantID:    "tenant-id-456",
		Version:     "1.0.0",
	}

	profileDTO := &dto.Profile{
		ProfileName: "test-profile",
		TenantID:    "tenant-id-456",
		Version:     "1.0.0",
		Tags:        []string{""},
	}

	tests := []test{
		{
			name: "successful retrieval",
			input: entity.Profile{
				ProfileName: "test-profile",
				TenantID:    "tenant-id-456",
			},
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByName(context.Background(), "test-profile", "tenant-id-456").
					Return(profile, nil)
			},
			res: profileDTO,
			err: nil,
		},
		{
			name: "profile not found",
			input: entity.Profile{
				ProfileName: "unknown-profile",
				TenantID:    "tenant-id-456",
			},
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByName(context.Background(), "unknown-profile", "tenant-id-456").
					Return(nil, nil)
			},
			res: (*dto.Profile)(nil),
			err: profiles.ErrNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo := profilesTest(t)

			tc.mock(repo)

			res, err := useCase.GetByName(context.Background(), tc.input.ProfileName, tc.input.TenantID)

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
			name:        "successful deletion",
			profileName: "example-profile",
			tenantID:    "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Delete(context.Background(), "example-profile", "tenant-id-456").
					Return(true, nil)
			},
			err: nil,
		},
		{
			name:        "deletion fails - profile not found",
			profileName: "nonexistent-profile",
			tenantID:    "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Delete(context.Background(), "nonexistent-profile", "tenant-id-456").
					Return(false, nil)
			},
			err: profiles.ErrNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo := profilesTest(t)

			tc.mock(repo)

			err := useCase.Delete(context.Background(), tc.profileName, tc.tenantID)

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

	profile := &entity.Profile{
		ProfileName: "example-profile",
		TenantID:    "tenant-id-456",
		Version:     "1.0.0",
	}

	profileDTO := &dto.Profile{
		ProfileName: "example-profile",
		TenantID:    "tenant-id-456",
		Version:     "1.0.0",
		Tags:        []string{""},
	}

	tests := []test{
		{
			name: "successful update",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Update(context.Background(), profile).
					Return(true, nil)
				repo.EXPECT().
					GetByName(context.Background(), profile.ProfileName, profile.TenantID).
					Return(profile, nil)
			},
			res: profileDTO,
			err: nil,
		},
		{
			name: "update fails - database error",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Update(context.Background(), profile).
					Return(false, profiles.ErrDatabase)
			},
			res: (*dto.Profile)(nil),
			err: profiles.ErrDatabase,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo := profilesTest(t)

			tc.mock(repo)

			result, err := useCase.Update(context.Background(), profileDTO)

			require.Equal(t, tc.res, result)
			require.IsType(t, err, tc.err)
		})
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	profile := &entity.Profile{
		ProfileName: "new-profile",
		TenantID:    "tenant-id-789",
		Version:     "1.0.0",
		Tags:        "",
	}

	profileDTO := &dto.Profile{
		ProfileName: "new-profile",
		TenantID:    "tenant-id-789",
		Version:     "1.0.0",
		Tags:        []string{""},
	}

	tests := []test{
		{
			name: "successful insertion",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Insert(context.Background(), profile).
					Return("unique-profile-id", nil)
				repo.EXPECT().
					GetByName(context.Background(), profile.ProfileName, profile.TenantID).
					Return(profile, nil)
			},
			res: profileDTO,
			err: nil,
		},
		{
			name: "insertion fails - database error",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Insert(context.Background(), profile).
					Return("", profiles.ErrDatabase)
			},
			res: (*dto.Profile)(nil),
			err: profiles.ErrDatabase,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo := profilesTest(t)

			tc.mock(repo)

			id, err := useCase.Insert(context.Background(), profileDTO)

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
