/*
Copyright © 2026 Andy Yoon
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setApiKeyCmd represents the setApiKey command
var setApiKeyCmd = &cobra.Command{
	Use:   "set-api-key [KEY]",
	Short: "Set your balldontlie API key",
	Long: `Sign up for a free API key at https://balldontlie.io, then save it with this command.

Alternatively you can set the WTS_API_KEY environment variable or edit the wts config file.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("api_key", args[0])
	},
}

func init() {
	rootCmd.AddCommand(setApiKeyCmd)
}
