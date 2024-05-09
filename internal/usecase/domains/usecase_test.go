package domains_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/domains"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type test struct {
	name         string
	top          int
	skip         int
	domainName   string
	domainSuffix string
	tenantID     string
	input        entity.Domain
	mock         func(repo *MockRepository)
	res          interface{}
	err          error
}

func domainsTest(t *testing.T) (*domains.UseCase, *MockRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockRepository(mockCtl)
	log := logger.New("error")
	useCase := domains.New(repo, log)

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
				repo.EXPECT().GetCount(context.Background(), "").Return(0, domains.ErrDatabase)
			},
			res: 0,
			err: domains.ErrDatabase,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := domainsTest(t)

			tc.mock(repo)

			res, err := useCase.GetCount(context.Background(), "")

			require.Equal(t, tc.res, res)
			require.IsType(t, tc.err, err)
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	testDomains := []entity.Domain{
		{
			ProfileName: "test-domain-1",
			TenantID:    "tenant-id-456",
		},
		{
			ProfileName: "test-domain-2",
			TenantID:    "tenant-id-456",
		},
	}

	testDomainDTOs := []dto.Domain{
		{
			ProfileName: "test-domain-1",
			TenantID:    "tenant-id-456",
		},
		{
			ProfileName: "test-domain-2",
			TenantID:    "tenant-id-456",
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
					Return(testDomains, nil)
			},
			res: testDomainDTOs,
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
					Return(nil, domains.ErrDatabase)
			},
			res: []dto.Domain(nil),
			err: domains.ErrDatabase,
		},
		{
			name:     "zero results",
			top:      10,
			skip:     20,
			tenantID: "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Get(context.Background(), 10, 20, "tenant-id-456").
					Return([]entity.Domain{}, nil)
			},
			res: []dto.Domain{},
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := domainsTest(t)

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

func TestGetDomainByDomainSuffix(t *testing.T) {
	t.Parallel()

	domain := entity.Domain{
		ProfileName:                   "test-domain",
		DomainSuffix:                  "test.com",
		ProvisioningCert:              "test-cert",
		ProvisioningCertStorageFormat: "test-format",
		ProvisioningCertPassword:      "test-password",
		TenantID:                      "tenant-id-456",
		Version:                       "1.0.0",
	}

	domainDTO := dto.Domain{
		ProfileName:                   "test-domain",
		DomainSuffix:                  "test.com",
		ProvisioningCert:              "test-cert",
		ProvisioningCertStorageFormat: "test-format",
		ProvisioningCertPassword:      "test-password",
		TenantID:                      "tenant-id-456",
		Version:                       "1.0.0",
	}

	tests := []test{
		{
			name:         "successful retrieval",
			domainSuffix: "test.com",
			tenantID:     "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					GetDomainByDomainSuffix(context.Background(), "test.com", "tenant-id-456").
					Return(&domain, nil)
			},
			res: &domainDTO,
			err: nil,
		},
		{
			name:         "domain not found",
			domainSuffix: "unknown.com",
			tenantID:     "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					GetDomainByDomainSuffix(context.Background(), "unknown.com", "tenant-id-456").
					Return(nil, nil)
			},
			res: nil,
			err: domains.ErrNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := domainsTest(t)

			tc.mock(repo)

			domain, err := useCase.GetDomainByDomainSuffix(context.Background(), tc.domainSuffix, tc.tenantID)

			if tc.res != nil {
				require.NotNil(t, domain)
				require.Equal(t, tc.res, domain)
			} else {
				require.Nil(t, domain)
			}

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

	domain := &entity.Domain{
		ProfileName:                   "test-domain",
		DomainSuffix:                  "test-domain",
		ProvisioningCert:              "test-cert",
		ProvisioningCertStorageFormat: "test-format",
		ProvisioningCertPassword:      "test-password",
		TenantID:                      "tenant-id-456",
		Version:                       "1.0.0",
	}

	domainDTO := &dto.Domain{
		ProfileName:                   "test-domain",
		DomainSuffix:                  "test-domain",
		ProvisioningCert:              "test-cert",
		ProvisioningCertStorageFormat: "test-format",
		ProvisioningCertPassword:      "test-password",
		TenantID:                      "tenant-id-456",
		Version:                       "1.0.0",
	}

	tests := []test{
		{
			name: "successful retrieval",
			input: entity.Domain{
				ProfileName: "test-domain",
				TenantID:    "tenant-id-456",
			},
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByName(context.Background(), "test-domain", "tenant-id-456").
					Return(domain, nil)
			},
			res: domainDTO,
			err: nil,
		},
		{
			name: "domain not found",
			input: entity.Domain{
				ProfileName: "unknown-domain",
				TenantID:    "tenant-id-456",
			},
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					GetByName(context.Background(), "unknown-domain", "tenant-id-456").
					Return(nil, nil)
			},
			res: (*dto.Domain)(nil),
			err: domains.ErrNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := domainsTest(t)

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
			name:       "successful deletion",
			domainName: "example-domain",
			tenantID:   "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Delete(context.Background(), "example-domain", "tenant-id-456").
					Return(true, nil)
			},
			err: nil,
		},
		{
			name:       "deletion fails - domain not found",
			domainName: "nonexistent-domain",
			tenantID:   "tenant-id-456",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Delete(context.Background(), "nonexistent-domain", "tenant-id-456").
					Return(false, nil)
			},
			err: domains.ErrNotFound,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := domainsTest(t)

			tc.mock(repo)

			err := useCase.Delete(context.Background(), tc.domainName, tc.tenantID)

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

	domain := &entity.Domain{
		ProfileName: "example-domain",
		TenantID:    "tenant-id-456",
		Version:     "1.0.0",
	}
	domainDTO := &dto.Domain{
		ProfileName: "example-domain",
		TenantID:    "tenant-id-456",
		Version:     "1.0.0",
	}

	tests := []test{
		{
			name: "successful update",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Update(context.Background(), domain).
					Return(true, nil)
				repo.EXPECT().
					GetByName(context.Background(), domain.ProfileName, domain.TenantID).
					Return(domain, nil)
			},
			res: domainDTO,
			err: nil,
		},
		{
			name: "update fails - database error",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Update(context.Background(), domain).
					Return(false, domains.ErrDatabase)
			},
			res: (*dto.Domain)(nil),
			err: domains.ErrDatabase,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			useCase, repo := domainsTest(t)

			tc.mock(repo)

			result, err := useCase.Update(context.Background(), domainDTO)

			require.Equal(t, tc.res, result)
			require.IsType(t, tc.err, err)
		})
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()

	domain := &entity.Domain{
		ProfileName:                   "new-domain",
		DomainSuffix:                  "newdomain.com",
		ProvisioningCert:              "cert-data",
		ProvisioningCertStorageFormat: "PEM",
		ProvisioningCertPassword:      "password",
		TenantID:                      "tenant-id-789",
		Version:                       "1.0.0",
	}
	domainDTO := &dto.Domain{
		ProfileName:                   "new-domain",
		DomainSuffix:                  "newdomain.com",
		ProvisioningCert:              "cert-data",
		ProvisioningCertStorageFormat: "PEM",
		ProvisioningCertPassword:      "password",
		TenantID:                      "tenant-id-789",
		Version:                       "1.0.0",
	}

	tests := []test{
		{
			name: "successful insertion",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Insert(context.Background(), domain).
					Return("unique-domain-id", nil)
				repo.EXPECT().
					GetByName(context.Background(), domain.ProfileName, domain.TenantID).
					Return(domain, nil)
			},
			res: domainDTO,
			err: nil,
		},
		{
			name: "insertion fails - database error",
			mock: func(repo *MockRepository) {
				repo.EXPECT().
					Insert(context.Background(), domain).
					Return("", domains.ErrDatabase)
			},
			res: (*dto.Domain)(nil),
			err: domains.ErrDatabase,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			useCase, repo := domainsTest(t)

			tc.mock(repo)

			id, err := useCase.Insert(context.Background(), domainDTO)

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
