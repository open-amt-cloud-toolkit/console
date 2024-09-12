package dtov2

type (
	Version struct {
		Flash               string `json:"FLash" example:"<major>.<minor>.<revision>.<build>"`
		Netstack            string `json:"Netstack" example:"<major>.<minor>.<revision>.<build>"`
		AMTApps             string `json:"AmtApps" example:"<major>.<minor>.<revision>.<build>"`
		AMT                 string `json:"Amt" example:"<major>.<minor>.<revision>.<build>"`
		Sku                 string `json:"Sku" example:"<major>.<minor>.<revision>.<build>"`
		VendorID            string `json:"VendorID" example:"<major>.<minor>.<revision>.<build>"`
		BuildNumber         string `json:"BuildNumber" example:"<major>.<minor>.<revision>.<build>"`
		RecoveryVersion     string `json:"Recovery" example:"<major>.<minor>.<revision>.<build>"`
		RecoveryBuildNumber string `json:"RecoveryBuildNumber" example:"<major>.<minor>.<revision>.<build>"`
		LegacyMode          *bool  `json:"LegacyMode" example:"false"`
		AmtFWCoreVersion    string `json:"AmtFWCore" example:"<major>.<minor>.<revision>.<build>"`
	}
)
