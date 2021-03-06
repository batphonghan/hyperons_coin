package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash         []byte
	Transactions []Transaction
	PrevHash     []byte
	Nonce        int
}

func (b *Block) HashTransactions() []byte {
	var txsHash [][]byte
	for _, tx := range b.Transactions {
		txsHash = append(txsHash, tx.ID)
	}
	info := bytes.Join(txsHash, []byte{})
	hash := sha256.Sum256(info)
	return hash[:]
}

func Genesis(coinbase Transaction) Block {
	return CreateBlock([]Transaction{coinbase}, []byte{})
}

func CreateBlock(txs []Transaction, prevHash []byte) Block {
	block := Block{[]byte{}, txs, prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	if err := encoder.Encode(b); err != nil {
		log.Panic(err)
	}

	return res.Bytes()
}

func Deserialize(data []byte) Block {
	var b Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&b); err != nil {
		log.Panic(err)
	}

	return b
}
