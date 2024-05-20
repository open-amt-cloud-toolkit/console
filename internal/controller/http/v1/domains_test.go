package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/domains"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

//nolint:gochecknoinits // required to avoid issues when running tests in parallel
func init() {
	gin.SetMode(gin.TestMode)
}

func domainsTest(t *testing.T) (*MockDomainsFeature, *gin.Engine) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	log := logger.New("error")
	domain := NewMockDomainsFeature(mockCtl)

	engine := gin.New()
	handler := engine.Group("/api/v1/admin")

	newDomainRoutes(handler, domain, log)

	return domain, engine
}

type test struct {
	name         string
	method       string
	url          string
	mock         func(repo *MockDomainsFeature)
	response     interface{}
	requestBody  dto.Domain
	expectedCode int
}

var (
	requestDomain  = dto.Domain{ProfileName: "newprofile", TenantID: "tenant1", DomainSuffix: "domain.com", ProvisioningCert: "cert", ProvisioningCertStorageFormat: "string", ProvisioningCertPassword: "password"}
	responseDomain = dto.Domain{ProfileName: "newprofile", TenantID: "tenant1", DomainSuffix: "domain.com", ProvisioningCert: "cert", ProvisioningCertStorageFormat: "string", ProvisioningCertPassword: "password"}
)

func TestDomainRoutes(t *testing.T) {
	t.Parallel()

	tests := []test{
		{
			name:   "get all domains",
			method: http.MethodGet,
			url:    "/api/v1/admin/domains",
			mock: func(domain *MockDomainsFeature) {
				domain.EXPECT().Get(context.Background(), 25, 0, "").Return([]dto.Domain{{
					ProfileName: "profile",
				}}, nil)
			},
			response:     []dto.Domain{{ProfileName: "profile"}},
			expectedCode: http.StatusOK,
		},
		{
			name:   "get all domains - with count",
			method: http.MethodGet,
			url:    "/api/v1/admin/domains?$top=10&$skip=1&$count=true",
			mock: func(domain *MockDomainsFeature) {
				domain.EXPECT().Get(context.Background(), 10, 1, "").Return([]dto.Domain{{
					ProfileName: "profile",
				}}, nil)
				domain.EXPECT().GetCount(context.Background(), "").Return(1, nil)
			},
			response:     DomainCountResponse{Count: 1, Data: []dto.Domain{{ProfileName: "profile"}}},
			expectedCode: http.StatusOK,
		},
		{
			name:   "get all domains - failed",
			method: http.MethodGet,
			url:    "/api/v1/admin/domains",
			mock: func(domain *MockDomainsFeature) {
				domain.EXPECT().Get(context.Background(), 25, 0, "").Return(nil, domains.ErrDatabase)
			},
			response:     domains.ErrDatabase,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "get domain by name",
			method: http.MethodGet,
			url:    "/api/v1/admin/domains/profile",
			mock: func(domain *MockDomainsFeature) {
				domain.EXPECT().GetByName(context.Background(), "profile", "").Return(&dto.Domain{
					ProfileName: "profile",
				}, nil)
			},
			response:     dto.Domain{ProfileName: "profile"},
			expectedCode: http.StatusOK,
		},
		{
			name:   "get domain by name - failed",
			method: http.MethodGet,
			url:    "/api/v1/admin/domains/profile",
			mock: func(domain *MockDomainsFeature) {
				domain.EXPECT().GetByName(context.Background(), "profile", "").Return(nil, domains.ErrDatabase)
			},
			response:     domains.ErrDatabase,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "insert domain",
			method: http.MethodPost,
			url:    "/api/v1/admin/domains",
			mock: func(domain *MockDomainsFeature) {
				domainTest := &dto.Domain{ProfileName: "newprofile", TenantID: "tenant1", DomainSuffix: "domain.com", ProvisioningCert: "cert", ProvisioningCertStorageFormat: "string", ProvisioningCertPassword: "password"}
				domain.EXPECT().Insert(context.Background(), domainTest).Return(domainTest, nil)
			},
			response:     responseDomain,
			requestBody:  requestDomain,
			expectedCode: http.StatusCreated,
		},
		{
			name:   "insert domain - failed",
			method: http.MethodPost,
			url:    "/api/v1/admin/domains",
			mock: func(domain *MockDomainsFeature) {
				domainTest := &dto.Domain{ProfileName: "newprofile", TenantID: "tenant1", DomainSuffix: "domain.com", ProvisioningCert: "cert", ProvisioningCertStorageFormat: "string", ProvisioningCertPassword: "password"}
				domain.EXPECT().Insert(context.Background(), domainTest).Return(nil, domains.ErrDatabase)
			},
			response:     domains.ErrDatabase,
			requestBody:  requestDomain,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "insert domain validation - failed",
			method: http.MethodPost,
			url:    "/api/v1/admin/domains",
			mock: func(domain *MockDomainsFeature) {
				domain400Test := &dto.Domain{ProfileName: "p1", TenantID: "t1", DomainSuffix: "domain1.com", ProvisioningCert: "cert1", ProvisioningCertStorageFormat: "string1"}
				domain.EXPECT().Insert(context.Background(), domain400Test).Return(nil, domains.ErrDatabase)
			},
			response:     domains.ErrDatabase,
			requestBody:  dto.Domain{ProfileName: "p1", TenantID: "t1", DomainSuffix: "domain1.com", ProvisioningCert: "cert1", ProvisioningCertStorageFormat: "string1"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "delete domain",
			method: http.MethodDelete,
			url:    "/api/v1/admin/domains/profile",
			mock: func(domain *MockDomainsFeature) {
				domain.EXPECT().Delete(context.Background(), "profile", "").Return(nil)
			},
			response:     nil,
			expectedCode: http.StatusNoContent,
		},
		{
			name:   "delete domain - failed",
			method: http.MethodDelete,
			url:    "/api/v1/admin/domains/profile",
			mock: func(domain *MockDomainsFeature) {
				domain.EXPECT().Delete(context.Background(), "profile", "").Return(domains.ErrDatabase)
			},
			response:     domains.ErrDatabase,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "update domain",
			method: http.MethodPatch,
			url:    "/api/v1/admin/domains",
			mock: func(domain *MockDomainsFeature) {
				domainTest := &dto.Domain{ProfileName: "newprofile", TenantID: "tenant1", DomainSuffix: "domain.com", ProvisioningCert: "cert", ProvisioningCertStorageFormat: "string", ProvisioningCertPassword: "password"}
				domain.EXPECT().Update(context.Background(), domainTest).Return(domainTest, nil)
			},
			response:     responseDomain,
			requestBody:  requestDomain,
			expectedCode: http.StatusOK,
		},
		{
			name:   "update domain - failed",
			method: http.MethodPatch,
			url:    "/api/v1/admin/domains",
			mock: func(domain *MockDomainsFeature) {
				domainTest := &dto.Domain{ProfileName: "newprofile", TenantID: "tenant1", DomainSuffix: "domain.com", ProvisioningCert: "cert", ProvisioningCertStorageFormat: "string", ProvisioningCertPassword: "password"}
				domain.EXPECT().Update(context.Background(), domainTest).Return(nil, domains.ErrDatabase)
			},
			response:     domains.ErrDatabase,
			requestBody:  requestDomain,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			domainFeature, engine := domainsTest(t)

			tc.mock(domainFeature)

			var req *http.Request

			var err error

			if tc.method == http.MethodPost || tc.method == http.MethodPatch {
				reqBody, _ := json.Marshal(tc.requestBody)
				req, err = http.NewRequest(tc.method, tc.url, bytes.NewBuffer(reqBody))
			} else {
				req, err = http.NewRequest(tc.method, tc.url, http.NoBody)
			}

			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}

			w := httptest.NewRecorder()

			engine.ServeHTTP(w, req)

			require.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedCode == http.StatusOK || tc.expectedCode == http.StatusCreated {
				jsonBytes, _ := json.Marshal(tc.response)
				require.Equal(t, string(jsonBytes), w.Body.String())
			}
		})
	}
}
