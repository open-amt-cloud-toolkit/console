package usecase

import (
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/amtexplorer"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/ciraconfigs"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/domains"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/ieee8021xconfigs"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/profiles"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/profilewificonfigs"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/sqldb"
	"github.com/open-amt-cloud-toolkit/console/internal/usecase/wificonfigs"
	"github.com/open-amt-cloud-toolkit/console/pkg/db"
	"github.com/open-amt-cloud-toolkit/console/pkg/logger"
)

// Usecases -.
type Usecases struct {
	Devices            devices.Feature
	Domains            domains.Feature
	AMTExplorer        amtexplorer.Feature
	Profiles           profiles.Feature
	ProfileWiFiConfigs profilewificonfigs.Feature
	IEEE8021xProfiles  ieee8021xconfigs.Feature
	CIRAConfigs        ciraconfigs.Feature
	WirelessProfiles   wificonfigs.Feature
}

// New -.
func NewUseCases(database *db.SQL, log logger.Interface) *Usecases {
	pwc := profilewificonfigs.New(sqldb.NewProfileWiFiConfigsRepo(database, log), log)
	ieee := ieee8021xconfigs.New(sqldb.NewIEEE8021xRepo(database, log), log)
	wificonfig := wificonfigs.New(sqldb.NewWirelessRepo(database, log), ieee, log)
	wsman1 := wsman.NewGoWSMANMessages(log)
	wsman2 := amtexplorer.NewGoWSMANMessages(log)

	return &Usecases{
		Domains:            domains.New(sqldb.NewDomainRepo(database, log), log),
		Devices:            devices.New(sqldb.NewDeviceRepo(database, log), wsman1, devices.NewRedirector(), log),
		AMTExplorer:        amtexplorer.New(sqldb.NewDeviceRepo(database, log), wsman2, log),
		Profiles:           profiles.New(sqldb.NewProfileRepo(database, log), wificonfig, pwc, ieee, log),
		IEEE8021xProfiles:  ieee,
		CIRAConfigs:        ciraconfigs.New(sqldb.NewCIRARepo(database, log), log),
		WirelessProfiles:   wificonfig,
		ProfileWiFiConfigs: pwc,
	}
}
