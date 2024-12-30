/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"ethwallet/config"
	"ethwallet/pkg/wallet"
	"fmt"

	"github.com/spf13/cobra"
)

// addWalletCmd represents the addWallet command
var addWalletCmd = &cobra.Command{
	Use:   "addWallet",
	Short: "Add an existing wallet to the wallet list",
	Long:  `Add an existing wallet to the wallet list by providing the private key and password.`,
	Run: func(cmd *cobra.Command, args []string) {
		publicKeyHex, seedPhrase, err := wallet.AddWallet(defaultWallet)

		if err != nil {
			fmt.Println("Failed to add wallet: ", err)
			panic(err)
		}

		if publicKeyHex == "" || seedPhrase == "" {
			fmt.Println("Failed to add wallet")
			panic("Failed to add wallet")
		}

		fmt.Println("Wallet added successfully!")
		fmt.Println("Public Address: ", publicKeyHex)
		fmt.Println("Seed Phrase:", seedPhrase)
	},
}

func init() {
	rootCmd.AddCommand(addWalletCmd)

	addWalletCmd.Flags().BoolVarP(&defaultWallet, "default", "d", false, "Add an existing wallet and set it as the default wallet")

	config.LoadViperConfig()

}
