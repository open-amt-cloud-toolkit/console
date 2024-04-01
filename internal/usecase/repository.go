package usecase

import (
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/postgresdb"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/wsman"
	"github.com/open-amt-cloud-toolkit/console/pkg/postgres"
)

// Repositories -.
type Repositories struct {
	Domains          Domain
	Devices          Device
	DeviceManagement DeviceManagement
	Profiles         Profile
}

// New -.
func New(pg *postgres.DB) *Repositories {
	return &Repositories{
		Devices:          postgresdb.NewDeviceRepo(pg),
		Domains:          postgresdb.NewDomainRepo(pg),
		DeviceManagement: wsman.NewGoWSMANMessages(),
		Profiles:         postgresdb.NewProfileRepo(pg),
	}
}
