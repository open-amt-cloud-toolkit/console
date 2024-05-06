package devices

import (
	"context"
	"fmt"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
)

var (
	ErrDomainsUseCase = consoleerrors.CreateConsoleError("DevicesUseCase")
	ErrDatabase       = consoleerrors.DatabaseError{Console: consoleerrors.CreateConsoleError("DevicesUseCase")}
	ErrNotFound       = consoleerrors.NotFoundError{Console: consoleerrors.CreateConsoleError("DevicesUseCase")}
)

// History - getting translate history from store.
func (uc *UseCase) GetCount(ctx context.Context, tenantID string) (int, error) {
	count, err := uc.repo.GetCount(ctx, tenantID)
	if err != nil {
		return 0, ErrDatabase.Wrap("Count", "uc.repo.GetCount", err)
	}

	return count, nil
}

func (uc *UseCase) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Device, error) {
	data, err := uc.repo.Get(ctx, top, skip, tenantID)
	if err != nil {
		return nil, ErrDatabase.Wrap("Get", "uc.repo.Get", err)
	}

	return data, nil
}

func (uc *UseCase) GetByID(ctx context.Context, guid, tenantID string) (*entity.Device, error) {
	data, err := uc.repo.GetByID(ctx, guid, tenantID)
	if err != nil {
		return nil, ErrDatabase.Wrap("GetByID", "uc.repo.GetByID", err)
	}

	if data == nil {
		return nil, ErrNotFound
	}

	return data, nil
}

func (uc *UseCase) GetDistinctTags(ctx context.Context, tenantID string) ([]string, error) {
	data, err := uc.repo.GetDistinctTags(ctx, tenantID)
	if err != nil {
		return nil, ErrDatabase.Wrap("GetDistinctTags", "uc.repo.GetDistinctTags", err)
	}

	return data, nil
}

func (uc *UseCase) GetByTags(ctx context.Context, tags []string, method string, limit, offset int, tenantID string) ([]entity.Device, error) {
	data, err := uc.repo.GetByTags(ctx, tags, method, limit, offset, tenantID)
	if err != nil {
		return nil, fmt.Errorf("DevicesUseCase - GetByTags - uc.repo.GetByTags: %w", err)
	}

	return data, nil
}

func (uc *UseCase) Delete(ctx context.Context, guid, tenantID string) error {
	isSuccessful, err := uc.repo.Delete(ctx, guid, tenantID)
	if err != nil {
		return ErrDatabase.Wrap("Delete", "uc.repo.Delete", err)
	}

	if !isSuccessful {
		return ErrNotFound
	}

	return nil
}

func (uc *UseCase) Update(ctx context.Context, d *entity.Device) (*entity.Device, error) {
	_, err := uc.repo.Update(ctx, d)
	if err != nil {
		return nil, ErrDatabase.Wrap("Update", "uc.repo.Update", err)
	}

	updateDevice, err := uc.repo.GetByID(ctx, d.GUID, d.TenantID)
	if err != nil {
		return nil, err
	}

	return updateDevice, nil
}

func (uc *UseCase) Insert(ctx context.Context, d *entity.Device) (*entity.Device, error) {
	_, err := uc.repo.Insert(ctx, d)
	if err != nil {
		return nil, ErrDatabase.Wrap("Insert", "uc.repo.Insert", err)
	}

	newDevice, err := uc.repo.GetByID(ctx, d.GUID, d.TenantID)
	if err != nil {
		return nil, err
	}

	return newDevice, nil
}
