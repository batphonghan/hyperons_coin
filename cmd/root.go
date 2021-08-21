package cmd

import (
	"fmt"
	"hyperon/wallet"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hyperon",
	Short: "hyperon coin deamon",
	Long:  "hyperon coin deamon",
}

func init() {
	walletCmd.AddCommand(creatCmd)
	walletCmd.AddCommand(listCmd)

	rootCmd.AddCommand(walletCmd)
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

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
