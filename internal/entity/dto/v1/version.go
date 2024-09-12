package dtov1

type (
	SoftwareIdentity struct {
		InstanceID    string `json:"instanceID"`
		VersionString string `json:"VersionString" example:"<major>.<minor>.<revision>.<build>"`
		IsEntity      bool   `json:"isEntity" example:"true"`
	}

	Version struct {
		CIMSoftwareIdentity             []SoftwareIdentity                   `json:"cimSoftwareIdentity"`
		AMTSetupAndConfigurationService SetupAndConfigurationServiceResponse `json:"amtSetupAndConfigurationService"`
	}
)
