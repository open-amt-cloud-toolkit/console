package entity

import "crypto/x509"

type TLSCerts struct {
	RootCertificate   CertCreationResult `json:"rootCertificate"`
	IssuedCertificate CertCreationResult `json:"issuedCertificate"`
	Version           string             `json:"version"`
}

type CertCreationResult struct {
	H             string           `json:"h:"`
	Cert          x509.Certificate `json:"cert"`
	Pem           string           `json:"pem"`
	CertBin       string           `json:"certBin"`
	PrivateKey    string           `json:"privateKey"`
	PrivateKeyBin string           `json:"privateKeyBin"`
	Checked       bool             `json:"checked" example:"true"`
	Key           []byte           `json:"key"`
}
