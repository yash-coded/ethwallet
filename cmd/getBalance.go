/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"ethwallet/pkg/wallet"
	"log"

	"github.com/spf13/cobra"
)

// getBalanceCmd represents the getBalance command
var getBalanceCmd = &cobra.Command{
	Use:   "getBalance",
	Short: "Get the balance of an ethereum wallet",
	Long:  `Get the balance of an ethereum wallet by providing the public key of the wallet.`,
	Run: func(cmd *cobra.Command, args []string) {
		publicKey, err := cmd.Flags().GetString("publicKey")

		if err != nil {
			log.Fatalf("Failed to get public key: %v", err)
			panic(err)
		}

		if publicKey == "" {
			log.Fatalf("Public key is required")
			panic("Public key is required")
		}

		ethvalue, err := wallet.GetBalance(publicKey)

		if err != nil {
			log.Fatalf("Failed to get balance: %v", err)
			panic(err)
		}

		log.Printf("Balance of %s is %f ETH", publicKey, ethvalue)
	},
}

func init() {
	rootCmd.AddCommand(getBalanceCmd)

	getBalanceCmd.Flags().StringP("publicKey", "p", "", "Public key of the wallet")
}
