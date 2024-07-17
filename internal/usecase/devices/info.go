package devices

import (
	"context"
	"unsafe"

	"github.com/open-amt-cloud-toolkit/console/internal/entity/dto"
	wsmanAPI "github.com/open-amt-cloud-toolkit/console/internal/usecase/devices/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/bios"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/card"
	// "github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/chassis"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/chip"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/mediaaccess"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/physical"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/wsman/cim/processor"
)

func (uc *UseCase) GetVersion(c context.Context, guid string) (map[string]interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	version, err := device.GetAMTVersion()
	if err != nil {
		return nil, err
	}

	data, err := device.GetSetupAndConfiguration()
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"CIM_SoftwareIdentity": map[string]interface{}{
			"responses": version,
		},
		"AMT_SetupAndConfigurationService": map[string]interface{}{
			"response": data[0],
		},
	}

	return response, nil
}

func (uc *UseCase) GetFeatures(c context.Context, guid string) (dto.Features, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return dto.Features{}, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	features, err := device.GetFeatures()
	if err != nil {
		return dto.Features{}, err
	}

	return features, nil
}

func (uc *UseCase) SetFeatures(c context.Context, guid string, features dto.Features) (dto.Features, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return features, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	features, err = device.SetFeatures(features)
	if err != nil {
		return features, err
	}

	return features, nil
}

func (uc *UseCase) GetHardwareInfo(c context.Context, guid string) (dto.HardwareInfoResults, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return dto.HardwareInfoResults{}, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	hwInfo, err := device.GetHardwareInfo()
	if err != nil {
		return dto.HardwareInfoResults{}, err
	}

	d1 := *uc.getHardwareInfoEntityToDTO(&hwInfo)

	return d1, nil
}

func (uc *UseCase) GetAuditLog(c context.Context, startIndex int, guid string) (dto.AuditLog, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return dto.AuditLog{}, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	response, err := device.GetAuditLog(startIndex)
	if err != nil {
		return dto.AuditLog{}, err
	}

	auditLogResponse := dto.AuditLog{}
	auditLogResponse.TotalCount = response.Body.ReadRecordsResponse.TotalRecordCount
	auditLogResponse.Records = response.Body.DecodedRecordsResponse

	return auditLogResponse, nil
}

func (uc *UseCase) GetEventLog(c context.Context, guid string) ([]dto.EventLog, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	eventLogs, err := device.GetEventLog()
	if err != nil {
		return nil, err
	}

	events := make([]dto.EventLog, len(eventLogs.RefinedEventData))

	for idx := range eventLogs.RefinedEventData {
		event := &eventLogs.RefinedEventData[idx]
		dtoEvent := dto.EventLog{
			// DeviceAddress:   event.DeviceAddress,
			// EventSensorType: event.EventSensorType,
			// EventType:       event.EventType,
			// EventOffset:     event.EventOffset,
			// EventSourceType: event.EventSourceType,
			EventSeverity: event.EventSeverity,
			// SensorNumber:    event.SensorNumber,
			Entity: event.Entity,
			// EntityInstance:  event.EntityInstance,
			// EventData:       event.EventData,
			Time: event.TimeStamp.String(),
			// EntityStr:       event.EntityStr,
			Description: event.Description,
			// EventTypeDesc:   event.EventTypeDesc,
		}

		events[idx] = dtoEvent
	}

	return events, nil
}

func (uc *UseCase) GetGeneralSettings(c context.Context, guid string) (interface{}, error) {
	item, err := uc.GetByID(c, guid, "")
	if err != nil {
		return nil, err
	}

	device := uc.device.SetupWsmanClient(*item, false, true)

	generalSettings, err := device.GetGeneralSettings()
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"Body": generalSettings,
	}

	return response, nil
}

func (uc *UseCase) getHardwareInfoEntityToDTO(d *wsmanAPI.HWResults) *dto.HardwareInfoResults {
	d1 := &dto.HardwareInfoResults{
		ComputerSystemPackage: dto.CIMComputerSystemPackage{
			PlatformGUID: d.CSPResult.Body.GetResponse.PlatformGUID,
		},
		SystemPackage: dto.CIMSystemPackage{},
		Chassis: dto.CIMChassis{
			Version:           d.ChassisResult.Body.PackageResponse.Version,
			SerialNumber:      d.ChassisResult.Body.PackageResponse.SerialNumber,
			Model:             d.ChassisResult.Body.PackageResponse.Model,
			Manufacturer:      d.ChassisResult.Body.PackageResponse.Manufacturer,
			ElementName:       d.ChassisResult.Body.PackageResponse.ElementName,
			CreationClassName: d.ChassisResult.Body.PackageResponse.CreationClassName,
			Tag:               d.ChassisResult.Body.PackageResponse.Tag,
			OperationalStatus: *(*[]int)(unsafe.Pointer(&d.ChassisResult.Body.PackageResponse.OperationalStatus)),
			// OperationalStatus:  operationalStatusToDTO(d.ChassisResult.Body.PackageResponse.OperationalStatus),
			PackageType:        dto.PackageType(d.ChassisResult.Body.PackageResponse.PackageType),
			ChassisPackageType: dto.ChassisPackageType(d.ChassisResult.Body.PackageResponse.ChassisPackageType),
		},
		Chip: dto.CIMChip{
			Get: dto.ChipItems{
				CanBeFRUed:        d.ChipResult.Body.PackageResponse.CanBeFRUed,
				CreationClassName: d.ChipResult.Body.PackageResponse.CreationClassName,
				ElementName:       d.ChipResult.Body.PackageResponse.ElementName,
				Manufacturer:      d.ChipResult.Body.PackageResponse.Manufacturer,
				OperationalStatus: *(*[]int)(unsafe.Pointer(&d.ChipResult.Body.PackageResponse.OperationalStatus)),
				// OperationalStatus: d.ChipResult.Body.PackageResponse.OperationalStatus,
				Tag:     d.ChipResult.Body.PackageResponse.Tag,
				Version: d.ChipResult.Body.PackageResponse.Version,
			},
			Pull: chipItemsToDTO(d.ChipResult.Body.PullResponse.ChipItems),
		},
		Card: dto.CIMCard{
			Get: dto.CardItems{
				CanBeFRUed:        d.CardResult.Body.PackageResponse.CanBeFRUed,
				CreationClassName: d.CardResult.Body.PackageResponse.CreationClassName,
				ElementName:       d.CardResult.Body.PackageResponse.ElementName,
				Manufacturer:      d.CardResult.Body.PackageResponse.Manufacturer,
				Model:             d.CardResult.Body.PackageResponse.Model,
				OperationalStatus: *(*[]int)(unsafe.Pointer(&d.CardResult.Body.PackageResponse.OperationalStatus)),
				// OperationalStatus: d.CardResult.Body.PackageResponse.OperationalStatus,
				PackageType:  dto.PackageType(d.CardResult.Body.PackageResponse.PackageType),
				SerialNumber: d.CardResult.Body.PackageResponse.SerialNumber,
				Tag:          d.CardResult.Body.PackageResponse.Tag,
				Version:      d.CardResult.Body.PackageResponse.Version,
			},
			Pull: cardItemsToDTO(d.CardResult.Body.PullResponse.CardItems),
		},
		BIOSElement: dto.CIMBIOSElement{
			// Get: dto.BiosElement{
				TargetOperatingSystem: dto.TargetOperatingSystem(d.BiosResult.Body.GetResponse.TargetOperatingSystem),
				SoftwareElementID:     d.BiosResult.Body.GetResponse.SoftwareElementID,
				SoftwareElementState:  dto.SoftwareElementState(d.BiosResult.Body.GetResponse.SoftwareElementState),
				Name:              d.BiosResult.Body.GetResponse.Name,
				OperationalStatus: *(*[]int)(unsafe.Pointer(&d.BiosResult.Body.GetResponse.OperationalStatus)),
				// OperationalStatus:     d.BiosResult.Body.GetResponse.OperationalStatus,
				ElementName:  d.BiosResult.Body.GetResponse.ElementName,
				Version:      d.BiosResult.Body.GetResponse.Version,
				Manufacturer: d.BiosResult.Body.GetResponse.Manufacturer,
				PrimaryBIOS:  d.BiosResult.Body.GetResponse.PrimaryBIOS,
				ReleaseDate:  dto.Time(d.BiosResult.Body.GetResponse.ReleaseDate),
			// },
			// Pull: biosItemsToDTO(d.BiosResult.Body.PullResponse.BiosElementItems),
		},
		Processor: dto.CIMProcessor{
			Get: dto.ProcessorItems{
				DeviceID:                d.ProcessorResult.Body.PackageResponse.DeviceID,
				CreationClassName:       d.ProcessorResult.Body.PackageResponse.CreationClassName,
				SystemName:              d.ProcessorResult.Body.PackageResponse.SystemName,
				SystemCreationClassName: d.ProcessorResult.Body.PackageResponse.SystemCreationClassName,
				ElementName:             d.ProcessorResult.Body.PackageResponse.ElementName,
				OperationalStatus:       *(*[]int)(unsafe.Pointer(&d.ProcessorResult.Body.PackageResponse.OperationalStatus)),
				// OperationalStatus:       d.ProcessorResult.Body.PackageResponse.OperationalStatus,
				HealthState:            dto.HealthState(d.ProcessorResult.Body.PackageResponse.HealthState),
				EnabledState:           dto.EnabledState(d.ProcessorResult.Body.PackageResponse.EnabledState),
				RequestedState:         dto.RequestedState(d.ProcessorResult.Body.PackageResponse.RequestedState),
				Role:                   d.ProcessorResult.Body.PackageResponse.Role,
				Family:                 d.ProcessorResult.Body.PackageResponse.Family,
				OtherFamilyDescription: d.ProcessorResult.Body.PackageResponse.OtherFamilyDescription,
				UpgradeMethod:          dto.UpgradeMethod(d.ProcessorResult.Body.PackageResponse.UpgradeMethod),
				MaxClockSpeed:          d.ProcessorResult.Body.PackageResponse.MaxClockSpeed,
				CurrentClockSpeed:      d.ProcessorResult.Body.PackageResponse.CurrentClockSpeed,
				Stepping:               d.ProcessorResult.Body.PackageResponse.Stepping,
				CPUStatus:              dto.CPUStatus(d.ProcessorResult.Body.PackageResponse.CPUStatus),
				ExternalBusClockSpeed:  d.ProcessorResult.Body.PackageResponse.ExternalBusClockSpeed,
			},
			Pull: processorItemsToDTO(d.ProcessorResult.Body.PullResponse.PackageItems),
		},
		PhysicalMemory: dto.CIMPhysicalMemory{
			Get: dto.PhysicalMemory{
				PartNumber:        d.PhysicalMemoryResult.Body.MemoryResponse.PartNumber,
				SerialNumber:      d.PhysicalMemoryResult.Body.MemoryResponse.SerialNumber,
				Manufacturer:      d.PhysicalMemoryResult.Body.MemoryResponse.Manufacturer,
				ElementName:       d.PhysicalMemoryResult.Body.MemoryResponse.ElementName,
				CreationClassName: d.PhysicalMemoryResult.Body.MemoryResponse.CreationClassName,
				Tag:               d.PhysicalMemoryResult.Body.MemoryResponse.Tag,
				OperationalStatus: *(*[]int)(unsafe.Pointer(&d.PhysicalMemoryResult.Body.MemoryResponse.OperationalStatus)),
				// OperationalStatus:          d.PhysicalMemoryResult.Body.MemoryResponse.OperationalStatus,
				FormFactor:                 d.PhysicalMemoryResult.Body.MemoryResponse.FormFactor,
				MemoryType:                 dto.MemoryType(d.PhysicalMemoryResult.Body.MemoryResponse.MemoryType),
				Speed:                      d.PhysicalMemoryResult.Body.MemoryResponse.Speed,
				Capacity:                   d.PhysicalMemoryResult.Body.MemoryResponse.Capacity,
				BankLabel:                  d.PhysicalMemoryResult.Body.MemoryResponse.BankLabel,
				ConfiguredMemoryClockSpeed: d.PhysicalMemoryResult.Body.MemoryResponse.ConfiguredMemoryClockSpeed,
				IsSpeedInMhz:               d.PhysicalMemoryResult.Body.MemoryResponse.IsSpeedInMhz,
				MaxMemorySpeed:             d.PhysicalMemoryResult.Body.MemoryResponse.MaxMemorySpeed,
			},
			Pull: physicalMemoryToDTO(d.PhysicalMemoryResult.Body.PullResponse.MemoryItems),
		},
		MediaAccessDevices: dto.CIMMediaAccessDevice{
			Pull: mediaAccessDeviceToDTO(d.MediaAccessPullResult.Body.PullResponse.MediaAccessDevices),
		},
		PhysicalPackage: dto.CIMPhysicalPackage{
			PullMemoryItems: ppPullResponseMemoryToDTO(d.PPPullResult.Body.PullResponse),
			PullCardItems:   ppPullResponseCardToDTO(d.PPPullResult.Body.PullResponse),
		},
	}

	return d1
}

func chipItemsToDTO(d []chip.PackageResponse) []dto.ChipItems {
	// iterate over the data and convert each entity to dto
	d2 := make([]dto.ChipItems, len(d))

	for i := range d {
		d2[i] = dto.ChipItems{
			CanBeFRUed:        d[i].CanBeFRUed,
			CreationClassName: d[i].CreationClassName,
			ElementName:       d[i].ElementName,
			Manufacturer:      d[i].Manufacturer,
			OperationalStatus: *(*[]int)(unsafe.Pointer(&d[i].OperationalStatus)),
			Tag:               d[i].Tag,
			Version:           d[i].Version,
		}
	}

	return d2
}

func cardItemsToDTO(d []card.PackageResponse) []dto.CardItems {
	// iterate over the data and convert each entity to dto
	d2 := make([]dto.CardItems, len(d))

	for i := range d {
		d2[i] = dto.CardItems{
			CanBeFRUed:        d[i].CanBeFRUed,
			CreationClassName: d[i].CreationClassName,
			ElementName:       d[i].ElementName,
			Manufacturer:      d[i].Manufacturer,
			Model:             d[i].Model,
			OperationalStatus: *(*[]int)(unsafe.Pointer(&d[i].OperationalStatus)),
			PackageType:       dto.PackageType(d[i].PackageType),
			SerialNumber:      d[i].SerialNumber,
			Tag:               d[i].Tag,
			Version:           d[i].Version,
		}
	}

	return d2
}

func biosItemsToDTO(d []bios.BiosElement) []dto.CIMBIOSElement {
	// iterate over the data and convert each entity to dto
	d2 := make([]dto.CIMBIOSElement, len(d))

	for i := range d {
		d2[i] = dto.CIMBIOSElement{
			TargetOperatingSystem: dto.TargetOperatingSystem(d[i].TargetOperatingSystem),
			SoftwareElementID:     d[i].SoftwareElementID,
			SoftwareElementState:  dto.SoftwareElementState(d[i].SoftwareElementState),
			Name:                  d[i].Name,
			OperationalStatus:     *(*[]int)(unsafe.Pointer(&d[i].OperationalStatus)),
			ElementName:           d[i].ElementName,
			Version:               d[i].Version,
			Manufacturer:          d[i].Manufacturer,
			PrimaryBIOS:           d[i].PrimaryBIOS,
			ReleaseDate:           dto.Time(d[i].ReleaseDate),
		}
	}

	return d2
}

func processorItemsToDTO(d []processor.PackageResponse) []dto.ProcessorItems {
	// iterate over the data and convert each entity to dto
	d2 := make([]dto.ProcessorItems, len(d))

	for i := range d {
		d2[i] = dto.ProcessorItems{
			DeviceID:                d[i].DeviceID,
			CreationClassName:       d[i].CreationClassName,
			SystemName:              d[i].SystemName,
			SystemCreationClassName: d[i].SystemCreationClassName,
			ElementName:             d[i].ElementName,
			OperationalStatus:       *(*[]int)(unsafe.Pointer(&d[i].OperationalStatus)),
			HealthState:             dto.HealthState(d[i].HealthState),
			EnabledState:            dto.EnabledState(d[i].EnabledState),
			RequestedState:          dto.RequestedState(d[i].RequestedState),
			Role:                    d[i].Role,
			Family:                  d[i].Family,
			OtherFamilyDescription:  d[i].OtherFamilyDescription,
			UpgradeMethod:           dto.UpgradeMethod(d[i].UpgradeMethod),
			MaxClockSpeed:           d[i].MaxClockSpeed,
			CurrentClockSpeed:       d[i].CurrentClockSpeed,
			Stepping:                d[i].Stepping,
			CPUStatus:               dto.CPUStatus(d[i].CPUStatus),
			ExternalBusClockSpeed:   d[i].ExternalBusClockSpeed,
		}
	}

	return d2
}

func physicalMemoryToDTO(d []physical.PhysicalMemory) []dto.PhysicalMemory {
	// iterate over the data and convert each entity to dto
	d2 := make([]dto.PhysicalMemory, len(d))

	for i := range d {
		d2[i] = dto.PhysicalMemory{
			PartNumber:                 d[i].PartNumber,
			SerialNumber:               d[i].SerialNumber,
			Manufacturer:               d[i].Manufacturer,
			ElementName:                d[i].ElementName,
			CreationClassName:          d[i].CreationClassName,
			Tag:                        d[i].Tag,
			OperationalStatus:          *(*[]int)(unsafe.Pointer(&d[i].OperationalStatus)),
			FormFactor:                 d[i].FormFactor,
			MemoryType:                 dto.MemoryType(d[i].MemoryType),
			Speed:                      d[i].Speed,
			Capacity:                   d[i].Capacity,
			BankLabel:                  d[i].BankLabel,
			ConfiguredMemoryClockSpeed: d[i].ConfiguredMemoryClockSpeed,
			IsSpeedInMhz:               d[i].IsSpeedInMhz,
			MaxMemorySpeed:             d[i].MaxMemorySpeed,
		}
	}

	return d2
}

func mediaAccessDeviceToDTO(d []mediaaccess.MediaAccessDevice) []dto.MediaAccessDevice {
	// iterate over the data and convert each entity to dto
	d2 := make([]dto.MediaAccessDevice, len(d))

	for i := range d {
		d2[i] = dto.MediaAccessDevice{
			// Capabilities:            d[1].Capabilities,
			CreationClassName:       d[1].CreationClassName,
			DeviceID:                d[1].DeviceID,
			ElementName:             d[1].ElementName,
			EnabledDefault:          dto.EnabledDefault(d[1].EnabledDefault),
			EnabledState:            dto.EnabledState(d[1].EnabledState),
			MaxMediaSize:            d[1].MaxMediaSize,
			OperationalStatus:       *(*[]int)(unsafe.Pointer(&d[i].OperationalStatus)),
			RequestedState:          dto.RequestedState(d[1].RequestedState),
			Security:                dto.Security(d[1].Security),
			SystemCreationClassName: d[1].SystemCreationClassName,
			SystemName:              d[1].SystemName}
	}

	return d2
}

func ppPullResponseMemoryToDTO(d physical.PullResponse) []dto.PhysicalMemory {
	// iterate over the data and convert each entity to dto
	d2 := make([]dto.PhysicalMemory, len(d.MemoryItems))

	for i := range d.MemoryItems {
		d2[i] = dto.PhysicalMemory{
			PartNumber:                 d.MemoryItems[i].PartNumber,
			SerialNumber:               d.MemoryItems[i].SerialNumber,
			Manufacturer:               d.MemoryItems[i].Manufacturer,
			ElementName:                d.MemoryItems[i].ElementName,
			CreationClassName:          d.MemoryItems[i].CreationClassName,
			Tag:                        d.MemoryItems[i].Tag,
			OperationalStatus:          *(*[]int)(unsafe.Pointer(&d.MemoryItems[i].OperationalStatus)),
			FormFactor:                 d.MemoryItems[i].FormFactor,
			MemoryType:                 dto.MemoryType(d.MemoryItems[i].MemoryType),
			Speed:                      d.MemoryItems[i].Speed,
			Capacity:                   d.MemoryItems[i].Capacity,
			BankLabel:                  d.MemoryItems[i].BankLabel,
			ConfiguredMemoryClockSpeed: d.MemoryItems[i].ConfiguredMemoryClockSpeed,
			IsSpeedInMhz:               d.MemoryItems[i].IsSpeedInMhz,
			MaxMemorySpeed:             d.MemoryItems[i].MaxMemorySpeed,
		}
	}

	return d2
}

func ppPullResponseCardToDTO(d physical.PullResponse) []dto.CardItems {
	// iterate over the data and convert each entity to dto
	d2 := make([]dto.CardItems, len(d.PhysicalPackage))

	for i := range d.PhysicalPackage {
		d2[i] = dto.CardItems{
			CanBeFRUed:        d.PhysicalPackage[i].CanBeFRUed,
			CreationClassName: d.PhysicalPackage[i].CreationClassName,
			ElementName:       d.PhysicalPackage[i].ElementName,
			Manufacturer:      d.PhysicalPackage[i].Manufacturer,
			Model:             d.PhysicalPackage[i].Model,
			OperationalStatus: *(*[]int)(unsafe.Pointer(&d.PhysicalPackage[i].OperationalStatus)),
			PackageType:       dto.PackageType(d.PhysicalPackage[i].PackageType),
			SerialNumber:      d.PhysicalPackage[i].SerialNumber,
			Tag:               d.PhysicalPackage[i].Tag,
			Version:           d.PhysicalPackage[i].Version,
		}
	}

	return d2
}

// func operationalStatusToDTO[T int](d []T) []int {
// 	// iterate over the data and convert each entity to dto
// 	d2 := make([]int, len(d))

// 	for i := range d {
// 		d2[i] = int(d[i])
// 	}
// 	return d2
// }

// func operationalStatusToDTO[T IntBase](d []T) []int {
// 	d2 := make([]int, len(d))
// 	for i, v := range d {
// 		d2[i] = int(v) // Explicit conversion to int
// 	}
// 	return d2
// }

// Define an interface for types that can convert to int
// type IntConvertible interface {
// 	ToInt() int
// }

// // Implement the interface for custom types
// func (v chassis.OperationalStatus) ToInt() int {
// 	return int(v)
// }

// // Generic function to convert slices of any custom int-based type to slices of int
// func operationalStatus[T IntConvertible](d []T) []int {
// 	d2 := make([]int, len(d))
// 	for i, v := range d {
// 		d2[i] = v.ToInt() // Use the interface method to convert to int
// 	}
// 	return d2
// }
