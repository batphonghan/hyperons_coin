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
	"fmt"
	"hyperon/wallet"

	"github.com/spf13/cobra"
)

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "main wallet cmd",
}

var creatCmd = &cobra.Command{
	Use:   "create",
	Short: "create a wallet",
	RunE: func(cmd *cobra.Command, args []string) error {
		wallets, err := wallet.CreateWallets()
		if err != nil {
			return fmt.Errorf("err create wallets %v", err)
		}
		address := wallets.AddWallet()
		wallets.SaveFile()

		fmt.Printf("New address is: %s\n", address)

		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		wallets, err := wallet.CreateWallets()
		if err != nil {
			return fmt.Errorf("err create wallet %v", err)
		}
		addresses := wallets.GetAllAddresses()

		for _, address := range addresses {
			fmt.Println(address)
		}

		return nil
	},
}

func init() {
	walletCmd.AddCommand(creatCmd)
	walletCmd.AddCommand(listCmd)

	rootCmd.AddCommand(walletCmd)
}
