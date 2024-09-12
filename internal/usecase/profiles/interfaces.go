package profiles

import (
	"context"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	dtov1 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
)

type (
	Repository interface {
		GetCount(ctx context.Context, tenantID string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Profile, error)
		GetByName(ctx context.Context, profileName, tenantID string) (*entity.Profile, error)
		Delete(ctx context.Context, profileName, tenantID string) (bool, error)
		Update(ctx context.Context, p *entity.Profile) (bool, error)
		Insert(ctx context.Context, p *entity.Profile) (string, error)
	}

	Feature interface {
		GetCount(ctx context.Context, tenantID string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]dtov1.Profile, error)
		GetByName(ctx context.Context, profileName, tenantID string) (*dtov1.Profile, error)
		Delete(ctx context.Context, profileName, tenantID string) error
		Update(ctx context.Context, p *dtov1.Profile) (*dtov1.Profile, error)
		Insert(ctx context.Context, p *dtov1.Profile) (*dtov1.Profile, error)
	}
)
