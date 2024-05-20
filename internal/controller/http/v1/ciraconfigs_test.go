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
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/ciraconfigs"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

func ciraconfigsTest(t *testing.T) (*MockCIRAConfigsFeature, *gin.Engine) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	log := logger.New("error")
	ciraconfig := NewMockCIRAConfigsFeature(mockCtl)

	engine := gin.New()
	handler := engine.Group("/api/v1/admin")

	newCIRAConfigRoutes(handler, ciraconfig, log)

	return ciraconfig, engine
}

type ciraconfigTest struct {
	name         string
	method       string
	url          string
	mock         func(repo *MockCIRAConfigsFeature)
	response     interface{}
	requestBody  dto.CIRAConfig
	expectedCode int
}

var (
	requestCIRAConfig  = dto.CIRAConfig{ConfigName: "ciraconfig", MPSAddress: "https://example.com", MPSPort: 4433, Username: "username", Password: "password", CommonName: "example.com", ServerAddressFormat: 201, AuthMethod: 2, MPSRootCertificate: "-----BEGIN CERTIFICATE-----\n...", ProxyDetails: "http://example.com", TenantID: "abc123", RegeneratePassword: true, Version: "1.0.0"}
	responseCIRAConfig = dto.CIRAConfig{ConfigName: "ciraconfig", MPSAddress: "https://example.com", MPSPort: 4433, Username: "username", Password: "password", CommonName: "example.com", ServerAddressFormat: 201, AuthMethod: 2, MPSRootCertificate: "-----BEGIN CERTIFICATE-----\n...", ProxyDetails: "http://example.com", TenantID: "abc123", RegeneratePassword: true, Version: "1.0.0"}
)

func TestCIRAConfigRoutes(t *testing.T) {
	t.Parallel()

	tests := []ciraconfigTest{
		{
			name:   "get all ciraconfigs",
			method: http.MethodGet,
			url:    "/api/v1/admin/ciraconfigs",
			mock: func(ciraconfig *MockCIRAConfigsFeature) {
				ciraconfig.EXPECT().Get(context.Background(), 25, 0, "").Return([]dto.CIRAConfig{{
					ConfigName: "config",
				}}, nil)
			},
			response:     []dto.CIRAConfig{{ConfigName: "config"}},
			expectedCode: http.StatusOK,
		},
		{
			name:   "get all ciraconfigs - with count",
			method: http.MethodGet,
			url:    "/api/v1/admin/ciraconfigs?$top=10&$skip=1&$count=true",
			mock: func(ciraconfig *MockCIRAConfigsFeature) {
				ciraconfig.EXPECT().Get(context.Background(), 10, 1, "").Return([]dto.CIRAConfig{{
					ConfigName: "config",
				}}, nil)
				ciraconfig.EXPECT().GetCount(context.Background(), "").Return(1, nil)
			},
			response:     CIRAConfigCountResponse{Count: 1, Data: []dto.CIRAConfig{{ConfigName: "config"}}},
			expectedCode: http.StatusOK,
		},
		{
			name:   "get all ciraconfigs - failed",
			method: http.MethodGet,
			url:    "/api/v1/admin/ciraconfigs",
			mock: func(ciraconfig *MockCIRAConfigsFeature) {
				ciraconfig.EXPECT().Get(context.Background(), 25, 0, "").Return(nil, ciraconfigs.ErrDatabase)
			},
			response:     ciraconfigs.ErrDatabase,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "get ciraconfig by name",
			method: http.MethodGet,
			url:    "/api/v1/admin/ciraconfigs/profile",
			mock: func(ciraconfig *MockCIRAConfigsFeature) {
				ciraconfig.EXPECT().GetByName(context.Background(), "profile", "").Return(&dto.CIRAConfig{
					ConfigName: "config",
				}, nil)
			},
			response:     dto.CIRAConfig{ConfigName: "config"},
			expectedCode: http.StatusOK,
		},
		{
			name:   "get ciraconfig by name - failed",
			method: http.MethodGet,
			url:    "/api/v1/admin/ciraconfigs/profile",
			mock: func(ciraconfig *MockCIRAConfigsFeature) {
				ciraconfig.EXPECT().GetByName(context.Background(), "profile", "").Return(nil, ciraconfigs.ErrDatabase)
			},
			response:     ciraconfigs.ErrDatabase,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "insert ciraconfig",
			method: http.MethodPost,
			url:    "/api/v1/admin/ciraconfigs",
			mock: func(ciraconfig *MockCIRAConfigsFeature) {
				ciraconfigTest := &dto.CIRAConfig{
					ConfigName:          "ciraconfig",
					MPSAddress:          "https://example.com",
					MPSPort:             4433,
					Username:            "username",
					Password:            "password",
					CommonName:          "example.com",
					ServerAddressFormat: 201,
					AuthMethod:          2,
					MPSRootCertificate:  "-----BEGIN CERTIFICATE-----\n...",
					ProxyDetails:        "http://example.com",
					TenantID:            "abc123",
					RegeneratePassword:  true,
					Version:             "1.0.0",
				}
				ciraconfig.EXPECT().Insert(context.Background(), ciraconfigTest).Return(ciraconfigTest, nil)
			},
			response:     responseCIRAConfig,
			requestBody:  requestCIRAConfig,
			expectedCode: http.StatusCreated,
		},
		{
			name:   "insert ciraconfig - failed",
			method: http.MethodPost,
			url:    "/api/v1/admin/ciraconfigs",
			mock: func(ciraconfig *MockCIRAConfigsFeature) {
				ciraconfigTest := &dto.CIRAConfig{
					ConfigName:          "ciraconfig",
					MPSAddress:          "https://example.com",
					MPSPort:             4433,
					Username:            "username",
					Password:            "password",
					CommonName:          "example.com",
					ServerAddressFormat: 201,
					AuthMethod:          2,
					MPSRootCertificate:  "-----BEGIN CERTIFICATE-----\n...",
					ProxyDetails:        "http://example.com",
					TenantID:            "abc123",
					RegeneratePassword:  true,
					Version:             "1.0.0",
				}
				ciraconfig.EXPECT().Insert(context.Background(), ciraconfigTest).Return(nil, ciraconfigs.ErrDatabase)
			},
			response:     ciraconfigs.ErrDatabase,
			requestBody:  requestCIRAConfig,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "insert ciraconfig validation - failed",
			method: http.MethodPost,
			url:    "/api/v1/admin/ciraconfigs",
			mock: func(ciraconfig *MockCIRAConfigsFeature) {
				ciraconfig400Test := &dto.CIRAConfig{
					ConfigName:          "ciraconfig",
					ServerAddressFormat: 201,
					AuthMethod:          2,
					MPSRootCertificate:  "-----BEGIN CERTIFICATE-----\n...",
					ProxyDetails:        "http://example.com",
					TenantID:            "abc123",
					RegeneratePassword:  true,
					Version:             "1.0.0",
				}
				ciraconfig.EXPECT().Insert(context.Background(), ciraconfig400Test).Return(nil, ciraconfigs.ErrDatabase)
			},
			response:     ciraconfigs.ErrDatabase,
			requestBody:  dto.CIRAConfig{ConfigName: "ciraconfig", ServerAddressFormat: 201, AuthMethod: 2, MPSRootCertificate: "-----BEGIN CERTIFICATE-----\n...", ProxyDetails: "http://example.com", TenantID: "abc123", RegeneratePassword: true, Version: "1.0.0"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "delete ciraconfig",
			method: http.MethodDelete,
			url:    "/api/v1/admin/ciraconfigs/profile",
			mock: func(ciraconfig *MockCIRAConfigsFeature) {
				ciraconfig.EXPECT().Delete(context.Background(), "profile", "").Return(nil)
			},
			response:     nil,
			expectedCode: http.StatusNoContent,
		},
		{
			name:   "delete ciraconfig - failed",
			method: http.MethodDelete,
			url:    "/api/v1/admin/ciraconfigs/profile",
			mock: func(ciraconfig *MockCIRAConfigsFeature) {
				ciraconfig.EXPECT().Delete(context.Background(), "profile", "").Return(ciraconfigs.ErrDatabase)
			},
			response:     ciraconfigs.ErrDatabase,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:   "update ciraconfig",
			method: http.MethodPatch,
			url:    "/api/v1/admin/ciraconfigs",
			mock: func(ciraconfig *MockCIRAConfigsFeature) {
				ciraconfigTest := &dto.CIRAConfig{
					ConfigName:          "ciraconfig",
					MPSAddress:          "https://example.com",
					MPSPort:             4433,
					Username:            "username",
					Password:            "password",
					CommonName:          "example.com",
					ServerAddressFormat: 201,
					AuthMethod:          2,
					MPSRootCertificate:  "-----BEGIN CERTIFICATE-----\n...",
					ProxyDetails:        "http://example.com",
					TenantID:            "abc123",
					RegeneratePassword:  true,
					Version:             "1.0.0",
				}
				ciraconfig.EXPECT().Update(context.Background(), ciraconfigTest).Return(ciraconfigTest, nil)
			},
			response:     responseCIRAConfig,
			requestBody:  requestCIRAConfig,
			expectedCode: http.StatusOK,
		},
		{
			name:   "update ciraconfig - failed",
			method: http.MethodPatch,
			url:    "/api/v1/admin/ciraconfigs",
			mock: func(ciraconfig *MockCIRAConfigsFeature) {
				ciraconfigTest := &dto.CIRAConfig{
					ConfigName:          "ciraconfig",
					MPSAddress:          "https://example.com",
					MPSPort:             4433,
					Username:            "username",
					Password:            "password",
					CommonName:          "example.com",
					ServerAddressFormat: 201,
					AuthMethod:          2,
					MPSRootCertificate:  "-----BEGIN CERTIFICATE-----\n...",
					ProxyDetails:        "http://example.com",
					TenantID:            "abc123",
					RegeneratePassword:  true,
					Version:             "1.0.0",
				}
				ciraconfig.EXPECT().Update(context.Background(), ciraconfigTest).Return(nil, ciraconfigs.ErrDatabase)
			},
			response:     ciraconfigs.ErrDatabase,
			requestBody:  requestCIRAConfig,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ciraconfigFeature, engine := ciraconfigsTest(t)

			tc.mock(ciraconfigFeature)

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
