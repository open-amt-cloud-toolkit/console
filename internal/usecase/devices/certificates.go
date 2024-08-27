package devices

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha1" //nolint:gosec // SHA-1 is used for thumbprint not signature
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/publickey"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/amt/publicprivate"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/concrete"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/credential"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
)

const (
	TypeWireless string = "Wireless"
	TypeTLS      string = "TLS"
	TypeWired    string = "Wired"
)

func processConcreteDependencies(certificateHandle string, profileAssociation *dto.ProfileAssociation, dependancyItems []concrete.ConcreteDependency, keyPairItems []publicprivate.RefinedPublicPrivateKeyPair) {
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

func buildCertificateAssociations(profileAssociation dto.ProfileAssociation, securitySettings *dto.SecuritySettings) {
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

func buildProfileAssociations(certificateHandle string, profileAssociation *dto.ProfileAssociation, response wsman.Certificates, securitySettings *dto.SecuritySettings) {
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

func processCertificates(contextItems []credential.CredentialContext, response wsman.Certificates, profileType string, securitySettings *dto.SecuritySettings) {
	for idx := range contextItems {
		var profileAssociation dto.ProfileAssociation

		profileAssociation.Type = profileType
		profileAssociation.ProfileID = strings.TrimPrefix(contextItems[idx].ElementProvidingContext.ReferenceParameters.SelectorSet.Selectors[0].Text, "Intel(r) AMT:IEEE 802.1x Settings ")
		certificateHandle := contextItems[idx].ElementInContext.ReferenceParameters.SelectorSet.Selectors[0].Text

		buildProfileAssociations(certificateHandle, &profileAssociation, response, securitySettings)
		buildCertificateAssociations(profileAssociation, securitySettings)
	}
}

func (uc *UseCase) GetCertificates(c context.Context, guid string) (dto.SecuritySettings, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return dto.SecuritySettings{}, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	response, err := device.GetCertificates()
	if err != nil {
		return dto.SecuritySettings{}, err
	}

	securitySettings := dto.SecuritySettings{
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

func (uc *UseCase) GetDeviceCertificate(c context.Context, guid string) (dto.Certificate, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return dto.Certificate{}, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	cert1, err := device.GetDeviceCertificate()
	if err != nil {
		return dto.Certificate{}, err
	}

	var certDTOs []dto.Certificate

	for _, certBytes := range cert1.Certificate {
		// Parse each certificate byte slice into an x509.Certificate
		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			uc.log.Warn(fmt.Sprintf("Failed to parse certificate: %v", err))

			continue
		}

		// Populate the DTO with certificate information
		certDTO := populateCertificateDTO(cert)
		certDTOs = append(certDTOs, certDTO)
	}

	return certDTOs[0], nil
}

func populateCertificateDTO(cert *x509.Certificate) dto.Certificate {
	// Compute the SHA-1 and SHA-256 fingerprints
	sha1Fingerprint := sha1.Sum(cert.Raw) //nolint:gosec // SHA-1 is used for thumbprint not signature
	sha256Fingerprint := sha256.Sum256(cert.Raw)

	// Determine the public key size
	var publicKeySize int
	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		publicKeySize = pub.N.BitLen()
	case *ecdsa.PublicKey:
		publicKeySize = pub.Curve.Params().BitSize
	default:
		publicKeySize = 0 // Unknown or unsupported key type
	}

	// Populate the dto.Certificate struct
	return dto.Certificate{
		CommonName:         cert.Subject.CommonName,
		IssuerName:         cert.Issuer.CommonName,
		SerialNumber:       cert.SerialNumber.String(),
		NotBefore:          cert.NotBefore,
		NotAfter:           cert.NotAfter,
		DNSNames:           cert.DNSNames,
		SHA1Fingerprint:    hex.EncodeToString(sha1Fingerprint[:]),
		SHA256Fingerprint:  hex.EncodeToString(sha256Fingerprint[:]),
		PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
		PublicKeySize:      publicKeySize,
	}
}
