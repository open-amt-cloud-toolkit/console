package devices

import (
	"strings"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

const (
	// MinAMTVersion - minimum AMT version required for certain features in power capabilities.
	MinAMTVersion = 9
)

// UseCase -.
type UseCase struct {
	repo   Repository
	device WSMAN
	log    logger.Interface
}

var ErrAMT = AMTError{Console: consoleerrors.CreateConsoleError("DevicesUseCase")}

// New -.
func New(r Repository, d WSMAN, log logger.Interface) *UseCase {
	uc := &UseCase{
		repo:   r,
		device: d,
		log:    log,
	}
	// start up the worker
	go d.Worker()

	return uc
}

// convert dto.Device to entity.Device.
func (uc *UseCase) dtoToEntity(d *dto.Device) *entity.Device {
	// convert []string to comma separated string
	if d.Tags == nil {
		d.Tags = []string{}
	}

	tags := strings.Join(d.Tags, ", ")

	d1 := &entity.Device{
		ConnectionStatus: d.ConnectionStatus,
		MPSInstance:      d.MPSInstance,
		Hostname:         d.Hostname,
		GUID:             d.GUID,
		MPSUsername:      d.MPSUsername,
		Tags:             tags,
		TenantID:         d.TenantID,
		FriendlyName:     d.FriendlyName,
		DNSSuffix:        d.DNSSuffix,
		LastConnected:    d.LastConnected,
		LastSeen:         d.LastSeen,
		LastDisconnected: d.LastDisconnected,
		// DeviceInfo:       d.DeviceInfo,
		Username:        d.Username,
		Password:        d.Password,
		UseTLS:          d.UseTLS,
		AllowSelfSigned: d.AllowSelfSigned,
	}

	return d1
}

// convert entity.Device to dto.Device.
func (uc *UseCase) entityToDTO(d *entity.Device) *dto.Device {
	// convert comma separated string to []string
	tags := strings.Split(d.Tags, ",")

	d1 := &dto.Device{
		ConnectionStatus: d.ConnectionStatus,
		MPSInstance:      d.MPSInstance,
		Hostname:         d.Hostname,
		GUID:             d.GUID,
		MPSUsername:      d.MPSUsername,
		Tags:             tags,
		TenantID:         d.TenantID,
		FriendlyName:     d.FriendlyName,
		DNSSuffix:        d.DNSSuffix,
		LastConnected:    d.LastConnected,
		LastSeen:         d.LastSeen,
		LastDisconnected: d.LastDisconnected,
		// DeviceInfo:       d.DeviceInfo,
		Username:        d.Username,
		Password:        d.Password,
		UseTLS:          d.UseTLS,
		AllowSelfSigned: d.AllowSelfSigned,
	}

	return d1
}
