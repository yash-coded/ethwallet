/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"ethwallet/config"
	"ethwallet/pkg/wallet"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var defaultWallet bool

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new ethereum wallet",
	Long:  `Create a new ethereum wallet with a new private key and public key.`,
	Run: func(cmd *cobra.Command, args []string) {
		publicKeyHex, seedPhrase, err := wallet.CreateWallet(defaultWallet)

		if err != nil {
			log.Fatalf("Failed to create wallet: %v", err)
			panic(err)
		}

		if publicKeyHex == "" || seedPhrase == "" {
			log.Fatalf("Failed to create wallet")
			panic("Failed to create wallet")
		}

		fmt.Println("🎉 New Wallet Created!")
		fmt.Println("-----------------------")
		fmt.Printf("Public Address: %s\n", publicKeyHex)
		fmt.Printf("Seed Phrase: %s\n", seedPhrase)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVarP(&defaultWallet, "default", "d", false, "Create a new wallet and set it as the default wallet")

	// load viper config

	config.LoadViperConfig()
}
