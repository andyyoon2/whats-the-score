/*
Copyright © 2026 Andy Yoon
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

// setApiKeyCmd represents the setApiKey command
var setApiKeyCmd = &cobra.Command{
	Use:   "set-api-key",
	Short: "Set your balldontlie API key",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setApiKey called")
	},
}

func init() {
	rootCmd.AddCommand(setApiKeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setApiKeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setApiKeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
