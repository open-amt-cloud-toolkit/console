package ciraconfigs

import (
	"context"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

// UseCase -.
type UseCase struct {
	repo Repository
	log  logger.Interface
}

var (
	ErrDomainsUseCase = consoleerrors.CreateConsoleError("CIRAConfigsUseCase")
	ErrDatabase       = consoleerrors.DatabaseError{Console: consoleerrors.CreateConsoleError("CIRAConfigsUseCase")}
	ErrNotFound       = consoleerrors.NotFoundError{Console: consoleerrors.CreateConsoleError("CIRAConfigsUseCase")}
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

func (uc *UseCase) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.CIRAConfig, error) {
	data, err := uc.repo.Get(ctx, top, skip, tenantID)
	if err != nil {
		return nil, ErrDatabase.Wrap("Get", "uc.repo.Get", err)
	}

	return data, nil
}

func (uc *UseCase) GetByName(ctx context.Context, configName, tenantID string) (*entity.CIRAConfig, error) {
	data, err := uc.repo.GetByName(ctx, configName, tenantID)
	if err != nil {
		return nil, ErrDatabase.Wrap("GetByName", "uc.repo.GetByName", err)
	}

	if data == nil {
		return nil, ErrNotFound
	}

	return data, nil
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

func (uc *UseCase) Update(ctx context.Context, d *entity.CIRAConfig) (*entity.CIRAConfig, error) {
	_, err := uc.repo.Update(ctx, d)
	if err != nil {
		return nil, ErrDatabase.Wrap("Update", "uc.repo.Update", err)
	}

	updatedCiraConfig, err := uc.repo.GetByName(ctx, d.ConfigName, d.TenantID)
	if err != nil {
		return nil, err
	}

	return updatedCiraConfig, nil
}

func (uc *UseCase) Insert(ctx context.Context, d *entity.CIRAConfig) (*entity.CIRAConfig, error) {
	_, err := uc.repo.Insert(ctx, d)
	if err != nil {
		return nil, ErrDatabase.Wrap("Insert", "uc.repo.Insert", err)
	}

	newConfig, err := uc.repo.GetByName(ctx, d.ConfigName, d.TenantID)
	if err != nil {
		return nil, err
	}

	return newConfig, nil
}
