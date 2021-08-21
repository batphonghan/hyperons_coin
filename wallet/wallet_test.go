package wallet_test

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"hyperon/wallet"
	"os"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func Test_ValidateAddress(t *testing.T) {
	err := wallet.Delete()
	require.Nil(t, err)

	w, err := wallet.CreateWallets()
	require.Nil(t, err)
	require.NotNil(t, w)

	w.AddWallet()
	// w.SaveFile()

	// addrs := w.GetAllAddresses()

	bitAdd, err := os.Open("bit.txt")
	require.Nil(t, err)
	sc := bufio.NewScanner(bitAdd)

	c := 0
	for sc.Scan() {
		add := sc.Text()
		c++
		fmt.Println("Validate: ", add)
		require.True(t, wallet.ValidateAddress(add))
		fmt.Println("OK =>", c)
	}
}

func Test_ur(t *testing.T) {
	a := []byte{0xff}

	v := utf8.Valid(a)
	fmt.Println(v)
	s := string(a)

	v = utf8.ValidString(s)
	fmt.Println(v)
	fmt.Println(s)
	for _, r := range s {
		fmt.Println(r)
	}
	rs := []rune(s)
	fmt.Println(rs)
}

func Test_new(t *testing.T) {
	gob.Register(elliptic.P256())
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	fmt.Printf("%x \n", pub)

	h := hex.EncodeToString(pub)

	fmt.Println(h)
}

//1GkQmKAmHtNfnD3LHhTkewJxKHVSta4m2a 50BTC

func Test_ReipMD(t *testing.T) {
	add := "1JSU7TXuEqrEFD5aUr5BDs6yPBPW7THYkE"
	decoded, err := wallet.Base58Decode([]byte(add))
	require.Nil(t, err)
	ripMD := decoded[1 : len(decoded)-wallet.ChecksumLength]
	s := hex.EncodeToString(ripMD)
	require.Equal(t, "bf4b2b34a6e1a57393820c737916b070ee0c7c3a", s)
}
