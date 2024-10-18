package dto

type (
	HardwareInfoResults struct {
		CIM_ComputerSystemPackage CIM_ComputerSystemPackage
		CIM_SystemPackage         CIM_SystemPackage
		CIM_Chassis               CIM_Chassis
		CIM_Chip                  CIM_Chip
		CIM_Card                  CIM_Card
		CIM_BIOSElement           CIM_BIOSElement
		CIM_Processor             CIM_Processor
		CIM_PhysicalPackage       CIM_PhysicalPackage
		CIM_PhysicalMemory        CIM_PhysicalMemory
		CIM_MediaAccessDevice     CIM_MediaAccessDevice
	}

	CIM_BIOSElement struct {
		Response  CIMBIOSElementResponse `json:"response"`
		Responses []any                  `json:"responses"`
	}

	CIM_Card struct {
		Response  CIMCardResponseGet `json:"response"`
		Responses []any              `json:"responses"`
	}

	CIM_Chassis struct {
		Response  CIMChassisResponse `json:"response"`
		Responses []any              `json:"responses"`
	}

	CIM_Chip struct {
		Responses []CIMChipGet `json:"responses"`
	}

	CIM_ComputerSystemPackage struct {
		// PlatformGUID string
		Response  string `json:"response"`
		Responses string `json:"responses"`
	}

	CIM_MediaAccessDevice struct {
		Responses []CIMMediaAccessDevice `json:"responses"`
	}

	CIM_PhysicalMemory struct {
		Responses []CIMPhysicalMemoryResponse `json:"responses"`
	}

	CIM_PhysicalPackage struct {
		Responses []CIMPhysicalPackageResponses
	}

	CIM_Processor struct {
		Responses []CIMProcessorResponse `json:"responses"`
	}

	CIM_SystemPackage struct {
		responses []CIMSystemPackagingResponses
	}

	CIMSystemPackagingResponses struct {
		SystemPackageItems any
	}

	CIMChassisResponse struct {
		Version            string
		SerialNumber       string
		Model              string
		Manufacturer       string
		ElementName        string
		CreationClassName  string
		Tag                string
		OperationalStatus  []int
		PackageType        int
		ChassisPackageType int
	}

	CIMBIOSElementResponse struct {
		TargetOperatingSystem TargetOperatingSystem
		SoftwareElementID     string
		SoftwareElementState  SoftwareElementState
		Name                  string
		OperationalStatus     []int
		ElementName           string
		Version               string
		Manufacturer          string
		PrimaryBIOS           bool
		ReleaseDate           Time
	}

	CIMChip struct {
		Pull []any
		Get  CIMChipGet
	}

	CIMChipGet struct {
		CanBeFRUed        bool
		CreationClassName string
		ElementName       string
		Manufacturer      string
		OperationalStatus []int
		Tag               string
		Version           string
	}

	CIMProcessorResponse struct {
		DeviceID                string
		CreationClassName       string
		SystemName              string
		SystemCreationClassName string
		ElementName             string
		OperationalStatus       []int
		HealthState             int
		EnabledState            int
		RequestedState          int
		Role                    string
		Family                  int
		OtherFamilyDescription  string
		UpgradeMethod           int
		MaxClockSpeed           int
		CurrentClockSpeed       int
		Stepping                string
		CPUStatus               int
		ExternalBusClockSpeed   int
	}

	CIMMediaAccessDevice struct {
		Pull []any
		Get  struct {
			Capabilities            []int
			CreationClassName       string
			DeviceID                string
			ElementName             string
			EnabledDefault          int
			EnabledState            int
			MaxMediaSize            int
			OperationalStatus       []int
			RequestedState          int
			Security                int
			SystemCreationClassName string
			SystemName              string
		}
	}

	CIMPhysicalMemoryResponse struct {
		PartNumber                 string
		SerialNumber               string
		Manufacturer               string
		ElementName                string
		CreationClassName          string
		Tag                        string
		OperationalStatus          []int
		FormFactor                 int
		MemoryType                 int
		Speed                      int
		Capacity                   int
		BankLabel                  string
		ConfiguredMemoryClockSpeed int
		IsSpeedInMhz               bool
		MaxMemorySpeed             int
	}

	CIMPhysicalPackageResponses struct {
		PullMemoryItems []any
		PullCardItems   []any
	}

	CIMCardResponse struct {
		Pull []any
		Get  CIMCardResponseGet
	}

	CIMCardResponseGet struct {
		CanBeFRUed        bool
		CreationClassName string
		ElementName       string
		Manufacturer      string
		Model             string
		OperationalStatus []int
		PackageType       int
		SerialNumber      string
		Tag               string
		Version           string
	}

	// CIMComputerSystemPackage struct {
	// 	PlatformGUID string
	// }

	// CIMSystemPackage struct {
	// 	responses []SystemPackageItems
	// }

	// CIMChassis struct {
	// 	Version            string
	// 	SerialNumber       string
	// 	Model              string
	// 	Manufacturer       string
	// 	ElementName        string
	// 	CreationClassName  string
	// 	Tag                string
	// 	OperationalStatus  []int
	// 	PackageType        PackageType
	// 	ChassisPackageType ChassisPackageType
	// }

	// CIMChip struct {
	// 	responses []interface{}
	// }

	// CIMCard struct {
	// 	response  hwInfo.Card
	// 	responses []interface{}
	// }
	// CIMBIOSElement struct {
	// 	response  hwInfo.BIOSElement
	// 	responses []interface{}
	// }

	// CIMProcessor struct {
	// 	responses []interface{}
	// }

	// CIMPhysicalMemory struct {
	// 	responses hwInfo.PhysicalMemory
	// }

	// CIMMediaAccessDevice struct {
	// 	responses []interface{}
	// }

	// CIMPhysicalPackage struct {
	// 	responses []interface{}
	// }

	// SystemPackageItems struct {
	// 	PlatformGUID string
	// }

	Time struct {
		DateTime string `xml:"Datetime"`
	}

	// ChassisPackageType is an enumeration defining the type of the PhysicalPackage.
	ChassisPackageType int

	// PackageType is the type of the PhysicalPackage.
	PackageType int

	// TargetOperatingSystem is the element's operating system environment.
	TargetOperatingSystem int

	// SoftwareElementState is defined in this model to identify various states of a SoftwareElement's life cycle.
	SoftwareElementState int
)
