// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/software"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/alarmclock"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Translation -.
	Domain interface {
		GetCount(context.Context, string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Domain, error)
		GetDomainByDomainSuffix(ctx context.Context, domainSuffix, tenantID string) (*entity.Domain, error)
		GetByName(ctx context.Context, name, tenantID string) (*entity.Domain, error)
		Delete(ctx context.Context, name, tenantID string) (bool, error)
		Update(ctx context.Context, d *entity.Domain) (bool, error)
		Insert(ctx context.Context, d *entity.Domain) (string, error)
	}
	Device interface {
		GetCount(context.Context, string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Device, error)
		GetByID(ctx context.Context, guid, tenantID string) (entity.Device, error)
		// GetByName(ctx context.Context, name, tenantID string) (*entity.Device, error)
		Delete(ctx context.Context, name, tenantID string) (bool, error)
		Update(ctx context.Context, d *entity.Device) (bool, error)
		Insert(ctx context.Context, d *entity.Device) (string, error)
	}
	Profile interface {
		GetCount(ctx context.Context, tenantID string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Profile, error)
		GetByName(ctx context.Context, profileName, tenantID string) (entity.Profile, error)
		Delete(ctx context.Context, profileName, tenantID string) (bool, error)
		Update(ctx context.Context, p *entity.Profile) (bool, error)
		Insert(ctx context.Context, p *entity.Profile) (string, error)
	}
	DeviceManagement interface {
		SetupWsmanClient(device entity.Device, logAMTMessages bool)
		GetAMTVersion() ([]software.SoftwareIdentity, error)
		GetFeatures() (interface{}, error)
		GetAlarmOccurrences() ([]alarmclock.AlarmClockOccurrence, error)
		GetHardwareInfo() (interface{}, error)
		GetPowerState() (interface{}, error)
		GetPowerCapabilities() (interface{}, error)
		GetGeneralSettings() (interface{}, error)
		CancelUserConsent() (interface{}, error)
		GetUserConsentCode() (interface{}, error)
		SendConsentCode(code int) (interface{}, error)
		// Unprovision(int) (setupandconfiguration.Response, error)
		// GetGeneralSettings() (general.Response, error)
		// HostBasedSetupService(digestRealm string, password string) (hostbasedsetup.Response, error)
		// GetHostBasedSetupService() (hostbasedsetup.Response, error)
		// AddNextCertInChain(cert string, isLeaf bool, isRoot bool) (hostbasedsetup.Response, error)
		// HostBasedSetupServiceAdmin(password string, digestRealm string, nonce []byte, signature string) (hostbasedsetup.Response, error)
		// SetupMEBX(string) (response setupandconfiguration.Response, err error)
		// GetPublicKeyCerts() ([]publickey.PublicKeyCertificateResponse, error)
		// GetPublicPrivateKeyPairs() ([]publicprivate.PublicPrivateKeyPair, error)
		// DeletePublicPrivateKeyPair(instanceId string) error
		// DeletePublicCert(instanceId string) error
		// GetCredentialRelationships() ([]credential.CredentialContext, error)
		// GetConcreteDependencies() ([]concrete.ConcreteDependency, error)
		// AddTrustedRootCert(caCert string) (string, error)
		// AddClientCert(clientCert string) (string, error)
		// AddPrivateKey(privateKey string) (string, error)
		// DeleteKeyPair(instanceID string) error
		// GetLowAccuracyTimeSynch() (response timesynchronization.Response, err error)
		// SetHighAccuracyTimeSynch(ta0 int64, tm1 int64, tm2 int64) (response timesynchronization.Response, err error)
		// GenerateKeyPair(keyAlgorithm publickey.KeyAlgorithm, keyLength publickey.KeyLength) (response publickey.Response, err error)
		// UpdateAMTPassword(passwordBase64 string) (authorization.Response, error)
		// // WiFi
		// GetWiFiSettings() ([]wifi.WiFiEndpointSettingsResponse, error)
		// DeleteWiFiSetting(instanceId string) error
		// EnableWiFi() error
		// AddWiFiSettings(wifiEndpointSettings wifi.WiFiEndpointSettingsRequest, ieee8021xSettings models.IEEE8021xSettings, wifiEndpoint, clientCredential, caCredential string) (wifiportconfiguration.Response, error)
		// // Wired
		// GetEthernetSettings() ([]ethernetport.SettingsResponse, error)
		// PutEthernetSettings(ethernetPortSettings ethernetport.SettingsRequest, instanceId string) (ethernetport.Response, error)
		// // TLS
		// CreateTLSCredentialContext(certHandle string) (response tls.Response, err error)
		// EnumerateTLSSettingData() (response tls.Response, err error)
		// PullTLSSettingData(enumerationContext string) (response tls.Response, err error)
		// PUTTLSSettings(instanceID string, tlsSettingData tls.SettingDataRequest) (response tls.Response, err error)

		// CommitChanges() (response setupandconfiguration.Response, err error)
		// GeneratePKCS10RequestEx(keyPair, nullSignedCertificateRequest string, signingAlgorithm publickey.SigningAlgorithm) (response publickey.Response, err error)

		// RequestRedirectionStateChange(requestedState redirection.RequestedState) (response redirection.Response, err error)
		// RequestKVMStateChange(requestedState kvm.KVMRedirectionSAPRequestedStateInputs) (response kvm.Response, err error)
		// PutRedirectionState(requestedState redirection.RedirectionRequest) (response redirection.Response, err error)
		// GetRedirectionService() (response redirection.Response, err error)
		// GetIpsOptInService() (response optin.Response, err error)
		// PutIpsOptInService(request optin.OptInServiceRequest) (response optin.Response, err error)
	}
	IEEE8021xProfile interface {
		CheckProfileExists(ctx context.Context, profileName string, tenantID string) (bool, error)
		GetCount(ctx context.Context, tenantID string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.IEEE8021xConfig, error)
		GetByName(ctx context.Context, profileName, tenantID string) (entity.IEEE8021xConfig, error)
		Delete(ctx context.Context, profileName, tenantID string) (bool, error)
		Update(ctx context.Context, p *entity.IEEE8021xConfig) (bool, error)
		Insert(ctx context.Context, p *entity.IEEE8021xConfig) (string, error)
	}
	CIRAConfig interface {
		GetCount(ctx context.Context, tenantID string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.CIRAConfig, error)
		GetByName(ctx context.Context, configName, tenantID string) (entity.CIRAConfig, error)
		Delete(ctx context.Context, profileName, tenantID string) (bool, error)
		Update(ctx context.Context, p *entity.CIRAConfig) (bool, error)
		Insert(ctx context.Context, p *entity.CIRAConfig) (string, error)
	}
	WirelessProfile interface {
		CheckProfileExists(ctx context.Context, profileName string, tenantID string) (bool, error)
		GetCount(ctx context.Context, tenantID string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.WirelessConfig, error)
		GetByName(ctx context.Context, guid, tenantID string) (entity.WirelessConfig, error)
		Delete(ctx context.Context, profileName, tenantID string) (bool, error)
		Update(ctx context.Context, p *entity.WirelessConfig) (bool, error)
		Insert(ctx context.Context, p *entity.WirelessConfig) (string, error)
	}
)
