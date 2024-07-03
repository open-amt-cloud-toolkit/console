package profiles

import (
	"context"
	"errors"
	"strings"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/ieee8021xconfigs"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/profilewificonfigs"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/sqldb"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/wificonfigs"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

// UseCase -.
type UseCase struct {
	repo              Repository
	wifiConfig        wificonfigs.Feature
	profileWifiConfig profilewificonfigs.Feature
	ieee              ieee8021xconfigs.Feature
	log               logger.Interface
}

var (
	ErrProfilesUseCase = consoleerrors.CreateConsoleError("ProfilesUseCase")
	ErrDatabase        = sqldb.DatabaseError{Console: consoleerrors.CreateConsoleError("ProfilesUseCase")}
	ErrNotFound        = sqldb.NotFoundError{Console: consoleerrors.CreateConsoleError("ProfilesUseCase")}
	ErrNotValid        = dto.NotValidError{Console: consoleerrors.CreateConsoleError("ProfilesUseCase")}
)

// New -.
func New(r Repository, wifiConfig wificonfigs.Feature, w profilewificonfigs.Feature, i ieee8021xconfigs.Feature, log logger.Interface) *UseCase {
	return &UseCase{
		repo:              r,
		wifiConfig:        wifiConfig,
		profileWifiConfig: w,
		ieee:              i,
		log:               log,
	}
}

// History - getting translate history from store.
func (uc *UseCase) GetCount(ctx context.Context, tenantID string) (int, error) {
	count, err := uc.repo.GetCount(ctx, tenantID)
	if err != nil {
		return 0, ErrDatabase.Wrap("Count", "uc.repo.GetCount", err)
	}

	return count, nil
}

func (uc *UseCase) Get(ctx context.Context, top, skip int, tenantID string) ([]dto.Profile, error) {
	data, err := uc.repo.Get(ctx, top, skip, tenantID)
	if err != nil {
		return nil, ErrDatabase.Wrap("Get", "uc.repo.Get", err)
	}

	if data == nil {
		return nil, ErrNotFound
	}

	// iterate over the data and convert each entity to dto
	d1 := make([]dto.Profile, len(data))

	for i := range data {
		tmpEntity := data[i] // create a new variable to avoid memory aliasing
		d1[i] = *uc.entityToDTO(&tmpEntity)
		associatedWiFiProfiles, _ := uc.profileWifiConfig.GetByProfileName(ctx, d1[i].ProfileName, tenantID)

		if len(associatedWiFiProfiles) > 0 {
			d1[i].WiFiConfigs = associatedWiFiProfiles
		}
	}

	return d1, nil
}

func (uc *UseCase) GetByName(ctx context.Context, profileName, tenantID string) (*dto.Profile, error) {
	data, err := uc.repo.GetByName(ctx, profileName, tenantID)
	if err != nil {
		return nil, ErrDatabase.Wrap("GetByName", "uc.repo.GetByName", err)
	}

	if data == nil {
		return nil, ErrNotFound
	}

	associatedWiFiProfiles, _ := uc.profileWifiConfig.GetByProfileName(ctx, profileName, tenantID)

	d2 := uc.entityToDTO(data)

	if len(associatedWiFiProfiles) > 0 {
		d2.WiFiConfigs = associatedWiFiProfiles
	}

	return d2, nil
}

func (uc *UseCase) Delete(ctx context.Context, profileName, tenantID string) error {
	// remove all wifi configs associated with the profile
	err := uc.profileWifiConfig.DeleteByProfileName(ctx, profileName, tenantID)
	if err != nil {
		return ErrDatabase.Wrap("Delete", "uc.repo.Delete", err)
	}

	isSuccessful, err := uc.repo.Delete(ctx, profileName, tenantID)
	if err != nil {
		return ErrDatabase.Wrap("Delete", "uc.repo.Delete", err)
	}

	if !isSuccessful {
		return ErrNotFound
	}

	return nil
}

func (uc *UseCase) isWifiProfileExists(ctx context.Context, d *dto.Profile, action string) error {
	if len(d.WiFiConfigs) > 0 {
		// check if the wireless profile is exists in the database
		wifiProfiles := []string{}

		for _, wifiConfig := range d.WiFiConfigs {
			result, err := uc.wifiConfig.CheckProfileExists(ctx, wifiConfig.WirelessProfileName, d.TenantID)
			if err != nil {
				return err
			}

			if !result {
				wifiProfiles = append(wifiProfiles, wifiConfig.WirelessProfileName)
			}
		}

		if len(wifiProfiles) > 0 {
			return ErrNotValid.Wrap(action, "uc.wifiConfig.CheckProfileExists", consoleerrors.CreateConsoleError("wifiProfiles are not found in the database"))
		}
	}

	return nil
}

func (uc *UseCase) Update(ctx context.Context, d *dto.Profile) (*dto.Profile, error) {
	d1 := uc.dtoToEntity(d)

	err := uc.isWifiProfileExists(ctx, d, "update")
	if err != nil {
		return nil, err
	}

	updated, err := uc.repo.Update(ctx, d1)
	if err != nil {
		return nil, ErrDatabase.Wrap("Update", "uc.repo.Update", err)
	}

	if !updated {
		return nil, ErrNotFound
	}
	// remove all wifi configs associated with the profile
	err = uc.profileWifiConfig.DeleteByProfileName(ctx, d.ProfileName, d.TenantID)
	if err != nil {
		return nil, ErrDatabase.Wrap("Delete", "uc.wifiRepo.DeleteByProfileName", err)
	}

	if d.DHCPEnabled {
		// insert new wifi configs
		if len(d.WiFiConfigs) > 0 {
			for _, wifiConfig := range d.WiFiConfigs {
				wifiConfig.ProfileName = d1.ProfileName

				tmpWifiConfig := wifiConfig // create a new variable to avoid memory aliasing

				err = uc.profileWifiConfig.Insert(ctx, &tmpWifiConfig)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	updatedProfile, err := uc.repo.GetByName(ctx, d.ProfileName, d.TenantID)
	if err != nil {
		return nil, err
	}

	d2 := uc.entityToDTO(updatedProfile)
	d2.WiFiConfigs = d.WiFiConfigs

	return d2, nil
}

func (uc *UseCase) Insert(ctx context.Context, d *dto.Profile) (*dto.Profile, error) {
	d1 := uc.dtoToEntity(d)

	if err := uc.isWifiProfileExists(ctx, d, "insert"); err != nil {
		return nil, err
	}

	if err := uc.validateIEEE8021xProfile(ctx, d1); err != nil {
		return nil, err
	}

	if err := uc.insertProfile(ctx, d1); err != nil {
		return nil, err
	}

	if err := uc.insertProfileWifiConfigs(ctx, d); err != nil {
		return nil, err
	}

	return uc.createdProfile(ctx, d)
}

func (uc *UseCase) validateIEEE8021xProfile(ctx context.Context, d1 *entity.Profile) error {
	if d1.IEEE8021xProfileName == nil || *d1.IEEE8021xProfileName == "" {
		return nil
	}

	return uc.checkIEEE8021xProfile(ctx, *d1.IEEE8021xProfileName, d1.TenantID)
}

func (uc *UseCase) checkIEEE8021xProfile(ctx context.Context, profileName, tenantID string) error {
	res, err := uc.ieee.GetByName(ctx, profileName, tenantID)
	if err != nil {
		var nfErr sqldb.NotFoundError
		if errors.As(err, &nfErr) {
			return ErrNotValid.Wrap("Insert", "uc.ieee.GetByName", consoleerrors.CreateConsoleError("IEEE profile is not found in the database"))
		}

		return err
	}

	if !res.WiredInterface {
		return ErrNotValid.Wrap("Insert", "uc.ieee.GetByName", consoleerrors.CreateConsoleError("Wired interface is required"))
	}

	return nil
}

func (uc *UseCase) insertProfile(ctx context.Context, d1 *entity.Profile) error {
	_, err := uc.repo.Insert(ctx, d1)

	return err
}

func (uc *UseCase) insertProfileWifiConfigs(ctx context.Context, d *dto.Profile) error {
	if len(d.WiFiConfigs) > 0 {
		for _, wifiConfig := range d.WiFiConfigs {
			wifiConfig.ProfileName = d.ProfileName
			tmpWifiConfig := wifiConfig // create a new variable to avoid memory aliasing

			err := uc.profileWifiConfig.Insert(ctx, &tmpWifiConfig)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (uc *UseCase) createdProfile(ctx context.Context, d *dto.Profile) (*dto.Profile, error) {
	newProfile, err := uc.repo.GetByName(ctx, d.ProfileName, d.TenantID)
	if err != nil {
		return nil, err
	}

	d2 := uc.entityToDTO(newProfile)
	d2.WiFiConfigs = d.WiFiConfigs

	return d2, nil
}

// convert dto.Profile to entity.Profile.
func (uc *UseCase) dtoToEntity(d *dto.Profile) *entity.Profile {
	// convert []string to comma separated string
	tags := strings.Join(d.Tags, ", ")

	d1 := &entity.Profile{
		ProfileName:                d.ProfileName,
		AMTPassword:                d.AMTPassword,
		CreationDate:               d.CreationDate,
		CreatedBy:                  d.CreatedBy,
		GenerateRandomPassword:     d.GenerateRandomPassword,
		CIRAConfigName:             d.CIRAConfigName,
		Activation:                 d.Activation,
		MEBXPassword:               d.MEBXPassword,
		GenerateRandomMEBxPassword: d.GenerateRandomMEBxPassword,
		Tags:                       tags,
		DHCPEnabled:                d.DHCPEnabled,
		IPSyncEnabled:              d.IPSyncEnabled,
		LocalWiFiSyncEnabled:       d.LocalWiFiSyncEnabled,
		TenantID:                   d.TenantID,
		TLSMode:                    d.TLSMode,
		TLSSigningAuthority:        d.TLSSigningAuthority,
		UserConsent:                d.UserConsent,
		IDEREnabled:                d.IDEREnabled,
		KVMEnabled:                 d.KVMEnabled,
		SOLEnabled:                 d.SOLEnabled,
		IEEE8021xProfileName:       d.IEEE8021xProfileName,
		Version:                    d.Version,
	}

	return d1
}

// convert entity.Profile to dto.Profile.
func (uc *UseCase) entityToDTO(d *entity.Profile) *dto.Profile {
	// convert comma separated string to []string
	tags := strings.Split(d.Tags, ",")
	d1 := &dto.Profile{
		ProfileName:                d.ProfileName,
		AMTPassword:                d.AMTPassword,
		CreationDate:               d.CreationDate,
		CreatedBy:                  d.CreatedBy,
		GenerateRandomPassword:     d.GenerateRandomPassword,
		CIRAConfigName:             d.CIRAConfigName,
		Activation:                 d.Activation,
		MEBXPassword:               d.MEBXPassword,
		GenerateRandomMEBxPassword: d.GenerateRandomMEBxPassword,
		Tags:                       tags,
		DHCPEnabled:                d.DHCPEnabled,
		IPSyncEnabled:              d.IPSyncEnabled,
		LocalWiFiSyncEnabled:       d.LocalWiFiSyncEnabled,
		TenantID:                   d.TenantID,
		TLSMode:                    d.TLSMode,
		TLSSigningAuthority:        d.TLSSigningAuthority,
		UserConsent:                d.UserConsent,
		IDEREnabled:                d.IDEREnabled,
		KVMEnabled:                 d.KVMEnabled,
		SOLEnabled:                 d.SOLEnabled,
		IEEE8021xProfileName:       d.IEEE8021xProfileName,
		Version:                    d.Version,
	}

	if d.IEEE8021xProfileName != nil && *d.IEEE8021xProfileName != "" {
		val := &dto.IEEE8021xConfig{
			ProfileName:            *d.IEEE8021xProfileName,
			AuthenticationProtocol: *d.AuthenticationProtocol,
			PXETimeout:             d.PXETimeout,
			WiredInterface:         *d.WiredInterface,
			TenantID:               d.TenantID,
			Version:                d.Version,
		}
		d1.IEEE8021xProfile = val
	}

	return d1
}
