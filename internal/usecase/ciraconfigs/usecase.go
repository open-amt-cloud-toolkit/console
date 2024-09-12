package ciraconfigs

import (
	"context"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	dtov1 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
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
	ErrCIRAConfigsUseCase = consoleerrors.CreateConsoleError("CIRAConfigsUseCase")
	ErrDatabase           = sqldb.DatabaseError{Console: consoleerrors.CreateConsoleError("CIRAConfigsUseCase")}
	ErrNotFound           = sqldb.NotFoundError{Console: consoleerrors.CreateConsoleError("CIRAConfigsUseCase")}
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

func (uc *UseCase) Get(ctx context.Context, top, skip int, tenantID string) ([]dtov1.CIRAConfig, error) {
	data, err := uc.repo.Get(ctx, top, skip, tenantID)
	if err != nil {
		return nil, ErrDatabase.Wrap("Get", "uc.repo.Get", err)
	}
	// iterate over the data and convert each entity to dto
	d1 := make([]dtov1.CIRAConfig, len(data))

	for i := range data {
		tmpEntity := data[i] // create a new variable to avoid memory aliasing
		d1[i] = *uc.entityToDTO(&tmpEntity)
	}

	return d1, nil
}

func (uc *UseCase) GetByName(ctx context.Context, configName, tenantID string) (*dtov1.CIRAConfig, error) {
	data, err := uc.repo.GetByName(ctx, configName, tenantID)
	if err != nil {
		return nil, ErrDatabase.Wrap("GetByName", "uc.repo.GetByName", err)
	}

	if data == nil {
		return nil, ErrNotFound
	}

	d2 := uc.entityToDTO(data)

	return d2, nil
}

func (uc *UseCase) Delete(ctx context.Context, configName, tenantID string) error {
	isSuccessful, err := uc.repo.Delete(ctx, configName, tenantID)
	if err != nil {
		return ErrDatabase.Wrap("Delete", "uc.repo.Delete", err)
	}

	if !isSuccessful {
		return ErrNotFound
	}

	return nil
}

func (uc *UseCase) Update(ctx context.Context, d *dtov1.CIRAConfig) (*dtov1.CIRAConfig, error) {
	d1 := uc.dtoToEntity(d)

	updated, err := uc.repo.Update(ctx, d1)
	if err != nil {
		return nil, ErrDatabase.Wrap("Update", "uc.repo.Update", err)
	}

	if !updated {
		return nil, ErrNotFound
	}

	updatedCiraConfig, err := uc.repo.GetByName(ctx, d.ConfigName, d.TenantID)
	if err != nil {
		return nil, err
	}

	d2 := uc.entityToDTO(updatedCiraConfig)

	return d2, nil
}

func (uc *UseCase) Insert(ctx context.Context, d *dtov1.CIRAConfig) (*dtov1.CIRAConfig, error) {
	d1 := uc.dtoToEntity(d)

	_, err := uc.repo.Insert(ctx, d1)
	if err != nil {
		return nil, ErrDatabase.Wrap("Insert", "uc.repo.Insert", err)
	}

	newConfig, err := uc.repo.GetByName(ctx, d.ConfigName, d.TenantID)
	if err != nil {
		return nil, err
	}

	d2 := uc.entityToDTO(newConfig)

	return d2, nil
}

// convert dtov1.CIRAConfig to entity.CIRAConfig.
func (uc *UseCase) dtoToEntity(d *dtov1.CIRAConfig) *entity.CIRAConfig {
	d1 := &entity.CIRAConfig{
		ConfigName:          d.ConfigName,
		MPSAddress:          d.MPSAddress,
		MPSPort:             d.MPSPort,
		Username:            d.Username,
		Password:            d.Password,
		CommonName:          d.CommonName,
		ServerAddressFormat: d.ServerAddressFormat,
		AuthMethod:          d.AuthMethod,
		MPSRootCertificate:  d.MPSRootCertificate,
		ProxyDetails:        d.ProxyDetails,
		TenantID:            d.TenantID,
		RegeneratePassword:  d.RegeneratePassword,
		Version:             d.Version,
	}

	return d1
}

// convert entity.CIRAConfig to dtov1.CIRAConfig.
func (uc *UseCase) entityToDTO(d *entity.CIRAConfig) *dtov1.CIRAConfig {
	d1 := &dtov1.CIRAConfig{
		ConfigName:          d.ConfigName,
		MPSAddress:          d.MPSAddress,
		MPSPort:             d.MPSPort,
		Username:            d.Username,
		Password:            d.Password,
		CommonName:          d.CommonName,
		ServerAddressFormat: d.ServerAddressFormat,
		AuthMethod:          d.AuthMethod,
		MPSRootCertificate:  d.MPSRootCertificate,
		ProxyDetails:        d.ProxyDetails,
		TenantID:            d.TenantID,
		RegeneratePassword:  d.RegeneratePassword,
		Version:             d.Version,
	}

	return d1
}
