package devices

import (
	"context"
	"fmt"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
)

// History - getting translate history from store.
func (uc *UseCase) GetCount(ctx context.Context, tenantID string) (int, error) {
	count, err := uc.repo.GetCount(ctx, tenantID)
	if err != nil {
		return 0, fmt.Errorf("DevicesUseCase - Count - s.repo.GetCount: %w", err)
	}

	return count, nil
}

func (uc *UseCase) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Device, error) {
	data, err := uc.repo.Get(ctx, top, skip, tenantID)
	if err != nil {
		return nil, fmt.Errorf("DevicesUseCase - Get - s.repo.Get: %w", err)
	}

	return data, nil
}

func (uc *UseCase) GetByID(ctx context.Context, guid, tenantID string) (entity.Device, error) {
	data, err := uc.repo.GetByID(ctx, guid, tenantID)
	if err != nil {
		return entity.Device{}, fmt.Errorf("DevicesUseCase - GetByID - s.repo.GetByID: %w", err)
	}

	return data, nil
}

func (uc *UseCase) GetDistinctTags(ctx context.Context, tenantID string) ([]string, error) {
	data, err := uc.repo.GetDistinctTags(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("DevicesUseCase - GetDistinctTags - s.repo.GetDistinctTags: %w", err)
	}

	return data, nil
}

func (uc *UseCase) GetByTags(ctx context.Context, tags []string, method string, limit, offset int, tenantID string) ([]entity.Device, error) {
	data, err := uc.repo.GetByTags(ctx, tags, method, limit, offset, tenantID)
	if err != nil {
		return nil, fmt.Errorf("DevicesUseCase - GetByTags - s.repo.GetByTags: %w", err)
	}

	return data, nil
}

func (uc *UseCase) Delete(ctx context.Context, guid, tenantID string) (bool, error) {
	data, err := uc.repo.Delete(ctx, guid, tenantID)
	if err != nil {
		return false, fmt.Errorf("DevicesUseCase - Delete - s.repo.Delete: %w", err)
	}

	return data, nil
}

func (uc *UseCase) Update(ctx context.Context, d *entity.Device) (bool, error) {
	data, err := uc.repo.Update(ctx, d)
	if err != nil {
		return false, fmt.Errorf("DevicesUseCase - Update - s.repo.Update: %w", err)
	}

	return data, nil
}

func (uc *UseCase) Insert(ctx context.Context, d *entity.Device) (string, error) {
	data, err := uc.repo.Insert(ctx, d)
	if err != nil {
		return "", fmt.Errorf("DevicesUseCase - Insert - s.repo.Insert: %w", err)
	}

	return data, nil
}
