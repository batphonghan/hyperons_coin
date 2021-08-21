/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

const (
	ChecksumLength = 4
	version        = byte(0x00)
	walletFile     = "./tmp/wallet.data"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func makeWallet() Wallet {
	private, public := newKeyPair()
	return Wallet{
		PrivateKey: private,
		PublicKey:  public,
	}
}

func newKeyPair() (*ecdsa.PrivateKey, []byte) {
	gob.Register(elliptic.P256())
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return private, pub
}

func Base58Encode(b []byte) []byte {
	s := base58.Encode(b)
	return []byte(s)
}

func Base58Decode(b []byte) ([]byte, error) {
	return base58.Decode(string(b))
}

// walletCmd represents the wallet command

// 1. Take our public key (in bytes)
// 2. Run a SHA256 hash on it, then run a RipeMD160 hash on that hash. This is called our PublicHash
// 3. Take our Publish Hash and append our Version (the global variable from earlier) to it. This is called the Versioned Hash
// 4. Run SHA256 on our Versioned Hash twice. Then take the first 4 bytes of that output. This is called the Checksum
// 5. Then we will add our Checksum to the end of our original Versioned Hash. We can call this FinalHash
// 6. Lastly, we will base58Encode our FinalHash. This is our wallet address!

func PublicKeyHash(publicKey []byte) []byte {
	hashPubKey := sha256.Sum256(publicKey)

	hasher := ripemd160.New()
	if _, err := hasher.Write(hashPubKey[:]); err != nil {
		log.Panic(err)
	}

	pubicRipeMd := hasher.Sum(nil)

	return pubicRipeMd
}

func CheckSum(versionedHash []byte) []byte {
	firstHash := sha256.Sum256(versionedHash)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:ChecksumLength]
}

func (w *Wallet) Address() []byte {
	publicHash := PublicKeyHash(w.PublicKey)

	// Step 3
	versionedHash := append([]byte{version}, publicHash...)
	// Step 4
	checksum := CheckSum(versionedHash)
	// Step 5
	finalHash := append(versionedHash, checksum...)

	// Step 6
	address := Base58Encode(finalHash)

	return address
}

func DecodePubKey(address []byte) ([]byte, error) {
	publicHash, err := Base58Decode([]byte(address))
	if err != nil {
		return nil, fmt.Errorf("Err decode address: %x", address)
	}

	pubkey := publicHash[1 : len(publicHash)-ChecksumLength]

	return pubkey, nil
}

func ValidateAddress(address string) bool {
	publicHash, err := Base58Decode([]byte(address))
	if err != nil {
		fmt.Println("Can not decode address")
		return false
	}

	checksum := publicHash[len(publicHash)-ChecksumLength:]
	version := publicHash[0]

	pubkey := publicHash[1 : len(publicHash)-ChecksumLength]
	actualCheckSum := CheckSum(append([]byte{version}, pubkey...))

	return bytes.Compare(checksum, actualCheckSum) == 0
}
