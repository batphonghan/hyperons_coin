package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewKeyPair() (*ecdsa.PrivateKey, []byte) {
	gob.Register(elliptic.P256())
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	HandleErr(err)

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return private, pub
}
