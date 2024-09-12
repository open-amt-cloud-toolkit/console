package ciraconfigs

import (
	"context"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	dtov1 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
)

type (
	Repository interface {
		GetCount(ctx context.Context, tenantID string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.CIRAConfig, error)
		GetByName(ctx context.Context, configName, tenantID string) (*entity.CIRAConfig, error)
		Delete(ctx context.Context, profileName, tenantID string) (bool, error)
		Update(ctx context.Context, p *entity.CIRAConfig) (bool, error)
		Insert(ctx context.Context, p *entity.CIRAConfig) (string, error)
	}
	Feature interface {
		GetCount(ctx context.Context, tenantID string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]dtov1.CIRAConfig, error)
		GetByName(ctx context.Context, configName, tenantID string) (*dtov1.CIRAConfig, error)
		Delete(ctx context.Context, profileName, tenantID string) error
		Update(ctx context.Context, p *dtov1.CIRAConfig) (*dtov1.CIRAConfig, error)
		Insert(ctx context.Context, p *dtov1.CIRAConfig) (*dtov1.CIRAConfig, error)
	}
)
