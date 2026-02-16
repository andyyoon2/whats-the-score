/*
Copyright © 2026 Andy Yoon
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/andyyoon2/whats-the-score/lib"
	"github.com/spf13/cobra"
)

func get() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.balldontlie.io/v1/teams/", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", os.Getenv("WTS_API_KEY"))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		teams := lib.GetTeams()
		if len(args) == 0 {
			fmt.Println(teams)
			return
		}

		query := strings.ToLower(args[0])
		for _, t := range teams {
			if strings.ToLower(t.Name) == query || strings.ToLower(t.City) == query || strings.ToLower(t.Abbreviation) == query {
				fmt.Println(t)
				fmt.Println(lib.GetGamesForTeam(t))
				return
			}
		}

		fmt.Println("Team not found!")
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
