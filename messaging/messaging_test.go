package messaging

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"
)

const (
	TestPublicKey string = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAElk30LFnrF48XLeEHrG3K/r7215xg
gOEmGeRDdJ7f86ByD7uK/Jxje79Jtn9HNjyQahd7bBBKUOfcWG3Kh927oA==
-----END PUBLIC KEY-----`

	// NB: make sure to use SPACES here in the test message instead of tabs,
	// otherwise validation will fail
	TestMessage string = `{
    "message": {
        "type": "issueTx",
        "userId": 1,
        "transaction": {
            "amount": 10123.50
        }
    },
    "signature": "MEUCIBkooxG2uFZeSEeaf5Xh5hWLxcKGMxCZzfnPshOh22y2AiEAwVLAaGhccUv8UhgC291qNWtxrGawX2pPsI7UUA/7QLM="
}`
)

func loadPublicKey(publicKey string) (*ecdsa.PublicKey, error) {
	// decode the key, assuming it's in PEM format
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("Failed to decode PEM public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("Failed to parse ECDSA public key")
	}
	switch pub := pub.(type) {
	case *ecdsa.PublicKey:
		return pub, nil
	}
	return nil, errors.New("Unsupported public key type")
}

func TestEnvelopeValidation(t *testing.T) {
	// our test message
	envelope, err := NewEnvelopeFromJSON(TestMessage)
	if err != nil {
		t.Error("Expected to be able to deserialise test message, but failed with err =", err)
	}
	// extract the public key from the test key string
	publicKey, err := loadPublicKey(TestPublicKey)
	if err != nil {
		t.Error("Failed to parse test public key:", err)
	}
	// now we validate the signature against the public key
	if err := envelope.Validate(publicKey); err != nil {
		t.Error("Expected nil error from message envelope validation routine, but got:", err)
	}
}
