package devices

type PowerState string

const (
	PowerOn                   PowerState = "Power On"
	SleepLight                PowerState = "Sleep Light (OS)"
	SleepDeep                 PowerState = "Sleep Deep (OS)"
	PowerCycleOffSoft         PowerState = "Soft Power Cycle (OS Graceful)"
	PowerOffHard              PowerState = "Hard Power Off"
	Hibernate                 PowerState = "Hibernate (OS)"
	PowerOffSoft              PowerState = "Soft Power Off (OS Graceful)"
	PowerCycleOffHard         PowerState = "Hard Power Cycle"
	MasterBusReset            PowerState = "Master Bus Reset"
	DiagnosticInterruptNMI    PowerState = "Diagnostic Interrupt NMI"
	PowerOffSoftGraceful      PowerState = "Soft Power Off (OS Graceful)"
	PowerCycleOffHardGraceful PowerState = "Hard Power Cycle (OS Graceful)"
)
