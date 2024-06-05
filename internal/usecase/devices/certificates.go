package devices

import (
	"context"
	"reflect"
	"strings"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/publickey"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/publicprivate"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/concrete"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/credential"

	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
)

const (
	TypeWireless string = "Wireless"
	TypeTLS      string = "TLS"
	TypeWired    string = "Wired"
)

type SecuritySettings struct {
	ProfileAssociation []ProfileAssociation `json:"ProfileAssociation"`
	Certificates       interface{}          `json:"Certificates"`
	Keys               interface{}          `json:"PublicKeys"`
}

type ProfileAssociation struct {
	Type              string      `json:"Type"`
	ProfileID         string      `json:"ProfileID"`
	RootCertificate   interface{} `json:"RootCertificate,omitempty"`
	ClientCertificate interface{} `json:"ClientCertificate,omitempty"`
	Key               interface{} `json:"PublicKey,omitempty"`
}

func processConcreteDependencies(certificateHandle string, profileAssociation *ProfileAssociation, dependancyItems []concrete.ConcreteDependency, keyPairItems []publicprivate.RefinedPublicPrivateKeyPair) {
	for x := range dependancyItems {
		if dependancyItems[x].Antecedent.ReferenceParameters.SelectorSet.Selectors[0].Text != certificateHandle {
			continue
		}

		keyHandle := dependancyItems[x].Dependent.ReferenceParameters.SelectorSet.Selectors[0].Text

		for i := range keyPairItems {
			if keyPairItems[i].InstanceID == keyHandle {
				profileAssociation.Key = keyPairItems[i]

				break
			}
		}
	}
}

func buildCertificateAssociations(profileAssociation ProfileAssociation, securitySettings *SecuritySettings) {
	var publicKeyHandle string

	// If a client cert, update the associated public key w/ the cert's handle
	if profileAssociation.ClientCertificate != nil {
		// Loop thru public keys looking for the one that matches the current profileAssociation's key
		for i, existingKeyPair := range securitySettings.Keys.(publicprivate.RefinedPullResponse).PublicPrivateKeyPairItems {
			// If found update that key with the profileAssociation's certificate handle
			if existingKeyPair.InstanceID == profileAssociation.Key.(publicprivate.RefinedPublicPrivateKeyPair).InstanceID {
				securitySettings.Keys.(publicprivate.RefinedPullResponse).PublicPrivateKeyPairItems[i].CertificateHandle = profileAssociation.ClientCertificate.(publickey.RefinedPublicKeyCertificateResponse).InstanceID
				// save this public key handle since we know it pairs with the profileAssociation's certificate
				publicKeyHandle = securitySettings.Keys.(publicprivate.RefinedPullResponse).PublicPrivateKeyPairItems[i].InstanceID

				break
			}
		}
	}

	// Loop thru certificates looking for the one that matches the current profileAssociation's certificate and append profile name
	for i := range securitySettings.Certificates.(publickey.RefinedPullResponse).PublicKeyCertificateItems {
		if (profileAssociation.ClientCertificate != nil && securitySettings.Certificates.(publickey.RefinedPullResponse).PublicKeyCertificateItems[i].InstanceID == profileAssociation.ClientCertificate.(publickey.RefinedPublicKeyCertificateResponse).InstanceID) ||
			(profileAssociation.RootCertificate != nil && securitySettings.Certificates.(publickey.RefinedPullResponse).PublicKeyCertificateItems[i].InstanceID == profileAssociation.RootCertificate.(publickey.RefinedPublicKeyCertificateResponse).InstanceID) {
			// if client cert found, associate the previously found key handle with it
			if !securitySettings.Certificates.(publickey.RefinedPullResponse).PublicKeyCertificateItems[i].TrustedRootCertificate {
				securitySettings.Certificates.(publickey.RefinedPullResponse).PublicKeyCertificateItems[i].PublicKeyHandle = publicKeyHandle
			}

			securitySettings.Certificates.(publickey.RefinedPullResponse).PublicKeyCertificateItems[i].AssociatedProfiles = append(securitySettings.Certificates.(publickey.RefinedPullResponse).PublicKeyCertificateItems[i].AssociatedProfiles, profileAssociation.ProfileID)

			break
		}
	}
}

func buildProfileAssociations(certificateHandle string, profileAssociation *ProfileAssociation, response wsman.Certificates, securitySettings *SecuritySettings) {
	isNewProfileAssociation := true

	for idx := range response.PublicKeyCertificateResponse.PublicKeyCertificateItems {
		if response.PublicKeyCertificateResponse.PublicKeyCertificateItems[idx].InstanceID != certificateHandle {
			continue
		}

		if response.PublicKeyCertificateResponse.PublicKeyCertificateItems[idx].TrustedRootCertificate {
			profileAssociation.RootCertificate = response.PublicKeyCertificateResponse.PublicKeyCertificateItems[idx]

			continue
		}

		profileAssociation.ClientCertificate = response.PublicKeyCertificateResponse.PublicKeyCertificateItems[idx]

		processConcreteDependencies(certificateHandle, profileAssociation, response.ConcreteDependencyResponse.Items, response.PublicPrivateKeyPairResponse.PublicPrivateKeyPairItems)
	}

	// Check if the certificate is already in the list
	for idx := range securitySettings.ProfileAssociation {
		if !(securitySettings.ProfileAssociation[idx].ProfileID == profileAssociation.ProfileID) {
			continue
		}

		if profileAssociation.RootCertificate != nil {
			securitySettings.ProfileAssociation[idx].RootCertificate = profileAssociation.RootCertificate
		}

		if profileAssociation.ClientCertificate != nil {
			securitySettings.ProfileAssociation[idx].ClientCertificate = profileAssociation.ClientCertificate
		}

		if profileAssociation.Key != nil {
			securitySettings.ProfileAssociation[idx].Key = profileAssociation.Key
		}

		isNewProfileAssociation = false

		break
	}

	// If the profile is not in the list, add it
	if isNewProfileAssociation {
		securitySettings.ProfileAssociation = append(securitySettings.ProfileAssociation, *profileAssociation)
	}
}

func processCertificates(contextItems []credential.CredentialContext, response wsman.Certificates, profileType string, securitySettings *SecuritySettings) {
	for idx := range contextItems {
		var profileAssociation ProfileAssociation

		profileAssociation.Type = profileType
		profileAssociation.ProfileID = strings.TrimPrefix(contextItems[idx].ElementProvidingContext.ReferenceParameters.SelectorSet.Selectors[0].Text, "Intel(r) AMT:IEEE 802.1x Settings ")
		certificateHandle := contextItems[idx].ElementInContext.ReferenceParameters.SelectorSet.Selectors[0].Text

		buildProfileAssociations(certificateHandle, &profileAssociation, response, securitySettings)
		buildCertificateAssociations(profileAssociation, securitySettings)
	}
}

func (uc *UseCase) GetCertificates(c context.Context, guid string) (interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil || item.GUID == "" {
		return nil, err
	}

	uc.device.SetupWsmanClient(*item, false, true)

	response, err := uc.device.GetCertificates()
	if err != nil {
		return nil, err
	}

	securitySettings := SecuritySettings{
		Certificates: response.PublicKeyCertificateResponse,
		Keys:         response.PublicPrivateKeyPairResponse,
	}

	if !reflect.DeepEqual(response.CIMCredentialContextResponse, credential.PullResponse{}) {
		processCertificates(response.CIMCredentialContextResponse.Items.CredentialContextTLS, response, TypeTLS, &securitySettings)
		processCertificates(response.CIMCredentialContextResponse.Items.CredentialContext, response, TypeWireless, &securitySettings)
		processCertificates(response.CIMCredentialContextResponse.Items.CredentialContext8021x, response, TypeWired, &securitySettings)
	}

	return securitySettings, nil
}
