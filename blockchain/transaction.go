package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"hyperon/wallet"
	"log"
	"math/big"
	"os"
	"strings"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
	LH          = "lasthash"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func DBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
func (tx *Transaction) setID() {
	var encoded bytes.Buffer
	var hash [32]byte
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	HandleErr(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

const Reward = 100

func CoinbaseTx(address, data string) Transaction {
	if data == "" {
		data = fmt.Sprintf("Coin to %s", address)
	}

	txIn := TxInput{
		ID:        []byte{},
		Out:       -1,
		Signature: []byte(data),
	}

	txOut := NewTXOutput(Reward, address)

	tx := Transaction{
		ID:      nil,
		Inputs:  []TxInput{txIn},
		Outputs: []TxOutput{txOut},
	}
	tx.setID()
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
	w := chain.Wallets.GetWallet(from)
	publicHash, err := wallet.DecodePubKey([]byte(from))
	HandleErr(err)
	acc, validOutputs := chain.FindSpendableOutputs(publicHash, amount)

	if acc < amount {
		log.Panic("Error: not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		HandleErr(err)

		for _, out := range outs {
			input := TxInput{
				ID:        txID,
				Out:       out,
				Signature: []byte{},
				PubKey:    w.PublicKey,
			}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, NewTXOutput(amount, to))

	if acc > amount {
		outputs = append(outputs, NewTXOutput(acc-amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.setID()

	return tx
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

func (tx *Transaction) Serialize() []byte {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func (tx *Transaction) Sign(privKey *ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}
	txCopy := tx.TrimmedCopy()
	for inId, input := range tx.Inputs {
		hexID := hex.EncodeToString(input.ID)
		if prevTxs[hexID].ID == nil {
			log.Panic("ERROR: previous transaction is not correct")
		}

		// set state
		prevTX := prevTxs[hex.EncodeToString(input.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTX.Outputs[input.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inId].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, privKey, txCopy.ID)

		fmt.Printf("Sign: %x r: %s s: %s\n", txCopy.ID, r.String(), s.String())
		HandleErr(err)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[inId].Signature = signature

		pubKeyHash := input.PubKey
		var x, y big.Int
		publen := len(pubKeyHash)
		x.SetBytes(pubKeyHash[:publen/2])
		y.SetBytes(pubKeyHash[publen/2:])

		rawPubKey := ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     &x,
			Y:     &y,
		}
		if !ecdsa.Verify(&rawPubKey, txCopy.ID, r, s) {
			panic("AAAAA")
		} else {
			fmt.Println("Okay ")
		}
	}
}

func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	for _, in := range tx.Inputs {
		IDToFind := hex.EncodeToString(in.ID)
		if prevTxs[IDToFind].ID == nil {
			log.Panic("Previous tx is not correct ID:", IDToFind)
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inId, in := range tx.Inputs {
		hexID := hex.EncodeToString(in.ID)
		prevTxs := prevTxs[hexID]

		// Clear
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTxs.Outputs[in.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inId].PubKey = nil

		var r, s big.Int
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:sigLen/2])
		s.SetBytes(in.Signature[sigLen/2:])

		var x, y big.Int
		publen := len(in.PubKey)
		x.SetBytes(in.PubKey[:publen/2])
		y.SetBytes(in.PubKey[publen/2:])

		rawPubKey := ecdsa.PublicKey{
			Curve: curve,
			X:     &x,
			Y:     &y,
		}
		fmt.Printf("Verify: %x r: %s s: %s\n", txCopy.ID, r.String(), s.String())

		match := ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s)
		if !match {
			return false
		}
	}
	return true
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{in.ID, in.Out, nil, nil})
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:     %x", input.ID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Out))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}
