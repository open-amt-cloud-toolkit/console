package devices

import (
	"context"
	"errors"
	"strings"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/amterror"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/redirection"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/kvm"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/ips/optin"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	dtov2 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v2"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
)

func (uc *UseCase) GetFeatures(c context.Context, guid string) (settingsResults dto.Features, settingsResultsV2 dtov2.Features, err error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return dto.Features{}, dtov2.Features{}, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	// Get redirection settings from AMT
	err = getRedirectionService(&settingsResultsV2, device)
	if err != nil {
		return settingsResults, settingsResultsV2, err
	}

	// translate v2 features to v1 - remove this once v1 is deprecated
	settingsResults.EnableSOL = settingsResultsV2.EnableSOL
	settingsResults.EnableIDER = settingsResultsV2.EnableIDER
	settingsResults.Redirection = settingsResultsV2.Redirection

	// Get optinservice settings from AMT
	err = getUserConsent(&settingsResultsV2, device)
	if err != nil {
		return settingsResults, settingsResultsV2, err
	}

	// translate v2 features to v1 - remove this once v1 is deprecated
	settingsResults.UserConsent = settingsResultsV2.UserConsent
	settingsResults.OptInState = settingsResultsV2.OptInState

	// Get KVM settings from AMT
	err = getKVM(&settingsResultsV2, device)
	if err != nil {
		return settingsResults, settingsResultsV2, err
	}

	// translate v2 features to v1 - remove this once v1 is deprecated
	settingsResults.EnableKVM = settingsResultsV2.EnableKVM
	settingsResults.KVMAvailable = settingsResultsV2.KVMAvailable

	return settingsResults, settingsResultsV2, nil
}

func (uc *UseCase) SetFeatures(c context.Context, guid string, features dto.Features) (settingsResults dto.Features, settingsResultsV2 dtov2.Features, err error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return dto.Features{}, dtov2.Features{}, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	// redirection
	state, listenerEnabled, err := redirectionRequestStateChange(features.EnableSOL, features.EnableIDER, &settingsResultsV2, device)
	if err != nil {
		return settingsResults, settingsResultsV2, err
	}

	// translate v2 features to v1 - remove this once v1 is deprecated
	settingsResults.EnableSOL = settingsResultsV2.EnableSOL
	settingsResults.EnableIDER = settingsResultsV2.EnableIDER

	// kvm
	kvmListenerEnabled, err := setKVM(features.EnableKVM, &settingsResultsV2, device)
	if err != nil {
		return settingsResults, settingsResultsV2, err
	}

	// translate v2 features to v1 - remove this once v1 is deprecated
	settingsResults.EnableKVM = settingsResultsV2.EnableKVM

	// get and put redirection
	err = setRedirectionService(state, listenerEnabled, kvmListenerEnabled, device)
	if err != nil {
		return settingsResults, settingsResultsV2, err
	}

	// Update Redirection, this is important when KVM, IDER and SOL are all false
	settingsResults.Redirection = listenerEnabled == 1 || kvmListenerEnabled == 1
	settingsResultsV2.Redirection = listenerEnabled == 1 || kvmListenerEnabled == 1

	// user consent
	err = setUserConsent(features.UserConsent, device)
	if err != nil {
		return settingsResults, settingsResultsV2, err
	}

	settingsResults.UserConsent = features.UserConsent
	settingsResultsV2.UserConsent = features.UserConsent

	return settingsResults, settingsResultsV2, err
}

func handleAMTKVMError(err error, results *dtov2.Features) bool {
	amtErr := &amterror.AMTError{}
	if errors.As(err, &amtErr) {
		if strings.Contains(amtErr.SubCode, "DestinationUnreachable") {
			results.EnableKVM = false
			results.KVMAvailable = false

			return true
		}
	}

	return false
}

func getSOLAndIDERState(enabledState redirection.EnabledState) (iderEnabled, solEnabled bool) {
	//nolint:exhaustive // we only care about IDER and SOL states. Other scenarios are handled by the default case.
	switch enabledState {
	case redirection.IDERAndSOLAreDisabled:
		return false, false
	case redirection.IDERIsEnabledAndSOLIsDisabled:
		return true, false
	case redirection.SOLIsEnabledAndIDERIsDisabled:
		return false, true
	case redirection.IDERAndSOLAreEnabled:
		return true, true
	default:
		return false, false // default case if state is invalid
	}
}

func redirectionRequestStateChange(enableSOL, enableIDER bool, results *dtov2.Features, w wsman.Management) (state redirection.EnabledState, listenerEnabled int, err error) {
	requestedState, listenerEnabled, err := w.RequestAMTRedirectionServiceStateChange(enableIDER, enableSOL)
	if err != nil {
		return 0, 0, err
	}

	state = redirection.EnabledState(requestedState)
	iderEnabled, solEnabled := getSOLAndIDERState(state)
	results.EnableSOL = solEnabled
	results.EnableIDER = iderEnabled
	results.Redirection = listenerEnabled == 1

	return state, listenerEnabled, nil
}

func getKVM(results *dtov2.Features, w wsman.Management) error {
	kvmResult, err := w.GetKVMRedirection()
	if err != nil {
		isAMTErr := handleAMTKVMError(err, results)
		if !isAMTErr {
			return err
		}
	} else {
		results.EnableKVM = kvmResult.Body.GetResponse.EnabledState == kvm.EnabledState(redirection.Enabled) || kvmResult.Body.GetResponse.EnabledState == kvm.EnabledState(redirection.EnabledButOffline)
		results.KVMAvailable = true
	}

	return nil
}

func setKVM(enableKVM bool, results *dtov2.Features, w wsman.Management) (kvmListenerEnabled int, err error) {
	kvmListenerEnabled, err = w.SetKVMRedirection(enableKVM)
	if err != nil {
		isAMTErr := handleAMTKVMError(err, results)
		if !isAMTErr {
			return 0, err
		}
	} else {
		results.EnableKVM = kvmListenerEnabled == 1
		results.KVMAvailable = true
	}

	return kvmListenerEnabled, nil
}

func getRedirectionService(results *dtov2.Features, w wsman.Management) error {
	redirectionResult, err := w.GetAMTRedirectionService()
	if err != nil {
		return err
	}

	iderEnabled, solEnabled := getSOLAndIDERState(redirectionResult.Body.GetAndPutResponse.EnabledState)
	results.EnableSOL = solEnabled
	results.EnableIDER = iderEnabled
	results.Redirection = redirectionResult.Body.GetAndPutResponse.ListenerEnabled

	return nil
}

func setRedirectionService(state redirection.EnabledState, listenerEnabled, kvmListenerEnabled int, w wsman.Management) error {
	currentRedirection, err := w.GetAMTRedirectionService()
	if err != nil {
		return err
	}

	request := redirection.RedirectionRequest{
		CreationClassName:       currentRedirection.Body.GetAndPutResponse.CreationClassName,
		ElementName:             currentRedirection.Body.GetAndPutResponse.ElementName,
		EnabledState:            state,
		ListenerEnabled:         listenerEnabled == 1 || kvmListenerEnabled == 1,
		Name:                    currentRedirection.Body.GetAndPutResponse.Name,
		SystemCreationClassName: currentRedirection.Body.GetAndPutResponse.SystemCreationClassName,
		SystemName:              currentRedirection.Body.GetAndPutResponse.SystemName,
	}

	_, err = w.SetAMTRedirectionService(request)
	if err != nil {
		return err
	}

	return nil
}

func getUserConsent(results *dtov2.Features, w wsman.Management) error {
	optServiceResult, err := w.GetIPSOptInService()
	if err != nil {
		return err
	}

	results.UserConsent = UserConsentOptions[int(optServiceResult.Body.GetAndPutResponse.OptInRequired)]
	results.OptInState = optServiceResult.Body.GetAndPutResponse.OptInState

	return nil
}

func setUserConsent(userConsent string, w wsman.Management) error {
	optInResponse, err := w.GetIPSOptInService()
	if err != nil {
		return err
	}

	optinRequest := optin.OptInServiceRequest{
		CreationClassName:       optInResponse.Body.GetAndPutResponse.CreationClassName,
		ElementName:             optInResponse.Body.GetAndPutResponse.ElementName,
		Name:                    optInResponse.Body.GetAndPutResponse.Name,
		OptInCodeTimeout:        optInResponse.Body.GetAndPutResponse.OptInCodeTimeout,
		OptInDisplayTimeout:     optInResponse.Body.GetAndPutResponse.OptInDisplayTimeout,
		OptInRequired:           determineConsentCode(userConsent),
		SystemName:              optInResponse.Body.GetAndPutResponse.SystemName,
		SystemCreationClassName: optInResponse.Body.GetAndPutResponse.SystemCreationClassName,
	}

	err = w.SetIPSOptInService(optinRequest)
	if err != nil {
		return err
	}

	return nil
}

func determineConsentCode(consent string) int {
	consentCode := optin.OptInRequiredAll // default to all if not valid user consent

	consent = strings.ToLower(consent)

	switch consent {
	case "kvm":
		consentCode = optin.OptInRequiredKVM
	case "all":
		consentCode = optin.OptInRequiredAll
	case "none":
		consentCode = optin.OptInRequiredNone
	}

	return int(consentCode)
}

var UserConsentOptions = map[int]string{
	0:          "none",
	1:          "kvm",
	4294967295: "all",
}
