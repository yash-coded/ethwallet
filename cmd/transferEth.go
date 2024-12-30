/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"ethwallet/config"
	"ethwallet/pkg/encryption"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// transferEthCmd represents the transferEth command
var transferEthCmd = &cobra.Command{
	Use:   "transferEth",
	Short: "Transfer ethereum from one wallet to another",
	Long:  `Transfer ethereum from one wallet to another by providing the sender's password, receiver's public key, and the amount of ethereum to transfer.`,
	Run:   transferEth,
}

func init() {
	rootCmd.AddCommand(transferEthCmd)

	config.LoadViperConfig()

}

func transferEth(cmd *cobra.Command, args []string) {

	prompt := promptui.Prompt{
		Label: "Enter Receiver's Public Key",
	}

	receiverPublicKey, _ := prompt.Run()

	prompt = promptui.Prompt{
		Label: "Enter Sender wallet Password",
		Mask:  '*',
	}

	senderPassword, _ := prompt.Run()

	prompt = promptui.Prompt{
		Label: "Enter Amount of ETH to Transfer",
	}

	ethAmountStr, _ := prompt.Run()

	ethAmount, err := strconv.ParseFloat(ethAmountStr, 64)

	if err != nil {
		panic(err)
	}

	amount := big.NewFloat(ethAmount)

	// transfer the ethereum

	var wallets []map[string]string

	if err := viper.UnmarshalKey("wallets", &wallets); err != nil {
		panic(err)
	}

	defaultWallet := viper.GetString("defaultWallet")

	for _, wallet := range wallets {
		if wallet["publicaddress"] == defaultWallet {
			privateKeyHex, err := encryption.DecryptData(wallet["privatekey"], senderPassword)

			if err != nil {
				log.Fatalf("Failed to decrypt private key: %v", err)
				panic(err)
			}

			client, err := ethclient.Dial(config.GetRPCUrl())

			if err != nil {
				log.Fatalf("Failed to connect to the Ethereum network: %v", err)
				panic(err)
			}

			nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(wallet["publickey"]))

			if err != nil {
				log.Fatalf("Failed to get nonce: %v", err)
				panic(err)
			}

			gasPrice, err := client.SuggestGasPrice(context.Background())

			if err != nil {
				log.Fatalf("Failed to get gas price: %v", err)
				panic(err)
			}

			wei := new(big.Int)
			amountWei, _ := amount.Mul(amount, big.NewFloat(1e18)).Int(wei)

			tx := types.NewTransaction(nonce, common.HexToAddress(receiverPublicKey), amountWei, uint64(21000), gasPrice, nil)

			prompt := promptui.Prompt{
				Label:     "Confirm Transfer",
				IsConfirm: true,
			}

			fmt.Printf("Transfer %f ETH to %s\n", ethAmount, receiverPublicKey)
			fmt.Println("Gas Price: ", gasPrice)
			fmt.Println("Nonce: ", nonce)
			fmt.Println("Max Gas: ", 21000)

			_, err = prompt.Run()

			if err != nil {
				log.Println("Transfer cancelled")
				return
			}

			chainID, err := client.NetworkID(context.Background())

			if err != nil {
				log.Fatalf("Failed to get chain ID: %v", err)
				panic(err)
			}

			privateKey, err := crypto.HexToECDSA(encryption.ConvertHexAddress(privateKeyHex))

			if err != nil {
				log.Fatalf("Failed to convert private key to ECDSA: %v", err)
				panic(err)
			}

			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

			if err != nil {
				log.Fatalf("Failed to sign transaction: %v", err)
				panic(err)
			}

			err = client.SendTransaction(context.Background(), signedTx)

			if err != nil {
				log.Fatalf("Failed to send transaction: %v", err)
				panic(err)
			}

			log.Printf("Transfer successful! Tx Hash: %s", signedTx.Hash().Hex())
			break
		}
	}

}
