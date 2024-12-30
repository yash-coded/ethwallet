package config

import (
	"ethwallet/pkg/encryption"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func LoadViperConfig() {
	home, err := os.UserHomeDir()

	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	viper.AddConfigPath(home)
	viper.SetConfigName(".ethereum-wallet")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error since we will create a new one
		} else {
			log.Fatalf("Failed to read config file: %v", err)
		}
	}
}

func WriteViperConfig() {
	home, err := os.UserHomeDir()

	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	configFile := filepath.Join(home, ".ethereum-wallet.json")

	if err = viper.WriteConfigAs(configFile); err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}
}

func SaveWalletInfo(publicAddressHex, privateKeyHex, password string, isDefault bool) {

	encryptedPrivateKey, err := encryption.EncryptData([]byte(privateKeyHex), password)

	if err != nil {
		log.Fatalf("Failed to encrypt private key: %v", err)
	}

	var wallets []map[string]string

	if err := viper.UnmarshalKey("wallets", &wallets); err != nil {
		log.Fatalf("Failed to unmarshal wallets: %v", err)
	}

	wallet := map[string]string{
		"publicAddress": publicAddressHex,
		"privateKey":    encryptedPrivateKey,
	}

	wallets = append(wallets, wallet)

	viper.Set("wallets", wallets)

	if isDefault {
		viper.Set("defaultWallet", publicAddressHex)
	}

	WriteViperConfig()
}

func GetRPCUrl() string {
	url := viper.GetString("url")

	if url == "" {
		url = "http://localhost:8545"
	}

	return url
}
