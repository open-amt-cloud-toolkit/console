package redfish

// Define the structure for Redfish ComputerSystemCollection
type ComputerSystemCollection struct {
	OdataContext string `json:"@odata.context"`
	OdataID      string `json:"@odata.id"`
	OdataType    string `json:"@odata.type"`
	Name         string `json:"Name"`
	Members      []struct {
		OdataID string `json:"@odata.id"`
	} `json:"Members"`
	MembersOdataCount int `json:"Members@odata.count"`
}

// Define the structure for Redfish ComputerSystem
type ComputerSystem struct {
	OdataContext string `json:"@odata.context"`
	OdataID      string `json:"@odata.id"`
	OdataType    string `json:"@odata.type"`
	ID           string `json:"Id"`
	Name         string `json:"Name"`
	Manufacturer string `json:"Manufacturer"`
}
