package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Wallets struct {
	Wallets map[string]Wallet
}

func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]Wallet)
	err := wallets.LoadFile()

	return &wallets, err
}

func (w *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	w.Wallets[address] = wallet

	return address
}

func (w *Wallets) GetWallet(address string) Wallet {
	return w.Wallets[address]
}

func (w *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range w.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

func (w *Wallets) SaveFile() {
	var content bytes.Buffer
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(w)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)

	if err != nil {
		log.Panic(err)
	}
}

func (w *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		return err
	}

	w.Wallets = wallets.Wallets

	return nil
}
