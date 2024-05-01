package domains

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
		return 0, fmt.Errorf("DomainsUseCase - Count - s.repo.GetCount: %w", err)
	}

	return count, nil
}

func (uc *UseCase) Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Domain, error) {
	data, err := uc.repo.Get(ctx, top, skip, tenantID)
	if err != nil {
		return nil, fmt.Errorf("DomainsUseCase - Get - s.repo.Get: %w", err)
	}

	return data, nil
}

func (uc *UseCase) GetDomainByDomainSuffix(ctx context.Context, domainSuffix, tenantID string) (*entity.Domain, error) {
	data, err := uc.repo.GetDomainByDomainSuffix(ctx, domainSuffix, tenantID)
	if err != nil {
		return nil, fmt.Errorf("DomainsUseCase - GetDomainByDomainSuffix - s.repo.GetDomainByDomainSuffix: %w", err)
	}

	return data, nil
}

func (uc *UseCase) GetByName(ctx context.Context, domainName, tenantID string) (*entity.Domain, error) {
	data, err := uc.repo.GetByName(ctx, domainName, tenantID)
	if err != nil {
		return nil, fmt.Errorf("DomainsUseCase - GetByName - s.repo.GetByName: %w", err)
	}

	return data, nil
}

func (uc *UseCase) Delete(ctx context.Context, domainName, tenantID string) (bool, error) {
	data, err := uc.repo.Delete(ctx, domainName, tenantID)
	if err != nil {
		return false, fmt.Errorf("DomainsUseCase - Delete - s.repo.Delete: %w", err)
	}

	return data, nil
}

func (uc *UseCase) Update(ctx context.Context, d *entity.Domain) (bool, error) {
	data, err := uc.repo.Update(ctx, d)
	if err != nil {
		return false, fmt.Errorf("DomainsUseCase - Update - s.repo.Update: %w", err)
	}

	return data, nil
}

func (uc *UseCase) Insert(ctx context.Context, d *entity.Domain) (string, error) {
	data, err := uc.repo.Insert(ctx, d)
	if err != nil {
		return "", fmt.Errorf("DomainsUseCase - Insert - s.repo.Insert: %w", err)
	}

	return data, nil
}
