package dto

type (
	SecuritySettings struct {
		ProfileAssociation []ProfileAssociation `json:"profileAssociation"`
		Certificates       interface{}          `json:"certificates"`
		Keys               interface{}          `json:"publicKeys"`
	}

	ProfileAssociation struct {
		Type              string      `json:"type"`
		ProfileID         string      `json:"profileID"`
		RootCertificate   interface{} `json:"rootCertificate,omitempty"`
		ClientCertificate interface{} `json:"clientCertificate,omitempty"`
		Key               interface{} `json:"publicKey,omitempty"`
	}
)
