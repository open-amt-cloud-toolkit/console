package usecase

import (
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/ciraconfigs"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/domains"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/ieee8021xconfigs"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/postgresdb"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/profiles"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/wificonfigs"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
	"github.com/open-amt-cloud-toolkit/console/pkg/postgres"
)

// Usecases -.
type Usecases struct {
	Devices           devices.Feature
	Domains           domains.Feature
	Profiles          profiles.Feature
	IEEE8021xProfiles ieee8021xconfigs.Feature
	CIRAConfigs       ciraconfigs.Feature
	WirelessProfiles  wificonfigs.Feature
}

// New -.
func NewUseCases(pg *postgres.DB, log logger.Interface) *Usecases {
	return &Usecases{
		Domains:           domains.New(postgresdb.NewDomainRepo(pg, log), log),
		Devices:           devices.New(postgresdb.NewDeviceRepo(pg, log), wsman.NewGoWSMANMessages(), devices.NewRedirector(), log),
		Profiles:          profiles.New(postgresdb.NewProfileRepo(pg, log), log),
		IEEE8021xProfiles: ieee8021xconfigs.New(postgresdb.NewIEEE8021xRepo(pg, log), log),
		CIRAConfigs:       ciraconfigs.New(postgresdb.NewCIRARepo(pg, log), log),
		WirelessProfiles:  wificonfigs.New(postgresdb.NewWirelessRepo(pg, log), log),
	}
}
