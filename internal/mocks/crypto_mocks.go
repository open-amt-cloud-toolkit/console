package mocks

import (
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/v2/pkg/config"
	"gopkg.in/yaml.v2"
)

type MockCrypto struct{}

// Encrypt encrypts a string.
func (c MockCrypto) Encrypt(_ string) (string, error) {
	return "encrypted", nil
}

func (c MockCrypto) GenerateKey() string {
	return "key"
}

func (c MockCrypto) Decrypt(_ string) (string, error) {
	return "decrypted", nil
}

// Read encrypted data from file and decrypt it.
func (c MockCrypto) ReadAndDecryptFile(_ string) (config.Configuration, error) {
	var configuration config.Configuration

	err := yaml.Unmarshal([]byte(""), &configuration)
	if err != nil {
		return config.Configuration{}, err
	}

	return configuration, nil
}
