/*
Copyright © 2026 Andy Yoon
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/andyyoon2/whats-the-score/lib"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func renderGame(g lib.Game) {
	boldStyle := lipgloss.NewStyle().Bold(true)
	homeScore := fmt.Sprintf("%-3d", g.HomeTeamScore) // left-justify home score
	visitorScore := fmt.Sprintf("%3d", g.VisitorTeamScore)
	if g.Status == "Final" {
		if g.HomeTeamScore > g.VisitorTeamScore {
			homeScore = boldStyle.Render(homeScore)
		} else {
			visitorScore = boldStyle.Render(visitorScore)
		}
	}

	fmt.Printf("%s: %s %s - %s %s\n", g.Date, g.VisitorTeam.Abbreviation, visitorScore, homeScore, g.HomeTeam.Abbreviation)
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
		var games []lib.Game
		teams := lib.GetTeams()
		if len(args) == 0 {
			games = lib.GetGames()
		} else {
			query := strings.ToLower(args[0])
			for _, t := range teams {
				if strings.ToLower(t.Name) == query || strings.ToLower(t.City) == query || strings.ToLower(t.Abbreviation) == query {
					games = lib.GetGamesForTeam(t)
				}
			}
		}

		if len(games) == 0 {
			fmt.Println("No recent games found")
			return
		}

		for _, g := range games {
			renderGame(g)
		}
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
