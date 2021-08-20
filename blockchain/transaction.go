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
