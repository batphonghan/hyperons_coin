/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"hyperon/blockchain"
	"hyperon/wallet"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
func start(cmd *cobra.Command, args []string) error {
	defer func() {
		if err := recover(); err != nil {
			log.Println("err ", err)
		}
	}()
	http.HandleFunc("/dump", dump)
	http.HandleFunc("/balance", getBalance)
	http.HandleFunc("/send", send)
	return http.ListenAndServe(":9787", nil)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start hyperchain serve at :9787",
	RunE:  start,
}

func init() {
	rootCmd.AddCommand(startCmd)

	add := startCmd.Flags().StringP("address", "a", "0x0", "The address to send genesis block reward to")
	createBlockChain(*add)
}

func dump(rw http.ResponseWriter, r *http.Request) {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	var blocks []blockchain.Block
	for {
		block := iter.Next()
		blocks = append(blocks, block)

		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
	b, err := json.Marshal(blocks)
	if err != nil {
		responseErr(rw, err)
		return
	}

	responseOk(rw, b)

}

func responseErr(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)

	fmt.Fprintf(rw, "Server err: %+v", err)
}

func responseOk(rw http.ResponseWriter, b []byte) {
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-type", "application/json")

	fmt.Fprintf(rw, string(b))
}

func createBlockChain(address string) {
	chain := blockchain.InitBlockChain(address)
	chain.Database.Close()
	fmt.Println("Finished!")
}

func getBalance(rw http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	chain := blockchain.ContinueBlockChain(address)
	defer chain.Database.Close()

	balance := 0
	publishHash, err := wallet.DecodePubKey([]byte(address))
	if err != nil {
		panic(err)
	}
	UTXOs := chain.FindUTXO(publishHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Fprintf(rw, "Balance of %s: %d\n", address, balance)
}

func send(rw http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from := q.Get("from")
	amount, err := strconv.Atoi(q.Get("amount"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	to := q.Get("to")
	chain := blockchain.ContinueBlockChain(from)
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]blockchain.Transaction{tx})
	fmt.Println("Success!")
}
