//nolint:gci // ignore import order
package sqldb_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/sqldb"
	"github.com/open-amt-cloud-toolkit/console/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestProfileRepo_GetCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func(dbConn *sql.DB)
		tenantID string
		expected int
		err      error
	}{
		{
			name: "Successful count",
			setup: func(dbConn *sql.DB) {
				_, err := dbConn.Exec(`INSERT INTO profiles (profile_name, tenant_id) VALUES (?, ?)`,
					"profile1", "tenant1")
				require.NoError(t, err)
			},
			tenantID: "tenant1",
			expected: 1,
			err:      nil,
		},
		{
			name:     "No profiles found",
			setup:    func(_ *sql.DB) {},
			tenantID: "tenant2",
			expected: 0,
			err:      nil,
		},
		{
			name:     "Query execution error",
			setup:    func(_ *sql.DB) {},
			tenantID: "tenant1",
			expected: 0,
			err:      &sqldb.DatabaseError{},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dbConn, err := sql.Open("sqlite", ":memory:")
			require.NoError(t, err)
			defer dbConn.Close()

			_, err = dbConn.Exec(`
				CREATE TABLE profiles (
					profile_name TEXT PRIMARY KEY,
					tenant_id TEXT NOT NULL
				);
			`)
			require.NoError(t, err)

			tc.setup(dbConn)

			sqlConfig := &db.SQL{
				Builder:    squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
				Pool:       dbConn,
				IsEmbedded: true,
			}

			if tc.name == QueryExecutionErrorTestName {
				sqlConfig.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.AtP)
			}

			mockLog := new(MockLogger)
			repo := sqldb.NewProfileRepo(sqlConfig, mockLog)

			count, err := repo.GetCount(context.Background(), tc.tenantID)

			if err == nil && tc.err != nil {
				t.Errorf("Expected error of type %T, got nil", tc.err)
			} else if err != nil {
				var dbError sqldb.DatabaseError

				if !errors.As(err, &dbError) {
					t.Errorf("Expected error of type %T, got %T", tc.err, err)
				}
			}

			if count != tc.expected {
				t.Errorf("Expected count %d, got %d", tc.expected, count)
			}
		})
	}
}

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func BoolPtr(b bool) *bool {
	return &b
}

func GetProfileHelper(t *testing.T, tc GetProfileTestCase, profiles []entity.Profile, err error) {
	t.Helper()

	if (err != nil) != (tc.err != nil) {
		t.Errorf("Expected error: %v, got: %v", tc.err, err)
	}

	if len(profiles) == 0 && len(tc.expected) > 0 {
		t.Errorf("Expected %d profiles, got %d", len(tc.expected), len(profiles))

		return
	}

	if len(profiles) != len(tc.expected) {
		t.Errorf("Expected %d profiles, got %d", len(tc.expected), len(profiles))
	}

	for i := range tc.expected {
		expectedProfile := &tc.expected[i]

		if i >= len(profiles) {
			t.Errorf("Expected profile %d, but got none", i)

			break
		}

		actualProfile := &profiles[i]
		assert.IsType(t, expectedProfile, actualProfile)
	}
}

type GetProfileTestCase struct {
	name     string
	setup    func(dbConn *sql.DB)
	top      int
	skip     int
	tenantID string
	expected []entity.Profile
	err      error
}

func TestProfileRepo_Get(t *testing.T) {
	t.Parallel()

	tests := []GetProfileTestCase{
		{
			name: "Successful query",
			setup: func(dbConn *sql.DB) {
				_, err := dbConn.Exec(`INSERT INTO profiles (
					profile_name, amt_password, creation_date, created_by, generate_random_password,
					cira_config_name, activation, mebx_password, generate_random_mebx_password, tags,
					dhcp_enabled, ip_sync_enabled, local_wifi_sync_enabled, tenant_id, tls_mode, 
					tls_signing_authority, user_consent, ider_enabled, kvm_enabled, sol_enabled, 
					ieee8021x_profile_name, version, auth_protocol, server_name, domain,
					username, password, roaming_identity, active_in_s0, pxe_timeout, wired_interface
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					"profile1", "password1", "2024-08-01", "user1", true,
					"cira1", "activation1", "mebx1", true, "tags1",
					true, true, true, "tenant1", 1,
					"authority1", "consent1", true, true, true,
					"ieee1", "v1", 1, "server1", "domain1",
					"user1", "pass1", "identity1", true, 10, true)
				require.NoError(t, err)

				_, err = dbConn.Exec(`
                CREATE TABLE profiles_wirelessconfigs (
                    profile_name TEXT NOT NULL,
                    authentication_method INTEGER,
                    encryption_method INTEGER,
                    ssid TEXT,
                    psk_value INTEGER,
                    psk_passphrase TEXT,
                    link_policy TEXT,
                    tenant_id TEXT NOT NULL,
                    ieee8021x_profile_name TEXT,
                    version TEXT,
                    authentication_protocol INTEGER,
                    pxe_timeout INTEGER,
                    wired_interface BOOLEAN,
                    PRIMARY KEY (profile_name, tenant_id)
                );`)
				require.NoError(t, err)

				_, err = dbConn.Exec(`
                INSERT INTO profiles_wirelessconfigs (
                    profile_name, authentication_method, encryption_method, ssid, 
                    psk_value, psk_passphrase, link_policy, tenant_id, 
                    ieee8021x_profile_name, version, authentication_protocol, 
                    pxe_timeout, wired_interface
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
					"profile1", 1, 1, "TestSSID", 12345, "TestPassphrase", "TestLinkPolicy", "tenant1",
					"ieee1", "v1", 1, 30, true)
				require.NoError(t, err)

				_, err = dbConn.Exec(`
                CREATE TABLE ieee8021xconfigs (
                    profile_name TEXT NOT NULL,
                    authentication_protocol INTEGER,
                    pxe_timeout INTEGER,
                    wired_interface BOOLEAN,
                    tenant_id TEXT NOT NULL,
                    version TEXT,
                    PRIMARY KEY (profile_name, tenant_id)
                );`)
				require.NoError(t, err)

				_, err = dbConn.Exec(`
                INSERT INTO ieee8021xconfigs (
                    profile_name, authentication_protocol, pxe_timeout, wired_interface, tenant_id, version
                ) VALUES (?, ?, ?, ?, ?, ?);`,
					"profile1", 1, 30, true, "tenant1", "v1")
				require.NoError(t, err)
			},
			top:      10,
			skip:     0,
			tenantID: "tenant1",
			expected: []entity.Profile{
				{
					ProfileName:                "profile1",
					AMTPassword:                "password1",
					CreationDate:               "2024-08-01",
					CreatedBy:                  "user1",
					GenerateRandomPassword:     true,
					CIRAConfigName:             StringPtr("cira1"),
					Activation:                 "activation1",
					MEBXPassword:               "mebx1",
					GenerateRandomMEBxPassword: true,
					Tags:                       "tags1",
					DHCPEnabled:                true,
					IPSyncEnabled:              true,
					LocalWiFiSyncEnabled:       true,
					TenantID:                   "tenant1",
					TLSMode:                    1,
					TLSSigningAuthority:        "authority1",
					UserConsent:                "consent1",
					IDEREnabled:                true,
					KVMEnabled:                 true,
					SOLEnabled:                 true,
					IEEE8021xProfileName:       StringPtr("ieee1"),
					Version:                    "v1",
					AuthenticationProtocol:     IntPtr(1),
					ServerName:                 "server1",
					Domain:                     "domain1",
					Username:                   "user1",
					Password:                   "pass1",
					RoamingIdentity:            "identity1",
					ActiveInS0:                 true,
					PXETimeout:                 IntPtr(10),
					WiredInterface:             BoolPtr(true),
				},
			},
			err: nil,
		},
		{
			name: "No profiles found",
			setup: func(dbConn *sql.DB) {
				_, err := dbConn.Exec(`
                CREATE TABLE profiles_wirelessconfigs (
                    profile_name TEXT NOT NULL,
                    authentication_method INTEGER,
                    encryption_method INTEGER,
                    ssid TEXT,
                    psk_value INTEGER,
                    psk_passphrase TEXT,
                    link_policy TEXT,
                    tenant_id TEXT NOT NULL,
                    ieee8021x_profile_name TEXT,
                    version TEXT,
                    authentication_protocol INTEGER,
                    pxe_timeout INTEGER,
                    wired_interface BOOLEAN,
                    PRIMARY KEY (profile_name, tenant_id)
                );`)
				require.NoError(t, err)

				_, err = dbConn.Exec(`
                INSERT INTO profiles_wirelessconfigs (
                    profile_name, authentication_method, encryption_method, ssid, 
                    psk_value, psk_passphrase, link_policy, tenant_id, 
                    ieee8021x_profile_name, version, authentication_protocol, 
                    pxe_timeout, wired_interface
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
					"profile1", 1, 1, "TestSSID", 12345, "TestPassphrase", "TestLinkPolicy", "tenant1",
					"ieee1", "v1", 1, 30, true)
				require.NoError(t, err)

				_, err = dbConn.Exec(`
                CREATE TABLE ieee8021xconfigs (
                    profile_name TEXT NOT NULL,
                    authentication_protocol INTEGER,
                    pxe_timeout INTEGER,
                    wired_interface BOOLEAN,
                    tenant_id TEXT NOT NULL,
                    version TEXT,
                    PRIMARY KEY (profile_name, tenant_id)
                );`)
				require.NoError(t, err)

				_, err = dbConn.Exec(`
                INSERT INTO ieee8021xconfigs (
                    profile_name, authentication_protocol, pxe_timeout, wired_interface, tenant_id, version
                ) VALUES (?, ?, ?, ?, ?, ?);`,
					"profile1", 1, 30, true, "tenant1", "v1")
				require.NoError(t, err)
			},
			top:      10,
			skip:     0,
			tenantID: "tenant2",
			expected: []entity.Profile{},
			err:      nil,
		},
		{
			name:     "Query execution error",
			setup:    func(_ *sql.DB) {},
			top:      0,
			skip:     0,
			tenantID: "tenant1",
			expected: nil,
			err:      ErrGeneral,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dbConn, err := sql.Open("sqlite", ":memory:")
			require.NoError(t, err)
			defer dbConn.Close()

			_, err = dbConn.Exec(`
				CREATE TABLE profiles (
					profile_name TEXT NOT NULL,
					amt_password TEXT,
					creation_date TEXT,
					created_by TEXT,
					generate_random_password BOOLEAN,
					cira_config_name TEXT,
					activation TEXT,
					mebx_password TEXT,
					generate_random_mebx_password BOOLEAN,
					tags TEXT,
					dhcp_enabled BOOLEAN,
					ip_sync_enabled BOOLEAN,
					local_wifi_sync_enabled BOOLEAN,
					tenant_id TEXT NOT NULL,
					tls_mode INTEGER,
					tls_signing_authority TEXT,
					user_consent TEXT,
					ider_enabled BOOLEAN,
					kvm_enabled BOOLEAN,
					sol_enabled BOOLEAN,
					ieee8021x_profile_name TEXT,
					version TEXT,
					auth_protocol INTEGER,
					server_name TEXT,
					domain TEXT,
					username TEXT,
					password TEXT,
					roaming_identity TEXT,
					active_in_s0 BOOLEAN,
					pxe_timeout INTEGER,
					wired_interface BOOLEAN
				);
			`)
			require.NoError(t, err)

			tc.setup(dbConn)

			sqlConfig := &db.SQL{
				Builder:    squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
				Pool:       dbConn,
				IsEmbedded: true,
			}

			mockLog := new(MockLogger)
			repo := sqldb.NewProfileRepo(sqlConfig, mockLog)

			profiles, err := repo.Get(context.Background(), tc.top, tc.skip, tc.tenantID)

			GetProfileHelper(t, tc, profiles, err)
		})
	}
}

func TestProfileRepo_GetByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(dbConn *sql.DB)
		profileName string
		tenantID    string
		expected    *entity.Profile
		expectError bool
	}{
		{
			name: "Successful retrieval",
			setup: func(dbConn *sql.DB) {
				_, err := dbConn.Exec(`INSERT INTO profiles (
					profile_name, amt_password, creation_date, created_by, generate_random_password,
					cira_config_name, activation, mebx_password, generate_random_mebx_password, tags,
					dhcp_enabled, ip_sync_enabled, local_wifi_sync_enabled, tenant_id, tls_mode,
					tls_signing_authority, user_consent, ider_enabled, kvm_enabled, sol_enabled,
					ieee8021x_profile_name, version, auth_protocol, server_name, domain,
					username, password, roaming_identity, active_in_s0, pxe_timeout, wired_interface
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					"profile1", "password1", "2024-08-01", "user1", true,
					"cira1", "activation1", "mebx1", true, "tags1",
					true, true, true, "tenant1", 1,
					"authority1", "consent1", true, true, true,
					"ieee1", "v1", 1, "server1", "domain1",
					"user1", "pass1", "identity1", true, 10, true)
				require.NoError(t, err)

				_, err = dbConn.Exec(`
		        CREATE TABLE ieee8021xconfigs (
		            profile_name TEXT NOT NULL,
		            authentication_protocol INTEGER,
		            pxe_timeout INTEGER,
		            wired_interface BOOLEAN,
		            tenant_id TEXT NOT NULL,
		            version TEXT,
		            PRIMARY KEY (profile_name, tenant_id)
		        );`)
				require.NoError(t, err)

				_, err = dbConn.Exec(`
		        INSERT INTO ieee8021xconfigs (
		            profile_name, authentication_protocol, pxe_timeout, wired_interface, tenant_id, version
		        ) VALUES (?, ?, ?, ?, ?, ?);`,
					"profile1", 1, 30, true, "tenant1", "v1")
				require.NoError(t, err)
			},
			profileName: "profile1",
			tenantID:    "tenant1",
			expected:    &entity.Profile{},
			expectError: false,
		},
		{
			name: "No Profile Found",
			setup: func(dbConn *sql.DB) {
				_, err := dbConn.Exec(`
		        CREATE TABLE ieee8021xconfigs (
		            profile_name TEXT NOT NULL,
		            authentication_protocol INTEGER,
		            pxe_timeout INTEGER,
		            wired_interface BOOLEAN,
		            tenant_id TEXT NOT NULL,
		            version TEXT,
		            PRIMARY KEY (profile_name, tenant_id)
		        );`)
				require.NoError(t, err)

				_, err = dbConn.Exec(`
		        INSERT INTO ieee8021xconfigs (
		            profile_name, authentication_protocol, pxe_timeout, wired_interface, tenant_id, version
		        ) VALUES (?, ?, ?, ?, ?, ?);`,
					"profile1", 1, 30, true, "tenant1", "v1")
				require.NoError(t, err)
			},
			profileName: "dontexist",
			tenantID:    "tenant1",
			expected: &entity.Profile{
				ProfileName:                "profile1",
				AMTPassword:                "password1",
				CreationDate:               "2024-08-01",
				CreatedBy:                  "user1",
				GenerateRandomPassword:     true,
				CIRAConfigName:             StringPtr("cira1"),
				Activation:                 "activation1",
				MEBXPassword:               "mebx1",
				GenerateRandomMEBxPassword: true,
				Tags:                       "tags1",
				DHCPEnabled:                true,
				IPSyncEnabled:              true,
				LocalWiFiSyncEnabled:       true,
				TenantID:                   "tenant1",
				TLSMode:                    1,
				TLSSigningAuthority:        "authority1",
				UserConsent:                "consent1",
				IDEREnabled:                true,
				KVMEnabled:                 true,
				SOLEnabled:                 true,
				IEEE8021xProfileName:       StringPtr("ieee1"),
				Version:                    "v1",
				AuthenticationProtocol:     IntPtr(1),
				ServerName:                 "server1",
				Domain:                     "domain1",
				Username:                   "user1",
				Password:                   "pass1",
				RoamingIdentity:            "identity1",
				ActiveInS0:                 true,
				PXETimeout:                 IntPtr(10),
				WiredInterface:             BoolPtr(true),
			},
			expectError: false,
		},
		{
			name:        "Query execution error",
			setup:       func(_ *sql.DB) {},
			profileName: "b",
			tenantID:    "tenant1",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dbConn, err := sql.Open("sqlite", ":memory:")
			require.NoError(t, err)
			defer dbConn.Close()

			_, err = dbConn.Exec(`
				CREATE TABLE profiles (
					profile_name TEXT NOT NULL,
					amt_password TEXT,
					creation_date TEXT,
					created_by TEXT,
					generate_random_password BOOLEAN,
					cira_config_name TEXT,
					activation TEXT,
					mebx_password TEXT,
					generate_random_mebx_password BOOLEAN,
					tags TEXT,
					dhcp_enabled BOOLEAN,
					ip_sync_enabled BOOLEAN,
					local_wifi_sync_enabled BOOLEAN,
					tenant_id TEXT NOT NULL,
					tls_mode INTEGER,
					tls_signing_authority TEXT,
					user_consent TEXT,
					ider_enabled BOOLEAN,
					kvm_enabled BOOLEAN,
					sol_enabled BOOLEAN,
					ieee8021x_profile_name TEXT,
					version TEXT,
					auth_protocol INTEGER,
					server_name TEXT,
					domain TEXT,
					username TEXT,
					password TEXT,
					roaming_identity TEXT,
					active_in_s0 BOOLEAN,
					pxe_timeout INTEGER,
					wired_interface BOOLEAN
				);
			`)
			require.NoError(t, err)

			tc.setup(dbConn)

			sqlConfig := &db.SQL{
				Builder:    squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
				Pool:       dbConn,
				IsEmbedded: true,
			}

			repo := sqldb.NewProfileRepo(sqlConfig, new(MockLogger))

			profile, err := repo.GetByName(context.Background(), tc.profileName, tc.tenantID)

			if (err != nil) != tc.expectError {
				t.Errorf("Expected error status %v, got %v", tc.expectError, err != nil)
			}

			if profile == nil && tc.expected == nil {
				return
			}

			assert.IsType(t, tc.expected, profile)
			assert.IsType(t, tc.expected, profile)
		})
	}
}

func TestProfileRepo_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(dbConn *sql.DB)
		profileName string
		tenantID    string
		expected    bool
		expectError bool
	}{
		{
			name: "Successful delete",
			setup: func(dbConn *sql.DB) {
				_, err := dbConn.Exec(`INSERT INTO profiles (profile_name, activation, generate_random_password, cira_config_name, generate_random_mebx_password, tags, dhcp_enabled, tenant_id, tls_mode, user_consent, ider_enabled, kvm_enabled, sol_enabled, tls_signing_authority, ip_sync_enabled, local_wifi_sync_enabled, ieee8021x_profile_name, auth_protocol, pxe_timeout, wired_interface) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					"profile1", true, true, "cira1", true, "tag1", true, "tenant1", "tls", true, true, true, true, "authority", true, true, "ieee1", 1, 30, true)
				require.NoError(t, err)
			},
			profileName: "profile1",
			tenantID:    "tenant1",
			expected:    true,
			expectError: false,
		},
		{
			name:        "No matching profile",
			setup:       func(_ *sql.DB) {},
			profileName: "nonexistent-profile",
			tenantID:    "tenant2",
			expected:    false,
			expectError: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dbConn, err := sql.Open("sqlite", ":memory:")
			require.NoError(t, err)
			defer dbConn.Close()

			_, err = dbConn.Exec(`
				CREATE TABLE profiles (
					profile_name TEXT NOT NULL,
					activation BOOLEAN,
					generate_random_password BOOLEAN,
					cira_config_name TEXT,
					generate_random_mebx_password BOOLEAN,
					tags TEXT,
					dhcp_enabled BOOLEAN,
					tenant_id TEXT NOT NULL,
					tls_mode TEXT,
					user_consent BOOLEAN,
					ider_enabled BOOLEAN,
					kvm_enabled BOOLEAN,
					sol_enabled BOOLEAN,
					tls_signing_authority TEXT,
					ip_sync_enabled BOOLEAN,
					local_wifi_sync_enabled BOOLEAN,
					ieee8021x_profile_name TEXT,
					auth_protocol INTEGER,
					pxe_timeout INTEGER,
					wired_interface BOOLEAN
				);
			`)
			require.NoError(t, err)

			tc.setup(dbConn)

			sqlConfig := &db.SQL{
				Builder:    squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
				Pool:       dbConn,
				IsEmbedded: true,
			}

			mockLog := new(MockLogger)
			repo := sqldb.NewProfileRepo(sqlConfig, mockLog)

			deleted, err := repo.Delete(context.Background(), tc.profileName, tc.tenantID)

			if (err != nil) != tc.expectError {
				t.Errorf("Expected error status %v, got %v", tc.expectError, err != nil)
			}

			if deleted != tc.expected {
				t.Errorf("Expected deleted status %v, got %v", tc.expected, deleted)
			}
		})
	}
}

func TestProfileRepo_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(dbConn *sql.DB)
		profile     *entity.Profile
		expected    bool
		expectError bool
	}{
		{
			name: "Successful update",
			setup: func(dbConn *sql.DB) {
				_, err := dbConn.Exec(`INSERT INTO profiles (
					profile_name, amt_password, creation_date, created_by, generate_random_password,
					cira_config_name, activation, mebx_password, generate_random_mebx_password, tags,
					dhcp_enabled, ip_sync_enabled, local_wifi_sync_enabled, tenant_id, tls_mode,
					tls_signing_authority, user_consent, ider_enabled, kvm_enabled, sol_enabled,
					ieee8021x_profile_name, version, auth_protocol, server_name, domain,
					username, password, roaming_identity, active_in_s0, pxe_timeout, wired_interface
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					"profile1", "password1", "2024-08-01", "user1", true,
					"cira1", "activation1", "mebx1", true, "tags1",
					true, true, true, "tenant1", 1,
					"authority1", "consent1", true, true, true,
					"ieee1", "v1", 1, "server1", "domain1",
					"user1", "pass1", "identity1", true, 10, true)
				require.NoError(t, err)
			},
			profile: &entity.Profile{
				ProfileName:                "profile1",
				AMTPassword:                "new-password",
				GenerateRandomPassword:     false,
				CIRAConfigName:             StringPtr("new-cira"),
				Activation:                 "new-activation",
				MEBXPassword:               "new-mebx",
				GenerateRandomMEBxPassword: false,
				Tags:                       "new-tags",
				DHCPEnabled:                false,
				IPSyncEnabled:              false,
				LocalWiFiSyncEnabled:       false,
				TenantID:                   "tenant1",
				TLSMode:                    2,
				TLSSigningAuthority:        "new-authority",
				UserConsent:                "new-consent",
				IDEREnabled:                false,
				KVMEnabled:                 false,
				SOLEnabled:                 false,
				IEEE8021xProfileName:       StringPtr("new-ieee"),
				Version:                    "v2",
				AuthenticationProtocol:     IntPtr(2),
				ServerName:                 "new-server",
				Domain:                     "new-domain",
				Username:                   "new-user",
				Password:                   "new-pass",
				RoamingIdentity:            "new-identity",
				ActiveInS0:                 false,
				PXETimeout:                 IntPtr(20),
				WiredInterface:             BoolPtr(false),
			},
			expected:    true,
			expectError: false,
		},
		{
			name:  "Update non-existent profile",
			setup: func(_ *sql.DB) {},
			profile: &entity.Profile{
				ProfileName:                "nonexistent-profile",
				AMTPassword:                "password",
				GenerateRandomPassword:     true,
				CIRAConfigName:             StringPtr("cira"),
				Activation:                 "activation",
				MEBXPassword:               "mebx",
				GenerateRandomMEBxPassword: true,
				Tags:                       "tags",
				DHCPEnabled:                true,
				IPSyncEnabled:              true,
				LocalWiFiSyncEnabled:       true,
				TenantID:                   "tenant1",
				TLSMode:                    1,
				TLSSigningAuthority:        "authority",
				UserConsent:                "consent",
				IDEREnabled:                true,
				KVMEnabled:                 true,
				SOLEnabled:                 true,
				IEEE8021xProfileName:       StringPtr("ieee"),
				Version:                    "v1",
				AuthenticationProtocol:     IntPtr(1),
				ServerName:                 "server",
				Domain:                     "domain",
				Username:                   "user",
				Password:                   "pass",
				RoamingIdentity:            "identity",
				ActiveInS0:                 true,
				PXETimeout:                 IntPtr(10),
				WiredInterface:             BoolPtr(true),
			},
			expected:    false,
			expectError: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dbConn, err := sql.Open("sqlite", ":memory:")
			require.NoError(t, err)
			defer dbConn.Close()

			_, err = dbConn.Exec(`
				CREATE TABLE profiles (
					profile_name TEXT NOT NULL,
					amt_password TEXT,
					creation_date TEXT,
					created_by TEXT,
					generate_random_password BOOLEAN,
					cira_config_name TEXT,
					activation TEXT,
					mebx_password TEXT,
					generate_random_mebx_password BOOLEAN,
					tags TEXT,
					dhcp_enabled BOOLEAN,
					ip_sync_enabled BOOLEAN,
					local_wifi_sync_enabled BOOLEAN,
					tenant_id TEXT NOT NULL,
					tls_mode INTEGER,
					tls_signing_authority TEXT,
					user_consent TEXT,
					ider_enabled BOOLEAN,
					kvm_enabled BOOLEAN,
					sol_enabled BOOLEAN,
					ieee8021x_profile_name TEXT,
					version TEXT,
					auth_protocol INTEGER,
					server_name TEXT,
					domain TEXT,
					username TEXT,
					password TEXT,
					roaming_identity TEXT,
					active_in_s0 BOOLEAN,
					pxe_timeout INTEGER,
					wired_interface BOOLEAN,
					PRIMARY KEY (profile_name, tenant_id)
				);
			`)
			require.NoError(t, err)

			tc.setup(dbConn)

			sqlConfig := &db.SQL{
				Builder:    squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
				Pool:       dbConn,
				IsEmbedded: true,
			}

			repo := sqldb.NewProfileRepo(sqlConfig, new(MockLogger))

			updated, err := repo.Update(context.Background(), tc.profile)

			if (err != nil) != tc.expectError {
				t.Errorf("Expected error status %v, got %v", tc.expectError, err != nil)
			}

			if updated != tc.expected {
				t.Errorf("Expected update status %v, got %v", tc.expected, updated)
			}
		})
	}
}

func TestProfileRepo_Insert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(dbConn *sql.DB)
		profile     *entity.Profile
		expectedErr bool
	}{
		{
			name: "Successful insertion",
			setup: func(dbConn *sql.DB) {
				_, err := dbConn.Exec(`
				CREATE TABLE profiles (
					profile_name TEXT NOT NULL PRIMARY KEY,
					activation TEXT,
					amt_password TEXT,
					generate_random_password BOOLEAN,
					cira_config_name TEXT,
					mebx_password TEXT,
					generate_random_mebx_password BOOLEAN,
					tags TEXT,
					dhcp_enabled BOOLEAN,
					tls_mode INTEGER,
					user_consent TEXT,
					ider_enabled BOOLEAN,
					kvm_enabled BOOLEAN,
					sol_enabled BOOLEAN,
					tls_signing_authority TEXT,
					ieee8021x_profile_name TEXT,
					ip_sync_enabled BOOLEAN,
					local_wifi_sync_enabled BOOLEAN,
					tenant_id TEXT NOT NULL
				);
				`)
				require.NoError(t, err)
			},
			profile: &entity.Profile{
				ProfileName:                "profile1",
				Activation:                 "activation1",
				AMTPassword:                "password1",
				GenerateRandomPassword:     true,
				CIRAConfigName:             StringPtr("cira1"),
				MEBXPassword:               "mebx1",
				GenerateRandomMEBxPassword: true,
				Tags:                       "tags1",
				DHCPEnabled:                true,
				TLSMode:                    1,
				UserConsent:                "consent1",
				IDEREnabled:                true,
				KVMEnabled:                 true,
				SOLEnabled:                 true,
				TLSSigningAuthority:        "authority1",
				IEEE8021xProfileName:       StringPtr("ieee1"),
				IPSyncEnabled:              true,
				LocalWiFiSyncEnabled:       true,
				TenantID:                   "tenant1",
			},
			expectedErr: false,
		},
		{
			name: "Insertion with non-unique profile name",
			setup: func(dbConn *sql.DB) {
				_, err := dbConn.Exec(`
				CREATE TABLE profiles (
					profile_name TEXT NOT NULL PRIMARY KEY,
					activation TEXT,
					amt_password TEXT,
					generate_random_password BOOLEAN,
					cira_config_name TEXT,
					mebx_password TEXT,
					generate_random_mebx_password BOOLEAN,
					tags TEXT,
					dhcp_enabled BOOLEAN,
					tls_mode INTEGER,
					user_consent TEXT,
					ider_enabled BOOLEAN,
					kvm_enabled BOOLEAN,
					sol_enabled BOOLEAN,
					tls_signing_authority TEXT,
					ieee8021x_profile_name TEXT,
					ip_sync_enabled BOOLEAN,
					local_wifi_sync_enabled BOOLEAN,
					tenant_id TEXT NOT NULL
				);
				`)
				require.NoError(t, err)

				_, err = dbConn.Exec(`
				INSERT INTO profiles (profile_name, activation, amt_password, generate_random_password, cira_config_name, mebx_password, generate_random_mebx_password, tags, dhcp_enabled, tls_mode, user_consent, ider_enabled, kvm_enabled, sol_enabled, tls_signing_authority, ieee8021x_profile_name, ip_sync_enabled, local_wifi_sync_enabled, tenant_id)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
					"profile1", "activation1", "password1", true, "cira1", "mebx1", true, "tags1", true, 1, "consent1", true, true, true, "authority1", "ieee1", true, true, "tenant1")
				require.NoError(t, err)
			},
			profile: &entity.Profile{
				ProfileName:                "profile1",
				Activation:                 "activation2",
				AMTPassword:                "password2",
				GenerateRandomPassword:     true,
				CIRAConfigName:             StringPtr("cira2"),
				MEBXPassword:               "mebx2",
				GenerateRandomMEBxPassword: true,
				Tags:                       "tags2",
				DHCPEnabled:                true,
				TLSMode:                    1,
				UserConsent:                "consent2",
				IDEREnabled:                true,
				KVMEnabled:                 true,
				SOLEnabled:                 true,
				TLSSigningAuthority:        "authority2",
				IEEE8021xProfileName:       StringPtr("ieee2"),
				IPSyncEnabled:              true,
				LocalWiFiSyncEnabled:       true,
				TenantID:                   "tenant1",
			},
			expectedErr: true,
		},
		{
			name:  "Query execution error",
			setup: func(_ *sql.DB) {},
			profile: &entity.Profile{
				ProfileName: "profile2",
			},
			expectedErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dbConn, err := sql.Open("sqlite", ":memory:")
			require.NoError(t, err)
			defer dbConn.Close()

			tc.setup(dbConn)

			sqlConfig := &db.SQL{
				Builder:    squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
				Pool:       dbConn,
				IsEmbedded: true,
			}

			repo := sqldb.NewProfileRepo(sqlConfig, new(MockLogger))

			_, err = repo.Insert(context.Background(), tc.profile)

			if (err != nil) != tc.expectedErr {
				t.Errorf("Expected error status %v, got %v", tc.expectedErr, err != nil)
			}

			if !tc.expectedErr {
				var count int
				err := dbConn.QueryRow(`SELECT COUNT(*) FROM profiles WHERE profile_name = ?`, tc.profile.ProfileName).Scan(&count)
				require.NoError(t, err)

				if count == 0 {
					t.Errorf("Expected profile to be inserted, but it was not found in the database")
				}
			}
		})
	}
}
