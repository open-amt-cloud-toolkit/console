package dto

type TLSCerts struct {
	RootCertificate   CertCreationResult `json:"rootCertificate"`
	IssuedCertificate CertCreationResult `json:"issuedCertificate"`
	Version           string             `json:"version"`
}

type CertCreationResult struct {
	H             string `json:"h:"`
	Cert          string `json:"cert"`
	PEM           string `json:"pem"`
	CertBin       string `json:"certBin"`
	PrivateKey    string `json:"privateKey"`
	PrivateKeyBin string `json:"privateKeyBin"`
	Checked       bool   `json:"checked" example:"true"`
	Key           []byte `json:"key"`
}
