package usecase

import (
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/ciraconfigs"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/domains"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/ieee8021xconfigs"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/profiles"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/sqldb"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/wificonfigs"
	"github.com/open-amt-cloud-toolkit/console/pkg/db"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
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
func NewUseCases(pg *db.SQL, log logger.Interface) *Usecases {
	return &Usecases{
		Domains:           domains.New(sqldb.NewDomainRepo(pg, log), log),
		Devices:           devices.New(sqldb.NewDeviceRepo(pg, log), wsman.NewGoWSMANMessages(), devices.NewRedirector(), log),
		Profiles:          profiles.New(sqldb.NewProfileRepo(pg, log), log),
		IEEE8021xProfiles: ieee8021xconfigs.New(sqldb.NewIEEE8021xRepo(pg, log), log),
		CIRAConfigs:       ciraconfigs.New(sqldb.NewCIRARepo(pg, log), log),
		WirelessProfiles:  wificonfigs.New(sqldb.NewWirelessRepo(pg, log), log),
	}
}
