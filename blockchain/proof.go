package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"time"
)

const Difficulty = 12

type ProofOfWork struct {
	Block  Block
	Target *big.Int
}

func NewProof(b Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	pow := ProofOfWork{b, target}
	return &pow
}

func (p *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			p.Block.PrevHash,
			p.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(Difficulty),
		},
		[]byte{},
	)

	return data
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.BigEndian, num); err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte
	var nonce int
	rand.Seed(time.Now().UnixNano())
	for {
		nonce = rand.Int()
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r %x", hash)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		}
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
