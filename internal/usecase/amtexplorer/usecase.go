package amtexplorer

import (
	"strings"

	"github.com/open-amt-cloud-toolkit/console/internal/entity"
	dtov1 "github.com/open-amt-cloud-toolkit/console/internal/entity/dto/v1"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/sqldb"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

var ErrDatabase = sqldb.DatabaseError{Console: consoleerrors.CreateConsoleError("DevicesUseCase")}

// UseCase -.
type UseCase struct {
	repo   Repository
	device WSMAN
	log    logger.Interface
}

var ErrAMT = AMTError{Console: consoleerrors.CreateConsoleError("DevicesUseCase")}

// New -.
func New(r Repository, d WSMAN, log logger.Interface) *UseCase {
	return &UseCase{
		repo:   r,
		device: d,
		log:    log,
	}
}

// convert entity.Device to dtov1.Device.
func (uc *UseCase) entityToDTO(d *entity.Device) *dtov1.Device {
	// convert comma separated string to []string
	tags := strings.Split(d.Tags, ",")

	d1 := &dtov1.Device{
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
