package profiles

import (
	"context"
	"strings"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/sqldb"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

// UseCase -.
type UseCase struct {
	repo Repository
	log  logger.Interface
}

var (
	ErrDomainsUseCase = consoleerrors.CreateConsoleError("ProfilesUseCase")
	ErrDatabase       = sqldb.DatabaseError{Console: consoleerrors.CreateConsoleError("ProfilesUseCase")}
	ErrNotFound       = sqldb.NotFoundError{Console: consoleerrors.CreateConsoleError("ProfilesUseCase")}
)

// New -.
func New(r Repository, log logger.Interface) *UseCase {
	return &UseCase{
		repo: r,
		log:  log,
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

	d2 := uc.entityToDTO(data)

	return d2, nil
}

func (uc *UseCase) Delete(ctx context.Context, profileName, tenantID string) error {
	isSuccessful, err := uc.repo.Delete(ctx, profileName, tenantID)
	if err != nil {
		return ErrDatabase.Wrap("Delete", "uc.repo.Delete", err)
	}

	if !isSuccessful {
		return ErrNotFound
	}

	return nil
}

func (uc *UseCase) Update(ctx context.Context, d *dto.Profile) (*dto.Profile, error) {
	d1 := uc.dtoToEntity(d)

	_, err := uc.repo.Update(ctx, d1)
	if err != nil {
		return nil, ErrDatabase.Wrap("Update", "uc.repo.Update", err)
	}

	updatedProfile, err := uc.repo.GetByName(ctx, d.ProfileName, d.TenantID)
	if err != nil {
		return nil, err
	}

	d2 := uc.entityToDTO(updatedProfile)

	return d2, nil
}

func (uc *UseCase) Insert(ctx context.Context, d *dto.Profile) (*dto.Profile, error) {
	d1 := uc.dtoToEntity(d)

	_, err := uc.repo.Insert(ctx, d1)
	if err != nil {
		return nil, ErrDatabase.Wrap("Insert", "uc.repo.Insert", err)
	}

	newProfile, err := uc.repo.GetByName(ctx, d.ProfileName, d.TenantID)
	if err != nil {
		return nil, err
	}

	d2 := uc.entityToDTO(newProfile)

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

	return d1
}
