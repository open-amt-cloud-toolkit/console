package devices

import "github.com/open-amt-cloud-toolkit/console/internal/features/amt"

type Device struct {
	Id                int
	Name              string
	Address           string
	Username          string
	Password          string
	UseTLS            bool
	SelfSignedAllowed bool
	PowerState        string
	AMTSpecific       amt.AMTSpecific
	BMCSpecific       BMCSpecific
	DASHSpecific      DASHSpecific
	RedfishSpecific   RedfishSpecific
}

type BMCSpecific struct {
}

type DASHSpecific struct {
}

type RedfishSpecific struct {
}
