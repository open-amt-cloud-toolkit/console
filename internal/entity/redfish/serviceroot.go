package redfish

// Define the structure for Redfish ServiceRoot
type ServiceRoot struct {
	OdataContext   string `json:"@odata.context"`
	OdataID        string `json:"@odata.id"`
	OdataType      string `json:"@odata.type"`
	ID             string `json:"Id"`
	Name           string `json:"Name"`
	RedfishVersion string `json:"RedfishVersion"`
	UUID           string `json:"UUID"`
	SessionService struct {
		OdataID string `json:"@odata.id"`
	} `json:"SessionService"`
	Systems struct {
		OdataID string `json:"@odata.id"`
	} `json:"Systems"`
	Links struct {
		Sessions struct {
			OdataID string `json:"@odata.id"`
		} `json:"Sessions"`
	} `json:"Links"`
}
