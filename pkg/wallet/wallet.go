package wallet

import (
	"context"
	"crypto/sha256"
	"ethwallet/config"
	"ethwallet/pkg/encryption"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
	"github.com/tyler-smith/go-bip39"
)

func CreateWallet(isDefault bool) (publicKeyHex string, seedPhrase string, err error) {
	prompt := promptui.Prompt{
		Label: "Enter Password",
		Mask:  '*',
	}

	password, err := prompt.Run()
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)

		return "", "", err
	}

	log.Println("Creating new wallet with password: ", password)

	privateKey, err := crypto.GenerateKey()

	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)

		return "", "", err
	}

	privateKeyHex := fmt.Sprintf("%x", privateKey.D.Bytes())
	publicAddressHex := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	entropy := sha256.Sum256([]byte(privateKeyHex))

	mnemonic, err := bip39.NewMnemonic(entropy[:])

	if err != nil {
		log.Fatalf("Failed to generate mnemonic: %v", err)

		return "", "", err
	}

	config.SaveWalletInfo(publicAddressHex, privateKeyHex, password, isDefault)

	return publicAddressHex, mnemonic, nil

}

func AddWallet(isDefault bool) (publicAddressHex string, seedPhrase string, err error) {
	prompt := promptui.Prompt{
		Label: "Enter Private Key",
		Mask:  '*',
	}

	privateKeyHex, err := prompt.Run()

	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
		return "", "", err
	}

	prompt = promptui.Prompt{
		Label: "Enter Password",
		Mask:  '*',
	}

	password, err := prompt.Run()

	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
		return "", "", err
	}

	// encrypt the private key with the password

	encryptedKey, err := encryption.EncryptData([]byte(privateKeyHex), password)

	if err != nil {
		log.Fatalf("Failed to encrypt private key: %v", err)
		return "", "", err
	}

	// save the encrypted key to the wallet list with viper

	var wallets []map[string]string

	if err := viper.UnmarshalKey("wallets", &wallets); err != nil {
		log.Fatalf("Failed to unmarshal wallets: %v", err)
		return "", "", err
	}

	privateKey, err := crypto.HexToECDSA(encryption.ConvertHexAddress(privateKeyHex))

	if err != nil {
		log.Fatalf("Failed to convert hex address: %v", err)
		return "", "", err
	}

	entropy := sha256.Sum256([]byte(privateKeyHex))
	mnemonic, err := bip39.NewMnemonic(entropy[:])
	if err != nil {
		log.Fatalf("Failed to generate mnemonic: %v", err)
	}

	publicKeyHex := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	wallet := map[string]string{
		"privateKey":    encryptedKey,
		"publicAddress": publicKeyHex,
	}

	wallets = append(wallets, wallet)

	viper.Set("wallets", wallets)

	if isDefault {
		viper.Set("defaultWallet", publicKeyHex)
	}

	config.WriteViperConfig()

	return publicKeyHex, mnemonic, nil
}

func GetBalance(publicAddress string) (balanceEth float64, err error) {
	client, err := ethclient.Dial(config.GetRPCUrl())

	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		panic(err)
	}

	balance, err := client.BalanceAt(context.Background(), common.HexToAddress(publicAddress), nil)

	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
		panic(err)
	}

	balanceInEther := new(big.Float)
	balanceInEther.SetString(balance.String())

	ethValue := new(big.Float).Quo(balanceInEther, big.NewFloat(1e18))

	ethValueFloat, _ := ethValue.Float64()

	return ethValueFloat, nil
}
