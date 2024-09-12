package profilewificonfigs

import (
	"context"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	dtov1 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/sqldb"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

type UseCase struct {
	repo Repository
	log  logger.Interface
}

var (
	ErrProfileWiFiConfigsUseCase = consoleerrors.CreateConsoleError("ProfilesWiFiUseCase")
	ErrDatabase                  = sqldb.DatabaseError{Console: consoleerrors.CreateConsoleError("ProfilesWiFiUseCase")}
	ErrNotFound                  = sqldb.NotFoundError{Console: consoleerrors.CreateConsoleError("ProfilesWiFiUseCase")}
)

func New(r Repository, log logger.Interface) *UseCase {
	return &UseCase{
		repo: r,
		log:  log,
	}
}

func (uc *UseCase) GetByProfileName(ctx context.Context, profileName, tenantID string) ([]dtov1.ProfileWiFiConfigs, error) {
	data, err := uc.repo.GetByProfileName(ctx, profileName, tenantID)
	if err != nil {
		return nil, ErrDatabase.Wrap("Get", "uc.repo.Get", err)
	}

	if data == nil {
		return nil, ErrNotFound
	}

	// iterate over the data and convert each entity to dto
	d1 := make([]dtov1.ProfileWiFiConfigs, len(data))

	for i := range data {
		tmpEntity := data[i] // create a new variable to avoid memory aliasing
		d1[i] = *uc.entityToDTO(&tmpEntity)
	}

	return d1, nil
}

func (uc *UseCase) DeleteByProfileName(ctx context.Context, profileName, tenantID string) error {
	_, err := uc.repo.DeleteByProfileName(ctx, profileName, tenantID)
	if err != nil {
		return ErrDatabase.Wrap("Delete", "uc.repo.Delete", err)
	}

	return nil
}

func (uc *UseCase) Insert(ctx context.Context, d *dtov1.ProfileWiFiConfigs) error {
	d1 := uc.dtoToEntity(d)

	_, err := uc.repo.Insert(ctx, d1)
	if err != nil {
		return ErrDatabase.Wrap("Insert", "uc.repo.Insert", err)
	}

	return nil
}

// convert dtov1.ProfileWiFiConfigs to entity.ProfileWiFiConfigs.
func (uc *UseCase) dtoToEntity(d *dtov1.ProfileWiFiConfigs) *entity.ProfileWiFiConfigs {
	d1 := &entity.ProfileWiFiConfigs{
		Priority:            d.Priority,
		WirelessProfileName: d.WirelessProfileName,
		ProfileName:         d.ProfileName,
		TenantID:            d.TenantID,
	}

	return d1
}

// convert entity.ProfileWiFiConfigs to dtov1.ProfileWiFiConfigs.
func (uc *UseCase) entityToDTO(d *entity.ProfileWiFiConfigs) *dtov1.ProfileWiFiConfigs {
	d1 := &dtov1.ProfileWiFiConfigs{
		Priority:            d.Priority,
		WirelessProfileName: d.WirelessProfileName,
		ProfileName:         d.ProfileName,
		TenantID:            d.TenantID,
	}

	return d1
}
