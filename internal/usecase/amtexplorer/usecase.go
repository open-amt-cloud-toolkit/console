package amtexplorer

import (
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/security"

	"github.com/open-amt-cloud-toolkit/console/internal/usecase/sqldb"
	"github.com/open-amt-cloud-toolkit/console/pkg/consoleerrors"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

var ErrDatabase = sqldb.DatabaseError{Console: consoleerrors.CreateConsoleError("DevicesUseCase")}

// UseCase -.
type UseCase struct {
	repo             Repository
	device           WSMAN
	log              logger.Interface
	safeRequirements security.Cryptor
}

var ErrAMT = AMTError{Console: consoleerrors.CreateConsoleError("DevicesUseCase")}

// New -.
func New(r Repository, d WSMAN, log logger.Interface, safeRequirements security.Cryptor) *UseCase {
	return &UseCase{
		repo:             r,
		device:           d,
		log:              log,
		safeRequirements: safeRequirements,
	}
}
