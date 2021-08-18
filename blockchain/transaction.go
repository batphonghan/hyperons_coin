package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
	LH          = "lasthash"
)

type Transaction struct {
	ID      []byte
	Intputs []TxInput
	Outputs []TxOutput
}

func DBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	HandleErr(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

type TxOutput struct {
	//Value would be representative of the amount of coins in a transaction
	Value int
	//The Pubkey is needed to "unlock" any coins within an Output. This indicated that YOU are the one that sent it.
	//You are indentifiable by your PubKey
	//PubKey in this iteration will be very straightforward, however in an actual application this is a more complex algorithm
	PubKey string
}

//TxInput is representative of a reference to a previous TxOutput
type TxInput struct {
	//ID will find the Transaction that a specific output is inside of
	ID []byte

	//Out will be the index of the specific output we found within a transaction.
	//For example if a transaction has 4 outputs, we can use this "Out" field to specify which output we are looking for
	Out int

	//This would be a script that adds data to an outputs' PubKey
	//however for this tutorial the Sig will be indentical to the PubKey.
	Sig string
}

const Reward = 100

func CoinbaseTx(ToAddress, data string) Transaction {
	if data == "" {
		data = fmt.Sprintf("Coin to %s", ToAddress)
	}

	txIn := TxInput{
		ID:  []byte{},
		Out: -1,
		Sig: data,
	}

	txOut := TxOutput{
		Value:  Reward,
		PubKey: ToAddress,
	}

	tx := Transaction{
		ID:      nil,
		Intputs: []TxInput{txIn},
		Outputs: []TxOutput{txOut},
	}
	tx.SetID()
	return tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Intputs) == 1 &&
		len(tx.Intputs[0].ID) == 0 &&
		tx.Intputs[0].Out == -1
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
