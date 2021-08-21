package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hyperon/wallet"
)

type TxOutput struct {
	Value int

	PubKeyHash []byte
}

//TxInput is representative of a reference to a previous TxOutput
type TxInput struct {
	ID  []byte
	Out int

	Signature []byte
	PubKey    []byte
}

func NewTXOutput(value int, address string) TxOutput {
	txo := TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

func (in *TxInput) CanUnlock(data string) bool {
	return string(in.Signature) == data
}

func (out *TxOutput) CanBeUnlocked(pubkeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubkeyHash) == 0
}

// Check if input use public hashed key
func (in *TxInput) IsUsesKey(publicKeyHashed []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Compare(lockingHash, publicKeyHashed) == 0
}

func (out *TxOutput) Lock(address []byte) error {
	pubKeyHash, err := wallet.DecodePubKey(address)
	if err != nil {
		return fmt.Errorf("Err during decode address %v", err)
	}

	out.PubKeyHash = pubKeyHash
	return nil
}

func NewTxOutput(value int, address string) (TxOutput, error) {
	txo := TxOutput{
		Value: value,
	}
	err := txo.Lock([]byte(address))
	if err != nil {
		return txo, fmt.Errorf("err lock txo %v", err)
	}

	return txo, nil
}

func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer

	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	HandleErr(err)

	return buffer.Bytes()
}

type TxOutputs struct {
	Outputs []TxOutput
}

func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs

	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&outputs)
	HandleErr(err)

	return outputs
}
