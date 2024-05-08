package entity

type TLSCerts struct {
	RootCertificate   CertCreationResult
	IssuedCertificate CertCreationResult
	Version           string
}

type CertCreationResult struct {
	H             string
	Cert          string
	Pem           string
	CertBin       string
	PrivateKey    string
	PrivateKeyBin string
	Checked       bool
	Key           []byte
}
