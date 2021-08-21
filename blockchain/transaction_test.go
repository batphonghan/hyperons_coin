package blockchain_test

import (
	"encoding/hex"
	"fmt"
	"hyperon/blockchain"
	"hyperon/wallet"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Sign(t *testing.T) {
	wallets, err := wallet.CreateWallets()
	require.Nil(t, err)

	aliceAddr := wallets.AddWallet()
	bodAddr := wallets.AddWallet()

	alice := wallets.GetWallet(aliceAddr)

	chain := blockchain.InitBlockChain(aliceAddr)
	chain.Wallets = wallets

	tx := blockchain.NewTransaction(aliceAddr, bodAddr, 10, chain)
	fmt.Println(tx.String())

	txs := make(map[string]blockchain.Transaction)

	publicHash, err := wallet.DecodePubKey([]byte(aliceAddr))
	require.Nil(t, err)
	for _, v := range chain.FindUnspentTransactions(publicHash) {
		txs[hex.EncodeToString(v.ID)] = v
	}

	tx.Sign(alice.PrivateKey, txs)
	fmt.Println("After sign")
	fmt.Println(tx.String())

	require.True(t, tx.Verify(txs))
}
