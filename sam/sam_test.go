package sam

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sam(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	handleErr(err)

	msg := "Hello hung nguyen"

	hash := sha256.Sum256([]byte(msg))
	fmt.Printf("Len hash : %d \n", len(hash))

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	handleErr(err)

	fmt.Printf("Signature: %x, len %d \n", sig, len(sig))

	valid := ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], sig)
	fmt.Println("Verify ASN1 IS valid: ", valid)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	handleErr(err)

	valid = ecdsa.Verify(&privateKey.PublicKey, hash[:], r, s)
	fmt.Println("Verify IS valid: ", valid)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Test_en(t *testing.T) {
	sum := sha256.Sum256([]byte("This is hung"))
	fmt.Printf("%x == len %d \n", sum, len(sum)*8)
}

func Test_Symmetric(t *testing.T) {
	key := sha256.Sum256([]byte("This is hung"))
	plainText := []byte("This is plain text")

	block, err := aes.NewCipher(key[:])
	handleErr(err)

	aesgcm, err := cipher.NewGCM(block)
	handleErr(err)
	nonce := []byte("gopostmedium")
	ciphertext := aesgcm.Seal(nil, nonce, plainText, nil)

	text, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	handleErr(err)

	assert.Equal(t, plainText, text)
}
