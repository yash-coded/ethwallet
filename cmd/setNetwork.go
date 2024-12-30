/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"ethwallet/config"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setNetworkCmd represents the setNetwork command
var setNetworkCmd = &cobra.Command{
	Use:   "setNetwork",
	Short: "Set the network to connect to",
	Long:  `Set the network to connect to by providing the network URL`,
	Run: func(cmd *cobra.Command, args []string) {
		url, err := cmd.Flags().GetString("url")

		if err != nil {
			log.Fatalf("Failed to get network URL: %v", err)
		}

		viper.Set("url", url)

		config.WriteViperConfig()

		log.Printf("Network URL set to %s", url)
	},
}

func init() {
	rootCmd.AddCommand(setNetworkCmd)

	setNetworkCmd.Flags().StringP("url", "u", "", "URL of the network to connect to")
	setNetworkCmd.MarkFlagRequired("url")

	config.LoadViperConfig()
}
