/*
Copyright © 2026 Andy Yoon
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setApiKeyCmd = &cobra.Command{
	Use:   "set-api-key [KEY]",
	Short: "Set your BALLDONTLIE API key",
	Long: `Sign up for a free API key at https://balldontlie.io, then save it with this command.

Alternatively you can set the WTS_API_KEY environment variable, or set api_key in the wts config file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		viper.Set("api_key", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setApiKeyCmd)
}
