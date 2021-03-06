package blockchain

import (
	"encoding/hex"
	"fmt"
	"hyperon/wallet"
	"log"
	"os"
	"runtime"

	"github.com/dgraph-io/badger/v3"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
	Wallets  *wallet.Wallets
}
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

var GenesisData = "123"

func InitBlockChain(address string) *BlockChain {
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	HandleErr(err)

	var lasthash []byte
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(LH)); err == badger.ErrKeyNotFound {
			fmt.Println("Not Exist DB")
			cbTx := CoinbaseTx(address, GenesisData)
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

func (iter *BlockChainIterator) Next() Block {
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

	return block
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

func DBExist() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
func ContinueBlockChain(address string) *BlockChain {
	if DBExist() == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	HandleErr(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(LH))
		HandleErr(err)
		lastHash, err = item.ValueCopy(nil)

		return err
	})
	HandleErr(err)

	chain := BlockChain{
		LastHash: lastHash,
		Database: db,
	}

	return &chain
}

func (chain *BlockChain) FindUTXO(pubkeyHash []byte) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(pubkeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(pubkeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

func (chain *BlockChain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					return accumulated, unspentOuts
				}
			}
		}
	}

	return accumulated, unspentOuts
}

func (chain *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	i := 0
	for {
		i++
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(pubKeyHash) {
					unspentTxs = append(unspentTxs, tx)
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlock(string(pubKeyHash)) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	fmt.Println("Time loop: ", i)
	return unspentTxs
}
