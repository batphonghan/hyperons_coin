package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Proof(t *testing.T) {
	v := big.NewInt(1)
	fmt.Printf("before v is: %b ", v)
	v.Lsh(v, 256-(256-10))
	fmt.Printf("after v is: %b ", v)
	fmt.Println(ToHex(1 << 32))
}

func Test_Me(t *testing.T) {
	tx := Transaction{
		ID:      []byte{},
		Intputs: []TxInput{TxInput{}},
		Outputs: []TxOutput{},
	}

	var encode bytes.Buffer
	encoder := gob.NewEncoder(&encode)
	err := encoder.Encode(tx)
	require.Nil(t, err)
	b := encode.Bytes()

	_b := sha256.Sum256(b)

	var intHash big.Int
	intHash.SetBytes(_b[:])

	//864691128522211840
	fmt.Println(intHash.Int64())

}
