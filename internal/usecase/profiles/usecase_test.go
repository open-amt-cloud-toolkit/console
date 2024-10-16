package profiles_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	"github.com/open-amt-cloud-toolkit/console/internal/mocks"
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
	mock        func(*mocks.MockProfilesRepository, *mocks.MockWiFiConfigsRepository, *mocks.MockProfileWiFiConfigsFeature)
	res         interface{}
	err         error
}

func profilesTest(t *testing.T) (*profiles.UseCase, *mocks.MockProfilesRepository, *mocks.MockWiFiConfigsRepository, *mocks.MockProfileWiFiConfigsFeature) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := mocks.NewMockProfilesRepository(mockCtl)
	wificonfigs := mocks.NewMockWiFiConfigsRepository(mockCtl)
	profilewificonfigs := mocks.NewMockProfileWiFiConfigsFeature(mockCtl)
	ieeeMock := mocks.NewMockIEEE8021xConfigsFeature(mockCtl)
	domains := mocks.NewMockDomainsRepository(mockCtl)
	security := mocks.MockCrypto{}
	log := logger.New("error")
	useCase := profiles.New(repo, wificonfigs, profilewificonfigs, ieeeMock, log, domains, security)

	return useCase, repo, wificonfigs, profilewificonfigs
}

func TestGetCount(t *testing.T) {
	t.Parallel()

	tests := []test{
		{
			name: "empty result",
			mock: func(repo *mocks.MockProfilesRepository, _ *mocks.MockWiFiConfigsRepository, _ *mocks.MockProfileWiFiConfigsFeature) {
				repo.EXPECT().GetCount(context.Background(), "").Return(0, nil)
			},
			res: 0,
			err: nil,
		},
		{
			name: "result with error",
			mock: func(repo *mocks.MockProfilesRepository, _ *mocks.MockWiFiConfigsRepository, _ *mocks.MockProfileWiFiConfigsFeature) {
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
			useCase, repo, wifiFeat, pwfFeat := profilesTest(t)

			tc.mock(repo, wifiFeat, pwfFeat)

			res, err := useCase.GetCount(context.Background(), "")

			require.Equal(t, tc.res, res)
			require.IsType(t, tc.err, err)
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	ieeeProfileName := "test-8021x-profile-1"

	testProfiles := []entity.Profile{
		{
			ProfileName:            "test-profile-1",
			TenantID:               "tenant-id-456",
			IEEE8021xProfileName:   &ieeeProfileName,
			AuthenticationProtocol: &[]int{1}[0],
			WiredInterface:         &[]bool{true}[0],
		},
		{
			ProfileName: "test-profile-2",
			TenantID:    "tenant-id-456",
		},
	}

	testProfileDTOs := []dto.Profile{
		{
			ProfileName:          "test-profile-1",
			TenantID:             "tenant-id-456",
			Tags:                 []string{""},
			IEEE8021xProfileName: &ieeeProfileName,
			IEEE8021xProfile: &dto.IEEE8021xConfig{
				ProfileName:            ieeeProfileName,
				TenantID:               "tenant-id-456",
				WiredInterface:         true,
				AuthenticationProtocol: 1,
			},
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
			mock: func(repo *mocks.MockProfilesRepository, _ *mocks.MockWiFiConfigsRepository, profileWifiRepo *mocks.MockProfileWiFiConfigsFeature) {
				repo.EXPECT().
					Get(context.Background(), 10, 0, "tenant-id-456").
					Return(testProfiles, nil)
				profileWifiRepo.EXPECT().
					GetByProfileName(context.Background(), "test-profile-1", "tenant-id-456").
					Return([]dto.ProfileWiFiConfigs{}, nil)
				profileWifiRepo.EXPECT().
					GetByProfileName(context.Background(), "test-profile-2", "tenant-id-456").
					Return([]dto.ProfileWiFiConfigs{}, nil)
			},
			res: testProfileDTOs,
			err: nil,
		},
		{
			name:     "database error",
			top:      5,
			skip:     0,
			tenantID: "tenant-id-456",
			mock: func(repo *mocks.MockProfilesRepository, _ *mocks.MockWiFiConfigsRepository, _ *mocks.MockProfileWiFiConfigsFeature) {
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
			mock: func(repo *mocks.MockProfilesRepository, _ *mocks.MockWiFiConfigsRepository, _ *mocks.MockProfileWiFiConfigsFeature) {
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
			useCase, repo, wifiFeat, pwfFeat := profilesTest(t)

			tc.mock(repo, wifiFeat, pwfFeat)

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
			mock: func(repo *mocks.MockProfilesRepository, _ *mocks.MockWiFiConfigsRepository, profilewificonfigfeat *mocks.MockProfileWiFiConfigsFeature) {
				repo.EXPECT().
					GetByName(context.Background(), "test-profile", "tenant-id-456").
					Return(profile, nil)
				profilewificonfigfeat.EXPECT().
					GetByProfileName(context.Background(), "test-profile", "tenant-id-456").
					Return([]dto.ProfileWiFiConfigs{}, nil)
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
			mock: func(repo *mocks.MockProfilesRepository, _ *mocks.MockWiFiConfigsRepository, _ *mocks.MockProfileWiFiConfigsFeature) {
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
			useCase, repo, wifiFeat, pwfFeat := profilesTest(t)

			tc.mock(repo, wifiFeat, pwfFeat)

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
			mock: func(repo *mocks.MockProfilesRepository, _ *mocks.MockWiFiConfigsRepository, profilewificonfigfeat *mocks.MockProfileWiFiConfigsFeature) {
				profilewificonfigfeat.EXPECT().
					DeleteByProfileName(context.Background(), "example-profile", "tenant-id-456").
					Return(nil)
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
			mock: func(repo *mocks.MockProfilesRepository, _ *mocks.MockWiFiConfigsRepository, profilewificonfigfeat *mocks.MockProfileWiFiConfigsFeature) {
				profilewificonfigfeat.EXPECT().
					DeleteByProfileName(context.Background(), "nonexistent-profile", "tenant-id-456").
					Return(nil)
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
			useCase, repo, wifiFeat, pwfFeat := profilesTest(t)

			tc.mock(repo, wifiFeat, pwfFeat)

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
		ProfileName:  "example-profile",
		TenantID:     "tenant-id-456",
		Version:      "1.0.0",
		AMTPassword:  "encrypted",
		MEBXPassword: "encrypted",
	}

	profileDTO := &dto.Profile{
		ProfileName: "example-profile",
		TenantID:    "tenant-id-456",
		Version:     "1.0.0",
		Tags:        []string{""},
		WiFiConfigs: []dto.ProfileWiFiConfigs{
			{
				ProfileName:         "example-profile",
				WirelessProfileName: "wireless-profile-1",
			},
		},
	}

	tests := []test{
		{
			name: "successful update",
			mock: func(repo *mocks.MockProfilesRepository, wifiConfig *mocks.MockWiFiConfigsRepository, profilewificonfigfeat *mocks.MockProfileWiFiConfigsFeature) {
				repo.EXPECT().
					Update(context.Background(), profile).
					Return(true, nil)
				profilewificonfigfeat.EXPECT().
					DeleteByProfileName(context.Background(), profile.ProfileName, profile.TenantID).
					Return(nil)
				repo.EXPECT().
					GetByName(context.Background(), profile.ProfileName, profile.TenantID).
					Return(profile, nil)
				wifiConfig.EXPECT().
					CheckProfileExists(context.Background(), profileDTO.WiFiConfigs[0].WirelessProfileName, profileDTO.TenantID).
					Return(true, nil)
			},
			res: profileDTO,
			err: nil,
		},
		{
			name: "update fails - not found",
			mock: func(repo *mocks.MockProfilesRepository, wifiConfig *mocks.MockWiFiConfigsRepository, _ *mocks.MockProfileWiFiConfigsFeature) {
				repo.EXPECT().
					Update(context.Background(), profile).
					Return(false, profiles.ErrNotFound)
				wifiConfig.EXPECT().
					CheckProfileExists(context.Background(), profileDTO.WiFiConfigs[0].WirelessProfileName, profileDTO.TenantID).
					Return(true, nil)
			},
			res: (*dto.Profile)(nil),
			err: profiles.ErrDatabase,
		},
		{
			name: "update fails - database error",
			mock: func(repo *mocks.MockProfilesRepository, wifiConfig *mocks.MockWiFiConfigsRepository, _ *mocks.MockProfileWiFiConfigsFeature) {
				repo.EXPECT().
					Update(context.Background(), profile).
					Return(false, profiles.ErrDatabase)
				wifiConfig.EXPECT().
					CheckProfileExists(context.Background(), profileDTO.WiFiConfigs[0].WirelessProfileName, profileDTO.TenantID).
					Return(true, nil)
			},
			res: (*dto.Profile)(nil),
			err: profiles.ErrDatabase,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo, wifiFeat, pwfFeat := profilesTest(t)

			tc.mock(repo, wifiFeat, pwfFeat)

			result, err := useCase.Update(context.Background(), profileDTO)

			require.Equal(t, tc.res, result)
			require.IsType(t, err, tc.err)
		})
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	profile := &entity.Profile{
		ProfileName:  "new-profile",
		TenantID:     "tenant-id-789",
		Version:      "1.0.0",
		Tags:         "",
		DHCPEnabled:  true,
		AMTPassword:  "encrypted",
		MEBXPassword: "encrypted",
	}

	profileDTO := &dto.Profile{
		ProfileName: "new-profile",
		TenantID:    "tenant-id-789",
		Version:     "1.0.0",
		Tags:        []string{""},
		DHCPEnabled: true,
		WiFiConfigs: []dto.ProfileWiFiConfigs{
			{
				ProfileName:         "new-profile",
				WirelessProfileName: "wireless-profile-1",
			},
		},
	}

	tests := []test{
		{
			name: "successful insertion",
			mock: func(repo *mocks.MockProfilesRepository, wifiRepo *mocks.MockWiFiConfigsRepository, profilewificonfigfeat *mocks.MockProfileWiFiConfigsFeature) {
				repo.EXPECT().
					Insert(context.Background(), profile).
					Return("unique-profile-id", nil)
				profilewificonfigfeat.EXPECT().
					Insert(context.Background(), &profileDTO.WiFiConfigs[0]).
					Return(nil)
				repo.EXPECT().
					GetByName(context.Background(), profile.ProfileName, profile.TenantID).
					Return(profile, nil)
				wifiRepo.EXPECT().
					CheckProfileExists(context.Background(), profileDTO.WiFiConfigs[0].WirelessProfileName, profileDTO.TenantID).
					Return(true, nil)
			},
			res: profileDTO,
			err: nil,
		},
		{
			name: "insertion fails - database error",
			mock: func(repo *mocks.MockProfilesRepository, wifiRepo *mocks.MockWiFiConfigsRepository, _ *mocks.MockProfileWiFiConfigsFeature) {
				repo.EXPECT().
					Insert(context.Background(), profile).
					Return("", profiles.ErrDatabase)
				wifiRepo.EXPECT().
					CheckProfileExists(context.Background(), profileDTO.WiFiConfigs[0].WirelessProfileName, profileDTO.TenantID).
					Return(true, nil)
			},
			res: (*dto.Profile)(nil),
			err: profiles.ErrDatabase,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo, wifiFeat, pwfFeat := profilesTest(t)

			tc.mock(repo, wifiFeat, pwfFeat)

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
