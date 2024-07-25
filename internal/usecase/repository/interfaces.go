package repository

import (
	"context"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	wsmanAPI "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
)

type (
	Repository interface {
		GetCount(context.Context, string) (int, error)
		Get(ctx context.Context, top, skip int, tenantID string) ([]entity.Device, error)
		GetByID(ctx context.Context, guid, tenantID string) (*entity.Device, error)
		GetDistinctTags(ctx context.Context, tenantID string) ([]string, error)
		GetByTags(ctx context.Context, tags []string, method string, limit, offset int, tenantID string) ([]entity.Device, error)
		Delete(ctx context.Context, guid, tenantID string) (bool, error)
		Update(ctx context.Context, d *entity.Device) (bool, error)
		Insert(ctx context.Context, d *entity.Device) (string, error)
		GetByColumn(ctx context.Context, columnName, queryValue, tenantID string) ([]entity.Device, error)
	}

	WSMAN interface {
		SetupWsmanClient(device dto.Device, isRedirection, logMessages bool) wsmanAPI.Management
		DestroyWsmanClient(device dto.Device)
		Worker()
	}
)
