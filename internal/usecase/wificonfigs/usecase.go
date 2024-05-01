package wificonfigs

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
func (uc *UseCase) CheckProfileExists(ctx context.Context, profileName, tenantID string) (bool, error) {
	data, err := uc.repo.CheckProfileExists(ctx, profileName, tenantID)
	if err != nil {
		return false, fmt.Errorf("WificonfigsUseCase - Count - s.repo.GetCount: %w", err)
	}

	return data, nil
}

func (uc *UseCase) GetCount(ctx context.Context, tenantID string) (int, error) {
	count, err := uc.repo.GetCount(ctx, tenantID)
	if err != nil {
		return 0, fmt.Errorf("WificonfigsUseCase - Count - s.repo.GetCount: %w", err)
	}

	return count, nil
}

func (uc *UseCase) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.WirelessConfig, error) {
	data, err := uc.repo.Get(ctx, top, skip, tenantID)
	if err != nil {
		return nil, fmt.Errorf("WificonfigsUseCase - Get - s.repo.Get: %w", err)
	}

	return data, nil
}

func (uc *UseCase) GetByName(ctx context.Context, profileName, tenantID string) (entity.WirelessConfig, error) {
	data, err := uc.repo.GetByName(ctx, profileName, tenantID)
	if err != nil {
		return entity.WirelessConfig{}, fmt.Errorf("WificonfigsUseCase - GetByName - s.repo.GetByName: %w", err)
	}

	return data, nil
}

func (uc *UseCase) Delete(ctx context.Context, profileName, tenantID string) (bool, error) {
	data, err := uc.repo.Delete(ctx, profileName, tenantID)
	if err != nil {
		return false, fmt.Errorf("WificonfigsUseCase - Delete - s.repo.Delete: %w", err)
	}

	return data, nil
}

func (uc *UseCase) Update(ctx context.Context, d *entity.WirelessConfig) (bool, error) {
	data, err := uc.repo.Update(ctx, d)
	if err != nil {
		return false, fmt.Errorf("WificonfigsUseCase - Update - s.repo.Update: %w", err)
	}

	return data, nil
}

func (uc *UseCase) Insert(ctx context.Context, d *entity.WirelessConfig) (string, error) {
	data, err := uc.repo.Insert(ctx, d)
	if err != nil {
		return "", fmt.Errorf("WificonfigsUseCase - Insert - s.repo.Insert: %w", err)
	}

	return data, nil
}
