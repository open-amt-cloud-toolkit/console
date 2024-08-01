package redirection

import (
	"context"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman"
)

type (
	Repository interface {
		GetByID(ctx context.Context, guid, tenantID string) (*entity.Device, error)
	}

	Feature interface {
		// SetupWsmanClient(device dto.Device, isRedirection, logMessages bool) wsman.Messages
		RedirectConnect(ctx context.Context, deviceConnection *DeviceConnection) error
		RedirectClose(ctx context.Context, deviceConnection *DeviceConnection) error
		RedirectListen(ctx context.Context, deviceConnection *DeviceConnection) ([]byte, error)
		RedirectSend(ctx context.Context, deviceConnection *DeviceConnection, message []byte) error
	}

	WSMAN interface {
		SetupWsmanClient(device dto.Device, logMessages bool) Feature
		DestroyWsmanClient(device dto.Device)
	}
)
