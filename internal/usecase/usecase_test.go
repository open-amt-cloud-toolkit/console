package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/open-amt-cloud-toolkit/console/internal/usecase/ciraconfigs"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/domains"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/ieee8021xconfigs"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/profiles"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/profilewificonfigs"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/sqldb"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/wificonfigs"
	"github.com/open-amt-cloud-toolkit/console/pkg/db"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type usecaseTest struct {
	name           string
	initializeFunc func() *Usecases
	expectedResult *Usecases
}

func TestUsecases(t *testing.T) {
	t.Parallel()

	tests := []usecaseTest{
		{
			name: "NewUseCases initializes correctly",
			initializeFunc: func() *Usecases {
				mockDB := NewMockDB()
				mockLogger := &MockLogger{}

				return NewUseCases(mockDB, mockLogger)
			},
			expectedResult: &Usecases{
				Domains:            domains.New(sqldb.NewDomainRepo(&db.SQL{}, &MockLogger{}), &MockLogger{}),
				Devices:            devices.New(sqldb.NewDeviceRepo(&db.SQL{}, &MockLogger{}), wsman.NewGoWSMANMessages(), devices.NewRedirector(), wsman.NewGoWSMANMessages(), &MockLogger{}),
				Profiles:           profiles.New(sqldb.NewProfileRepo(&db.SQL{}, &MockLogger{}), wificonfigs.New(sqldb.NewWirelessRepo(&db.SQL{}, &MockLogger{}), ieee8021xconfigs.New(sqldb.NewIEEE8021xRepo(&db.SQL{}, &MockLogger{}), &MockLogger{}), &MockLogger{}), profilewificonfigs.New(sqldb.NewProfileWiFiConfigsRepo(&db.SQL{}, &MockLogger{}), &MockLogger{}), ieee8021xconfigs.New(sqldb.NewIEEE8021xRepo(&db.SQL{}, &MockLogger{}), &MockLogger{}), &MockLogger{}),
				IEEE8021xProfiles:  ieee8021xconfigs.New(sqldb.NewIEEE8021xRepo(&db.SQL{}, &MockLogger{}), &MockLogger{}),
				CIRAConfigs:        ciraconfigs.New(sqldb.NewCIRARepo(&db.SQL{}, &MockLogger{}), &MockLogger{}),
				WirelessProfiles:   wificonfigs.New(sqldb.NewWirelessRepo(&db.SQL{}, &MockLogger{}), ieee8021xconfigs.New(sqldb.NewIEEE8021xRepo(&db.SQL{}, &MockLogger{}), &MockLogger{}), &MockLogger{}),
				ProfileWiFiConfigs: profilewificonfigs.New(sqldb.NewProfileWiFiConfigsRepo(&db.SQL{}, &MockLogger{}), &MockLogger{}),
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := tc.initializeFunc()

			require.NotNil(t, uc)
			assert.NotNil(t, uc.Devices)
			assert.NotNil(t, uc.Domains)
			assert.NotNil(t, uc.Profiles)
			assert.NotNil(t, uc.ProfileWiFiConfigs)
			assert.NotNil(t, uc.IEEE8021xProfiles)
			assert.NotNil(t, uc.CIRAConfigs)
			assert.NotNil(t, uc.WirelessProfiles)

			assert.Equal(t, tc.expectedResult.Domains, uc.Domains)
			assert.Equal(t, tc.expectedResult.Devices, uc.Devices)
			assert.Equal(t, tc.expectedResult.Profiles, uc.Profiles)
			assert.Equal(t, tc.expectedResult.ProfileWiFiConfigs, uc.ProfileWiFiConfigs)
			assert.Equal(t, tc.expectedResult.IEEE8021xProfiles, uc.IEEE8021xProfiles)
			assert.Equal(t, tc.expectedResult.CIRAConfigs, uc.CIRAConfigs)
			assert.Equal(t, tc.expectedResult.WirelessProfiles, uc.WirelessProfiles)
		})
	}
}

type MockDB struct {
	*db.SQL
}

func NewMockDB() *db.SQL {
	return &db.SQL{}
}

type MockLogger struct {
	logger.Interface
}

func TestInitialization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level         string
		expectedLevel string
	}{
		{"debug", "debug"},
		{"info", "info"},
		{"warn", "warn"},
		{"error", "error"},
		{"invalid", "info"},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.level, func(t *testing.T) {
			t.Parallel()

			mockDB := NewMockDB()
			mockLogger := &MockLogger{}

			uc := NewUseCases(mockDB, mockLogger)

			require.NotNil(t, uc)
			assert.NotNil(t, uc.Devices)
		})
	}
}
