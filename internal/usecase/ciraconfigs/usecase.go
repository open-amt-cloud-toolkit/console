package ciraconfigs

import (
	"context"
	"fmt"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

// UseCase -.
type UseCase struct {
	repo Repository
	log  logger.Interface
}

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
		return 0, fmt.Errorf("CIRAConfigsUseCase - Count - s.repo.GetCount: %w", err)
	}

	return count, nil
}

func (uc *UseCase) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.CIRAConfig, error) {
	data, err := uc.repo.Get(ctx, top, skip, tenantID)
	if err != nil {
		return nil, fmt.Errorf("CIRAConfigsUseCase - Get - s.repo.Get: %w", err)
	}

	return data, nil
}

func (uc *UseCase) GetByName(ctx context.Context, configName, tenantID string) (entity.CIRAConfig, error) {
	data, err := uc.repo.GetByName(ctx, configName, tenantID)
	if err != nil {
		return entity.CIRAConfig{}, fmt.Errorf("CIRAConfigsUseCase - GetByName - s.repo.GetByName: %w", err)
	}

	return data, nil
}

func (uc *UseCase) Delete(ctx context.Context, configName, tenantID string) (bool, error) {
	data, err := uc.repo.Delete(ctx, configName, tenantID)
	if err != nil {
		return false, fmt.Errorf("CIRAConfigsUseCase - Delete - s.repo.Delete: %w", err)
	}

	return data, nil
}

func (uc *UseCase) Update(ctx context.Context, d *entity.CIRAConfig) (bool, error) {
	data, err := uc.repo.Update(ctx, d)
	if err != nil {
		return false, fmt.Errorf("CIRAConfigsUseCase - Update - s.repo.Update: %w", err)
	}

	return data, nil
}

func (uc *UseCase) Insert(ctx context.Context, d *entity.CIRAConfig) (string, error) {
	data, err := uc.repo.Insert(ctx, d)
	if err != nil {
		return "", fmt.Errorf("CIRAConfigsUseCase - Insert - s.repo.Insert: %w", err)
	}

	return data, nil
}
