package devicemanagement

import (
	"strings"
	"time"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman"
	amtAlarmClock "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/alarmclock"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/authorization"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/boot"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/ethernetport"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/messagelog"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/publickey"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/publicprivate"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/redirection"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/setupandconfiguration"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/timesynchronization"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/tls"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/wifiportconfiguration"
	cimBoot "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/boot"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/concrete"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/credential"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/kvm"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/models"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/power"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/software"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/wifi"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/client"
	ipsAlarmClock "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/alarmclock"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/optin"
)

type GoWSMANMessages struct {
	wsmanMessages wsman.Messages
}

func New() *GoWSMANMessages {
	return &GoWSMANMessages{}
}

func (g *GoWSMANMessages) SetupWsmanClient(device dto.Device, logAMTMessages bool) {
	clientParams := client.Parameters{
		Target:            device.Hostname,
		Username:          device.Username,
		Password:          device.Password,
		UseDigest:         true,
		UseTLS:            device.UseTLS,
		SelfSignedAllowed: device.AllowSelfSigned,
		LogAMTMessages:    logAMTMessages,
	}

	g.wsmanMessages = wsman.NewMessages(clientParams)
}

func (g *GoWSMANMessages) GetAMTVersion() ([]software.SoftwareIdentity, error) {
	response, err := g.wsmanMessages.CIM.SoftwareIdentity.Enumerate()
	if err != nil {
		return []software.SoftwareIdentity{}, err
	}

	response, err = g.wsmanMessages.CIM.SoftwareIdentity.Pull(response.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return []software.SoftwareIdentity{}, err
	}

	return response.Body.PullResponse.SoftwareIdentityItems, nil
}

func (g *GoWSMANMessages) GetSetupAndConfiguration() ([]setupandconfiguration.SetupAndConfigurationServiceResponse, error) {
	response, err := g.wsmanMessages.AMT.SetupAndConfigurationService.Enumerate()
	if err != nil {
		return []setupandconfiguration.SetupAndConfigurationServiceResponse{}, err
	}

	response, err = g.wsmanMessages.AMT.SetupAndConfigurationService.Pull(response.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return []setupandconfiguration.SetupAndConfigurationServiceResponse{}, err
	}

	return response.Body.PullResponse.SetupAndConfigurationServiceItems, nil
}

var UserConsentOptions = map[int]string{
	0:          "none",
	1:          "kvm",
	4294967295: "all",
}

func (g *GoWSMANMessages) GetFeatures() (interface{}, error) {
	redirectionResult, err := g.wsmanMessages.AMT.RedirectionService.Get()
	if err != nil {
		return nil, err
	}

	optServiceResult, err := g.wsmanMessages.IPS.OptInService.Get()
	if err != nil {
		return nil, err
	}

	kvmResult, err := g.wsmanMessages.CIM.KVMRedirectionSAP.Get()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"redirection": redirectionResult.Body.GetAndPutResponse.ListenerEnabled,
		"KVM":         kvmResult.Body.GetResponse.EnabledState == kvm.EnabledState(redirection.Enabled) || kvmResult.Body.GetResponse.EnabledState == kvm.EnabledState(redirection.EnabledButOffline),
		"SOL":         (redirectionResult.Body.GetAndPutResponse.EnabledState & redirection.Enabled) != 0,
		"IDER":        (redirectionResult.Body.GetAndPutResponse.EnabledState & redirection.Other) != 0,
		"optInState":  optServiceResult.Body.GetAndPutResponse.OptInState,
		"userConsent": UserConsentOptions[int(optServiceResult.Body.GetAndPutResponse.OptInRequired)],
	}, nil
}

func (g *GoWSMANMessages) SetFeatures(features dto.Features) (dto.Features, error) {
	// redirection
	requestedState := redirection.DisableIDERAndSOL // 32768
	listenerEnabled := 0

	if features.EnableIDER {
		requestedState++
		listenerEnabled = 1
	}

	if features.EnableSOL {
		requestedState += 2
		listenerEnabled = 1
	}

	_, err := g.wsmanMessages.AMT.RedirectionService.RequestStateChange(requestedState)
	if err != nil {
		return features, err
	}

	// kvm
	kvmRequestedState := kvm.RedirectionSAP_Disable // disabled
	if features.EnableKVM {
		kvmRequestedState = kvm.RedirectionSAP_Enable // enabled
		listenerEnabled = 1
	}

	_, err = g.wsmanMessages.CIM.KVMRedirectionSAP.RequestStateChange(kvmRequestedState)
	if err != nil {
		return features, err
	}

	// get and put redirection
	currentRedirection, err := g.wsmanMessages.AMT.RedirectionService.Get()
	if err != nil {
		return features, err
	}

	request := redirection.RedirectionRequest{
		CreationClassName:       currentRedirection.Body.GetAndPutResponse.CreationClassName,
		ElementName:             currentRedirection.Body.GetAndPutResponse.ElementName,
		EnabledState:            redirection.EnabledState(requestedState),
		ListenerEnabled:         listenerEnabled == 1,
		Name:                    currentRedirection.Body.GetAndPutResponse.Name,
		SystemCreationClassName: currentRedirection.Body.GetAndPutResponse.SystemCreationClassName,
		SystemName:              currentRedirection.Body.GetAndPutResponse.SystemName,
	}

	_, err = g.wsmanMessages.AMT.RedirectionService.Put(request)
	if err != nil {
		return features, err
	}

	// user consent
	optInResponse, err := g.wsmanMessages.IPS.OptInService.Get()
	if err != nil {
		return features, err
	}

	consentCode := optin.OptInRequiredAll // default to all if not valid user consent

	consent := strings.ToLower(features.UserConsent)

	switch consent {
	case "kvm":
		consentCode = optin.OptInRequiredKVM
	case "all":
		consentCode = optin.OptInRequiredAll
	case "none":
		consentCode = optin.OptInRequiredNone
	}

	optinRequest := optin.OptInServiceRequest{
		CreationClassName:       optInResponse.Body.GetAndPutResponse.CreationClassName,
		ElementName:             optInResponse.Body.GetAndPutResponse.ElementName,
		Name:                    optInResponse.Body.GetAndPutResponse.Name,
		OptInCodeTimeout:        optInResponse.Body.GetAndPutResponse.OptInCodeTimeout,
		OptInDisplayTimeout:     optInResponse.Body.GetAndPutResponse.OptInDisplayTimeout,
		OptInRequired:           int(consentCode),
		SystemName:              optInResponse.Body.GetAndPutResponse.SystemName,
		SystemCreationClassName: optInResponse.Body.GetAndPutResponse.SystemCreationClassName,
	}

	_, err = g.wsmanMessages.IPS.OptInService.Put(optinRequest)
	if err != nil {
		return features, err
	}

	return features, nil
}

func (g *GoWSMANMessages) GetAlarmOccurrences() ([]ipsAlarmClock.AlarmClockOccurrence, error) {
	response, err := g.wsmanMessages.IPS.AlarmClockOccurrence.Enumerate()
	if err != nil {
		return []ipsAlarmClock.AlarmClockOccurrence{}, err
	}

	response, err = g.wsmanMessages.IPS.AlarmClockOccurrence.Pull(response.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return []ipsAlarmClock.AlarmClockOccurrence{}, err
	}

	return response.Body.PullResponse.Items, nil
}

func (g *GoWSMANMessages) CreateAlarmOccurrences(name string, startTime time.Time, interval int, deleteOnCompletion bool) (amtAlarmClock.AddAlarmOutput, error) {
	alarmOccurrence := amtAlarmClock.AlarmClockOccurrence{
		InstanceID:         name,
		StartTime:          startTime,
		Interval:           interval,
		DeleteOnCompletion: deleteOnCompletion,
	}

	response, err := g.wsmanMessages.AMT.AlarmClockService.AddAlarm(alarmOccurrence)
	if err != nil {
		return amtAlarmClock.AddAlarmOutput{}, err
	}

	return response.Body.AddAlarmOutput, nil
}

func (g *GoWSMANMessages) DeleteAlarmOccurrences(instanceID string) error {
	_, err := g.wsmanMessages.IPS.AlarmClockOccurrence.Delete(instanceID)
	if err != nil {
		return err
	}

	return nil
}

func (g *GoWSMANMessages) GetHardwareInfo() (interface{}, error) {
	cspResult, err := g.wsmanMessages.CIM.ComputerSystemPackage.Get()
	if err != nil {
		return nil, err
	}

	chassisResult, err := g.wsmanMessages.CIM.Chassis.Get()
	if err != nil {
		return nil, err
	}

	cardResult, err := g.wsmanMessages.CIM.Card.Get()
	if err != nil {
		return nil, err
	}

	chipResult, err := g.wsmanMessages.CIM.Chip.Get()
	if err != nil {
		return nil, err
	}

	spEnumerateResult, err := g.wsmanMessages.CIM.SystemPackaging.Enumerate()
	if err != nil {
		return nil, err
	}

	spPullResult, err := g.wsmanMessages.CIM.SystemPackaging.Pull(spEnumerateResult.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return nil, err
	}

	ppEnumerateResult, err := g.wsmanMessages.CIM.PhysicalPackage.Enumerate()
	if err != nil {
		return nil, err
	}

	ppPullResult, err := g.wsmanMessages.CIM.PhysicalPackage.Pull(ppEnumerateResult.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return nil, err
	}

	mediaAccessEnumerateResult, err := g.wsmanMessages.CIM.MediaAccessDevice.Enumerate()
	if err != nil {
		return nil, err
	}

	mediaAccessPullResult, err := g.wsmanMessages.CIM.MediaAccessDevice.Pull(mediaAccessEnumerateResult.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return nil, err
	}

	pmEnumerateResult, err := g.wsmanMessages.CIM.PhysicalMemory.Enumerate()
	if err != nil {
		return nil, err
	}

	physicalMemoryResult, err := g.wsmanMessages.CIM.PhysicalMemory.Pull(pmEnumerateResult.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return nil, err
	}

	biosResult, err := g.wsmanMessages.CIM.BIOSElement.Get()
	if err != nil {
		return nil, err
	}

	processorResult, err := g.wsmanMessages.CIM.Processor.Get()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"CIM_ComputerSystemPackage": map[string]interface{}{
			"response":  cspResult.Body.GetResponse,
			"responses": cspResult.Body.GetResponse,
		},
		"CIM_SystemPackaging": map[string]interface{}{
			"responses": []interface{}{spPullResult.Body.PullResponse.SystemPackageItems},
		},
		"CIM_Chassis": map[string]interface{}{
			"response":  chassisResult.Body.PackageResponse,
			"responses": []interface{}{},
		}, "CIM_Chip": map[string]interface{}{
			"responses": []interface{}{chipResult.Body.PackageResponse},
		}, "CIM_Card": map[string]interface{}{
			"response":  cardResult.Body.PackageResponse,
			"responses": []interface{}{},
		}, "CIM_BIOSElement": map[string]interface{}{
			"response":  biosResult.Body.GetResponse,
			"responses": []interface{}{},
		}, "CIM_Processor": map[string]interface{}{
			"responses": []interface{}{processorResult.Body.PackageResponse},
		}, "CIM_PhysicalMemory": map[string]interface{}{
			"responses": physicalMemoryResult.Body.PullResponse.MemoryItems,
		}, "CIM_MediaAccessDevice": map[string]interface{}{
			"responses": []interface{}{mediaAccessPullResult.Body.PullResponse.MediaAccessDevices},
		}, "CIM_PhysicalPackage": map[string]interface{}{
			"responses": []interface{}{ppPullResult.Body.PullResponse.PhysicalPackage},
		},
	}, nil
}

func (g *GoWSMANMessages) GetPowerState() (interface{}, error) {
	response, err := g.wsmanMessages.CIM.ServiceAvailableToElement.Enumerate()
	if err != nil {
		return []amtAlarmClock.AlarmClockOccurrence{}, err
	}

	response, err = g.wsmanMessages.CIM.ServiceAvailableToElement.Pull(response.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return []amtAlarmClock.AlarmClockOccurrence{}, err
	}

	return map[string]interface{}{
		"powerstate": response.Body.PullResponse.AssociatedPowerManagementService[0].PowerState,
	}, nil
}

func (g *GoWSMANMessages) GetPowerCapabilities() (boot.BootCapabilitiesResponse, error) {
	response, err := g.wsmanMessages.AMT.BootCapabilities.Get()
	if err != nil {
		return boot.BootCapabilitiesResponse{}, err
	}

	return response.Body.BootCapabilitiesGetResponse, nil
}

func (g *GoWSMANMessages) GetGeneralSettings() (interface{}, error) {
	response, err := g.wsmanMessages.AMT.GeneralSettings.Get()
	if err != nil {
		return nil, err
	}

	return response.Body.GetResponse, nil
}

func (g *GoWSMANMessages) CancelUserConsent() (interface{}, error) {
	response, err := g.wsmanMessages.IPS.OptInService.CancelOptIn()
	if err != nil {
		return nil, err
	}

	return response.Body.CancelOptInResponse, nil
}

func (g *GoWSMANMessages) GetUserConsentCode() (optin.StartOptIn_OUTPUT, error) {
	response, err := g.wsmanMessages.IPS.OptInService.StartOptIn()
	if err != nil {
		return optin.StartOptIn_OUTPUT{}, err
	}

	return response.Body.StartOptInResponse, nil
}

func (g *GoWSMANMessages) SendConsentCode(code int) (interface{}, error) {
	response, err := g.wsmanMessages.IPS.OptInService.SendOptInCode(code)
	if err != nil {
		return nil, err
	}

	return response.Body.SendOptInCodeResponse, nil
}

func (g *GoWSMANMessages) GetBootData() (boot.BootCapabilitiesResponse, error) {
	bootSettingData, err := g.wsmanMessages.AMT.BootSettingData.Get()
	if err != nil {
		return boot.BootCapabilitiesResponse{}, err
	}

	return bootSettingData.Body.BootCapabilitiesGetResponse, nil
}

func (g *GoWSMANMessages) SetBootData(data boot.BootSettingDataRequest) (interface{}, error) {
	bootSettingData, err := g.wsmanMessages.AMT.BootSettingData.Put(data)
	if err != nil {
		return nil, err
	}

	return bootSettingData.Body, nil
}

func (g *GoWSMANMessages) SetBootConfigRole(_ int) (interface{}, error) {
	return nil, nil
}

func (g *GoWSMANMessages) ChangeBootOrder(bootSource string) (cimBoot.ChangeBootOrder_OUTPUT, error) {
	response, err := g.wsmanMessages.CIM.BootConfigSetting.ChangeBootOrder(cimBoot.Source(bootSource))
	if err != nil {
		return cimBoot.ChangeBootOrder_OUTPUT{}, err
	}

	return response.Body.ChangeBootOrder_OUTPUT, nil
}

func (g *GoWSMANMessages) GetAuditLog(startIndex int) (dto.AuditLog, error) {
	response, err := g.wsmanMessages.AMT.AuditLog.ReadRecords(startIndex)
	if err != nil {
		return dto.AuditLog{}, err
	}

	auditLogResponse := dto.AuditLog{}

	auditLogResponse.TotalCount = response.Body.ReadRecordsResponse.TotalRecordCount

	auditLogResponse.Records = response.Body.DecodedRecordsResponse

	return auditLogResponse, nil
}

func (g *GoWSMANMessages) GetEventLog() (messagelog.GetRecordsResponse, error) {
	response, err := g.wsmanMessages.AMT.MessageLog.GetRecords(1)
	if err != nil {
		return messagelog.GetRecordsResponse{}, err
	}

	return response.Body.GetRecordsResponse, nil
}

func (g *GoWSMANMessages) SendPowerAction(action int) (power.PowerActionResponse, error) {
	response, err := g.wsmanMessages.CIM.PowerManagementService.RequestPowerStateChange(power.PowerState(action))
	if err != nil {
		return power.PowerActionResponse{}, err
	}

	return response.Body.RequestPowerStateChangeResponse, nil
}

func (g *GoWSMANMessages) GetPublicKeyCerts() ([]publickey.PublicKeyCertificateResponse, error) {
	response, err := g.wsmanMessages.AMT.PublicKeyCertificate.Enumerate()
	if err != nil {
		return nil, err
	}

	response, err = g.wsmanMessages.AMT.PublicKeyCertificate.Pull(response.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return nil, err
	}

	return response.Body.PullResponse.PublicKeyCertificateItems, nil
}

func (g *GoWSMANMessages) GenerateKeyPair(keyAlgorithm publickey.KeyAlgorithm, keyLength publickey.KeyLength) (response publickey.Response, err error) {
	return g.wsmanMessages.AMT.PublicKeyManagementService.GenerateKeyPair(keyAlgorithm, keyLength)
}

func (g *GoWSMANMessages) UpdateAMTPassword(digestPassword string) (authorization.Response, error) {
	return g.wsmanMessages.AMT.AuthorizationService.SetAdminAclEntryEx("admin", digestPassword)
}

func (g *GoWSMANMessages) CreateTLSCredentialContext(certHandle string) (response tls.Response, err error) {
	return g.wsmanMessages.AMT.TLSCredentialContext.Create(certHandle)
}

// GetPublicPrivateKeyPairs

// NOTE: RSA Key encoded as DES PKCS#1. The Exponent (E) is 65537 (0x010001).

// When this structure is used as an output parameter (GET or PULL method),

// only the public section of the key is exported.

func (g *GoWSMANMessages) GetPublicPrivateKeyPairs() ([]publicprivate.PublicPrivateKeyPair, error) {
	response, err := g.wsmanMessages.AMT.PublicPrivateKeyPair.Enumerate()
	if err != nil {
		return nil, err
	}

	response, err = g.wsmanMessages.AMT.PublicPrivateKeyPair.Pull(response.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return nil, err
	}

	return response.Body.PullResponse.PublicPrivateKeyPairItems, nil
}

func (g *GoWSMANMessages) GetWiFiSettings() ([]wifi.WiFiEndpointSettingsResponse, error) {
	response, err := g.wsmanMessages.CIM.WiFiEndpointSettings.Enumerate()
	if err != nil {
		return nil, err
	}

	response, err = g.wsmanMessages.CIM.WiFiEndpointSettings.Pull(response.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return nil, err
	}

	return response.Body.PullResponse.EndpointSettingsItems, nil
}

func (g *GoWSMANMessages) GetEthernetSettings() ([]ethernetport.SettingsResponse, error) {
	response, err := g.wsmanMessages.AMT.EthernetPortSettings.Enumerate()
	if err != nil {
		return nil, err
	}

	response, err = g.wsmanMessages.AMT.EthernetPortSettings.Pull(response.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return nil, err
	}

	return response.Body.PullResponse.EthernetPortItems, nil
}

func (g *GoWSMANMessages) PutEthernetSettings(ethernetPortSettings ethernetport.SettingsRequest, instanceID string) (ethernetport.Response, error) {
	return g.wsmanMessages.AMT.EthernetPortSettings.Put(instanceID, ethernetPortSettings)
}

func (g *GoWSMANMessages) DeletePublicPrivateKeyPair(instanceID string) error {
	_, err := g.wsmanMessages.AMT.PublicPrivateKeyPair.Delete(instanceID)

	return err
}

func (g *GoWSMANMessages) DeletePublicCert(instanceID string) error {
	_, err := g.wsmanMessages.AMT.PublicKeyCertificate.Delete(instanceID)

	return err
}

func (g *GoWSMANMessages) GetCredentialRelationships() ([]credential.CredentialContext, error) {
	response, err := g.wsmanMessages.CIM.CredentialContext.Enumerate()
	if err != nil {
		return nil, err
	}

	response, err = g.wsmanMessages.CIM.CredentialContext.Pull(response.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return nil, err
	}

	return response.Body.PullResponse.Items, nil
}

func (g *GoWSMANMessages) GetConcreteDependencies() ([]concrete.ConcreteDependency, error) {
	response, err := g.wsmanMessages.CIM.ConcreteDependency.Enumerate()
	if err != nil {
		return nil, err
	}

	response, err = g.wsmanMessages.CIM.ConcreteDependency.Pull(response.Body.EnumerateResponse.EnumerationContext)
	if err != nil {
		return nil, err
	}

	return response.Body.PullResponse.Items, nil
}

func (g *GoWSMANMessages) DeleteWiFiSetting(instanceID string) error {
	_, err := g.wsmanMessages.CIM.WiFiEndpointSettings.Delete(instanceID)

	return err
}

func (g *GoWSMANMessages) AddTrustedRootCert(caCert string) (handle string, err error) {
	response, err := g.wsmanMessages.AMT.PublicKeyManagementService.AddTrustedRootCertificate(caCert)
	if err != nil {
		return "", err
	}

	if len(response.Body.AddTrustedRootCertificate_OUTPUT.CreatedCertificate.ReferenceParameters.SelectorSet.Selectors) > 0 {
		handle = response.Body.AddTrustedRootCertificate_OUTPUT.CreatedCertificate.ReferenceParameters.SelectorSet.Selectors[0].Text
	}

	return handle, nil
}

func (g *GoWSMANMessages) AddClientCert(clientCert string) (handle string, err error) {
	response, err := g.wsmanMessages.AMT.PublicKeyManagementService.AddCertificate(clientCert)
	if err != nil {
		return "", err
	}

	if len(response.Body.AddCertificate_OUTPUT.CreatedCertificate.ReferenceParameters.SelectorSet.Selectors) > 0 {
		handle = response.Body.AddCertificate_OUTPUT.CreatedCertificate.ReferenceParameters.SelectorSet.Selectors[0].Text
	}

	return handle, nil
}

func (g *GoWSMANMessages) AddPrivateKey(privateKey string) (handle string, err error) {
	response, err := g.wsmanMessages.AMT.PublicKeyManagementService.AddKey(privateKey)
	if err != nil {
		return "", err
	}

	if len(response.Body.AddKey_OUTPUT.CreatedKey.ReferenceParameters.SelectorSet.Selectors) > 0 {
		handle = response.Body.AddKey_OUTPUT.CreatedKey.ReferenceParameters.SelectorSet.Selectors[0].Text
	}

	return handle, nil
}

func (g *GoWSMANMessages) DeleteKeyPair(instanceID string) error {
	_, err := g.wsmanMessages.AMT.PublicKeyManagementService.Delete(instanceID)

	return err
}

func (g *GoWSMANMessages) EnableWiFi() error {
	response, err := g.wsmanMessages.AMT.WiFiPortConfigurationService.Get()
	if err != nil {
		return err
	}

	// if local sync not enable, enable it
	if response.Body.WiFiPortConfigurationService.LocalProfileSynchronizationEnabled == wifiportconfiguration.LocalSyncDisabled {
		putRequest := wifiportconfiguration.WiFiPortConfigurationServiceRequest{
			RequestedState:                     response.Body.WiFiPortConfigurationService.RequestedState,
			EnabledState:                       response.Body.WiFiPortConfigurationService.EnabledState,
			HealthState:                        response.Body.WiFiPortConfigurationService.HealthState,
			ElementName:                        response.Body.WiFiPortConfigurationService.ElementName,
			SystemCreationClassName:            response.Body.WiFiPortConfigurationService.SystemCreationClassName,
			SystemName:                         response.Body.WiFiPortConfigurationService.SystemName,
			CreationClassName:                  response.Body.WiFiPortConfigurationService.CreationClassName,
			Name:                               response.Body.WiFiPortConfigurationService.Name,
			LocalProfileSynchronizationEnabled: wifiportconfiguration.UnrestrictedSync,
			LastConnectedSsidUnderMeControl:    response.Body.WiFiPortConfigurationService.LastConnectedSsidUnderMeControl,
			NoHostCsmeSoftwarePolicy:           response.Body.WiFiPortConfigurationService.NoHostCsmeSoftwarePolicy,
			UEFIWiFiProfileShareEnabled:        response.Body.WiFiPortConfigurationService.UEFIWiFiProfileShareEnabled,
		}

		_, err = g.wsmanMessages.AMT.WiFiPortConfigurationService.Put(putRequest)
		if err != nil {
			return err
		}
	}

	// always turn wifi on via state change request
	// Enumeration 32769 - WiFi is enabled in S0 + Sx/AC
	_, err = g.wsmanMessages.CIM.WiFiPort.RequestStateChange(int(wifi.EnabledStateWifiEnabledS0SxAC))
	if err != nil {
		return err // utils.WSMANMessageError
	}

	return nil
}

func (g *GoWSMANMessages) AddWiFiSettings(wifiEndpointSettings wifi.WiFiEndpointSettingsRequest, ieee8021xSettings models.IEEE8021xSettings, wifiEndpoint, clientCredential, caCredential string) (response wifiportconfiguration.Response, err error) {
	return g.wsmanMessages.AMT.WiFiPortConfigurationService.AddWiFiSettings(wifiEndpointSettings, ieee8021xSettings, wifiEndpoint, clientCredential, caCredential)
}

func (g *GoWSMANMessages) PUTTLSSettings(instanceID string, tlsSettingData tls.SettingDataRequest) (response tls.Response, err error) {
	return g.wsmanMessages.AMT.TLSSettingData.Put(instanceID, tlsSettingData)
}

func (g *GoWSMANMessages) GetLowAccuracyTimeSynch() (response timesynchronization.Response, err error) {
	return g.wsmanMessages.AMT.TimeSynchronizationService.GetLowAccuracyTimeSynch()
}

func (g *GoWSMANMessages) SetHighAccuracyTimeSynch(ta0, tm1, tm2 int64) (response timesynchronization.Response, err error) {
	return g.wsmanMessages.AMT.TimeSynchronizationService.SetHighAccuracyTimeSynch(ta0, tm1, tm2)
}

func (g *GoWSMANMessages) EnumerateTLSSettingData() (response tls.Response, err error) {
	return g.wsmanMessages.AMT.TLSSettingData.Enumerate()
}

func (g *GoWSMANMessages) PullTLSSettingData(enumerationContext string) (response tls.Response, err error) {
	return g.wsmanMessages.AMT.TLSSettingData.Pull(enumerationContext)
}

func (g *GoWSMANMessages) CommitChanges() (response setupandconfiguration.Response, err error) {
	return g.wsmanMessages.AMT.SetupAndConfigurationService.CommitChanges()
}

func (g *GoWSMANMessages) GeneratePKCS10RequestEx(keyPair, nullSignedCertificateRequest string, signingAlgorithm publickey.SigningAlgorithm) (response publickey.Response, err error) {
	return g.wsmanMessages.AMT.PublicKeyManagementService.GeneratePKCS10RequestEx(keyPair, nullSignedCertificateRequest, signingAlgorithm)
}

func (g *GoWSMANMessages) RequestRedirectionStateChange(requestedState redirection.RequestedState) (response redirection.Response, err error) {
	return g.wsmanMessages.AMT.RedirectionService.RequestStateChange(requestedState)
}

func (g *GoWSMANMessages) RequestKVMStateChange(requestedState kvm.KVMRedirectionSAPRequestStateChangeInput) (response kvm.Response, err error) {
	return g.wsmanMessages.CIM.KVMRedirectionSAP.RequestStateChange(requestedState)
}

func (g *GoWSMANMessages) PutRedirectionState(requestedState redirection.RedirectionRequest) (response redirection.Response, err error) {
	return g.wsmanMessages.AMT.RedirectionService.Put(requestedState)
}

func (g *GoWSMANMessages) GetRedirectionService() (response redirection.Response, err error) {
	return g.wsmanMessages.AMT.RedirectionService.Get()
}

func (g *GoWSMANMessages) GetIpsOptInService() (response optin.Response, err error) {
	return g.wsmanMessages.IPS.OptInService.Get()
}

func (g *GoWSMANMessages) PutIpsOptInService(request optin.OptInServiceRequest) (response optin.Response, err error) {
	return g.wsmanMessages.IPS.OptInService.Put(request)
}
