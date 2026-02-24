/*
Copyright © 2026 Andy Yoon
*/
package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/andyyoon2/whats-the-score/lib"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

var (
	tableRowStyle     = lipgloss.NewStyle().Padding(0, 1)
	tableHeadingStyle = tableRowStyle.Foreground(lipgloss.Color("250"))
	teamColStyle      = tableRowStyle.Width(24) // Longest team name (22) plus padding
	dateCellStyle     = tableHeadingStyle.Width(24)
	boldStyle         = lipgloss.NewStyle().Bold(true)
)

func renderGame(g lib.Game) {
	homeScore := fmt.Sprintf("%3d", g.HomeTeamScore)
	visitorScore := fmt.Sprintf("%3d", g.VisitorTeamScore)
	if g.Status == "Final" {
		if g.HomeTeamScore > g.VisitorTeamScore {
			homeScore = boldStyle.Render(homeScore)
		} else {
			visitorScore = boldStyle.Render(visitorScore)
		}
	}

	var timeDisplay string
	if g.Period == 0 {
		// Game hasn't started yet, show the start time
		datetime, err := time.Parse(time.RFC3339, g.Datetime)
		if err != nil {
			log.Printf("[warning] Unable to parse date %s, %v", g.Datetime, err)
			timeDisplay = g.Date
		}

		timeDisplay = datetime.Local().Format("03:04 PM MST")
	} else {
		timeDisplay = g.Time
	}

	rows := [][]string{
		{g.Date, timeDisplay},
		{g.VisitorTeam.FullName, visitorScore},
		{g.HomeTeam.FullName, homeScore},
	}
	t := table.New().StyleFunc(func(row, col int) lipgloss.Style {
		if row == 0 {
			if col == 0 {
				return dateCellStyle
			} else {
				return tableHeadingStyle
			}
		}
		if col == 0 {
			return teamColStyle
		}
		return tableRowStyle
	}).Rows(rows...)

	fmt.Println(t)
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "A brief description of your command",
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
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
