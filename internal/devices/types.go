package devices

import "github.com/jritsema/go-htmx-starter/internal/features/amt"

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
