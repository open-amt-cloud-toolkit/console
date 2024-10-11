package usecase

import (
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/security"

	"github.com/open-amt-cloud-toolkit/console/config"
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
	wifiConfigRepo := sqldb.NewWirelessRepo(database, log)
	key := config.ConsoleConfig.EncryptionKey
	safeRequirements := security.Crypto{
		EncryptionKey: key,
	}
	wsman1 := wsman.NewGoWSMANMessages(log, safeRequirements)
	wsman2 := amtexplorer.NewGoWSMANMessages(log, safeRequirements)
	domainRepo := sqldb.NewDomainRepo(database, log)
	deviceRepo := sqldb.NewDeviceRepo(database, log)
	ciraRepo := sqldb.NewCIRARepo(database, log)
	profileRepo := sqldb.NewProfileRepo(database, log)

	domains1 := domains.New(domainRepo, log, safeRequirements)
	wificonfig := wificonfigs.New(wifiConfigRepo, ieee, log, safeRequirements)

	return &Usecases{
		Domains:            domains1,
		Devices:            devices.New(deviceRepo, wsman1, devices.NewRedirector(safeRequirements), log, safeRequirements),
		AMTExplorer:        amtexplorer.New(deviceRepo, wsman2, log, safeRequirements),
		Profiles:           profiles.New(profileRepo, wifiConfigRepo, pwc, ieee, log, domainRepo, safeRequirements),
		IEEE8021xProfiles:  ieee,
		CIRAConfigs:        ciraconfigs.New(ciraRepo, log, safeRequirements),
		WirelessProfiles:   wificonfig,
		ProfileWiFiConfigs: pwc,
	}
}
