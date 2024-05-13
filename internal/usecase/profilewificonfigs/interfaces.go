package profilewificonfigs

import (
	"context"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
)

type (
	Repository interface {
		GetByProfileName(ctx context.Context, profileName, tenantID string) ([]entity.ProfileWiFiConfigs, error)
		DeleteByProfileName(ctx context.Context, profileName, tenantID string) (bool, error)
		Insert(ctx context.Context, p *entity.ProfileWiFiConfigs) (string, error)
	}

	Feature interface {
		GetByProfileName(ctx context.Context, profileName, tenantID string) ([]dto.ProfileWiFiConfigs, error)
		DeleteByProfileName(ctx context.Context, profileName, tenantID string) error
		Insert(ctx context.Context, p *dto.ProfileWiFiConfigs) error
	}
)
