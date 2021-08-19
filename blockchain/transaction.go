package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
	LH          = "lasthash"
)

type Transaction struct {
	ID     []byte
	Inputs []TxInput
	Ouputs []TxOutput
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
		ID:     nil,
		Inputs: []TxInput{txIn},
		Ouputs: []TxOutput{txOut},
	}
	tx.SetID()
	return tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 &&
		len(tx.Inputs[0].ID) == 0 &&
		tx.Inputs[0].Out == -1
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}

func NewTransaction(from, to string, amount int, chain *BlockChain) Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("Error: not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		HandleErr(err)

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return tx
}
