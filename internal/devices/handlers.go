package devices

func ProvisioningModeLookup(mode int) string {
	valueMap := map[int]string{
		1: "Admin Control Mode",
		4: "Client Control Mode",
	}

	result, ok := valueMap[mode]
	if !ok {
		result = "invalid provisioning mode"
	}

	return result
}

func ProvisioningStateLookup(state int) string {
	valueMap := map[int]string{
		0: "Pre-Provisioning",
		1: "In Provisioning",
		2: "Post Provisioning",
	}

	result, ok := valueMap[state]
	if !ok {
		result = "invalid provisoining state"
	}

	return result
}

func PowerControlLookup(value int) string {
	valueMap := map[int]string{
		2:  "PowerOn",
		3:  "SleepLight",
		4:  "SleepDeep",
		5:  "PowerCycleOffSoft",
		6:  "PowerOffHard",
		7:  "Hibernate",
		8:  "PowerOffSoft",
		9:  "PowerCycleOffHard",
		10: "MasterBusReset",
		11: "DiagnosticInterruptNMI",
		12: "PowerOffSoftGraceful",
		13: "PowerOffHardGraceful",
		14: "MasterBusResetGraceful",
		15: "PowerCycleOffSoftGraceful",
		16: "PowerCycleOffHardGraceful",
	}

	result, ok := valueMap[value]
	if !ok {
		result = "invalid power control value"
	}

	return result
}
