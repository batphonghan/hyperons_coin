package blockchain

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v3"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

var Add = "123"
var GenesisData = "123"

func InitBlockChain() *BlockChain {
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	HandleErr(err)

	var lasthash []byte
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(LH)); err == badger.ErrKeyNotFound {
			fmt.Println("Not Exist DB")
			cbTx := CoinbaseTx(Add, GenesisData)
			gen := Genesis(cbTx)

			err = txn.Set(gen.Hash, gen.Serialize())
			if err != nil {
				return fmt.Errorf("set hash data err: %v", err)
			}

			err = txn.Set([]byte(LH), gen.Hash)
			if err != nil {
				return fmt.Errorf("set lasthast key err: %v", err)
			}
			lasthash = gen.Hash

			return nil
		} else {
			item, err := txn.Get([]byte(LH))
			if err != nil {
				return fmt.Errorf("Get lasthast key err: %v", err)
			}
			lasthash, err = item.ValueCopy(nil)
			HandleErr(err)
			return nil
		}
	})

	HandleErr(err)

	return &BlockChain{
		LastHash: lasthash,
		Database: db,
	}
}

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		HandleErr(err)
		var encodedBlock []byte
		encodedBlock, err = item.ValueCopy(nil)
		block = Deserialize(encodedBlock)

		return err
	})
	HandleErr(err)

	iter.CurrentHash = block.PrevHash

	return &block
}

func (chain *BlockChain) AddBlock(txs []Transaction) {
	var lastHash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(LH))
		HandleErr(err)
		lastHash, err = item.ValueCopy(nil)
		HandleErr(err)
		return nil
	})

	HandleErr(err)

	newBlock := CreateBlock(txs, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = txn.Set([]byte(LH), newBlock.Hash)
		if err != nil {
			return err
		}
		return nil
	})
	HandleErr(err)
}
