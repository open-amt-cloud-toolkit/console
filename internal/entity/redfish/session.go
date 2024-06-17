package redfish

// Define the structure for Redfish SessionService
type SessionService struct {
	OdataContext string `json:"@odata.context"`
	OdataID      string `json:"@odata.id"`
	OdataType    string `json:"@odata.type"`
	ID           string `json:"Id"`
	Name         string `json:"Name"`
}

// Define the structure for Redfish SessionCollection
type SessionCollection struct {
	OdataContext string `json:"@odata.context"`
	OdataID      string `json:"@odata.id"`
	OdataType    string `json:"@odata.type"`
	Name         string `json:"Name"`
	Members      []struct {
		OdataID string `json:"@odata.id"`
	} `json:"Members"`
	MembersOdataCount int `json:"Members@odata.count"`
}

type Session struct {
	OdataContext string `json:"@odata.context"`
	OdataID      string `json:"@odata.id"`
	OdataType    string `json:"@odata.type"`
	Name         string `json:"Name"`
	ID           string `json:"Id"`
}
