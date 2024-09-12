package profilewificonfigs

import (
	"context"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	dtov1 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
)

type (
	Repository interface {
		GetByProfileName(ctx context.Context, profileName, tenantID string) ([]entity.ProfileWiFiConfigs, error)
		DeleteByProfileName(ctx context.Context, profileName, tenantID string) (bool, error)
		Insert(ctx context.Context, p *entity.ProfileWiFiConfigs) (string, error)
	}

	Feature interface {
		GetByProfileName(ctx context.Context, profileName, tenantID string) ([]dtov1.ProfileWiFiConfigs, error)
		DeleteByProfileName(ctx context.Context, profileName, tenantID string) error
		Insert(ctx context.Context, p *dtov1.ProfileWiFiConfigs) error
	}
)
