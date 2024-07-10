package dto

type (
	SecuritySettings struct {
		ProfileAssociation []ProfileAssociation `json:"ProfileAssociation"`
		Certificates       interface{}          `json:"Certificates"`
		Keys               interface{}          `json:"PublicKeys"`
	}

	ProfileAssociation struct {
		Type              string      `json:"Type"`
		ProfileID         string      `json:"ProfileID"`
		RootCertificate   interface{} `json:"RootCertificate,omitempty"`
		ClientCertificate interface{} `json:"ClientCertificate,omitempty"`
		Key               interface{} `json:"PublicKey,omitempty"`
	}
)
