/*
Copyright © 2026 Andy Yoon
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/andyyoon2/whats-the-score/lib"
	"github.com/spf13/cobra"
)

var (
	hasDarkBG           = lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark           = lipgloss.LightDark(hasDarkBG)
	pxStyle             = lipgloss.NewStyle().Padding(0, 1)
	rightTableCellStyle = pxStyle.Align(lipgloss.Center)
	secondaryTextStyle  = lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("243"), lipgloss.Color("250")))
	tableHeadingStyle   = pxStyle.Foreground(lightDark(lipgloss.Color("243"), lipgloss.Color("250")))
	rightColumnStyle    = tableHeadingStyle.Align(lipgloss.Center)
	teamCellStyle       = lipgloss.NewStyle().Width(24) // Longest team name (22) plus padding
	dateCellStyle       = tableHeadingStyle.Width(24)
	boldStyle           = lipgloss.NewStyle().Bold(true)
)

func renderGame(g lib.Game) []string {
	var homeScore string
	var visitorScore string
	// Hide scores if game hasn't started yet. Align it with 3-digit scores
	if g.Period == 0 {
		homeScore = "   "
		visitorScore = "   "
	} else {
		homeScore = fmt.Sprintf("%3d", g.HomeTeamScore)
		visitorScore = fmt.Sprintf("%3d", g.VisitorTeamScore)
	}

	home := teamCellStyle.Render(g.HomeTeam.FullName) + homeScore
	visitor := teamCellStyle.Render(g.VisitorTeam.FullName) + visitorScore

	// Bold the winner if the game has ended.
	if g.Status == "Final" {
		if g.HomeTeamScore > g.VisitorTeamScore {
			home = boldStyle.Render(home)
		} else {
			visitor = boldStyle.Render(visitor)
		}
	}

	// Display the game date.
	var dateDisplay string
	datetime, err := time.Parse(time.RFC3339, g.Datetime)
	if err != nil {
		log.Printf("[warning] Unable to parse date %s, %v", g.Datetime, err)
		dateDisplay = g.Date
	} else {
		dateDisplay = datetime.Local().Format("Mon Jan 02")
	}

	// Display the game time.
	var timeDisplay string
	if g.Status == "Final" {
		// Game is ended, check for OTs
		timeDisplay = "Final"
		if g.Period > 4 {
			timeDisplay += "/OT"
		}
		if g.Period > 5 {
			timeDisplay += strconv.Itoa(g.Period - 4)
		}
	} else if g.Period == 0 {
		// Game hasn't started
		timeDisplay = datetime.Local().Format("03:04 PM")
	} else {
		// Game in progress
		timeDisplay = g.Time
	}

	// Return strings to be rendered in a single table row
	return []string{visitor + "\n" + home, timeDisplay + "\n" + dateDisplay}
}

func renderGamesTable(rows [][]string) {
	// rows := [][]string{
	// 	{visitor, timeDisplay},
	// 	{home, dateDisplay},
	// }
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderColumn(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				if col == 0 {
					return pxStyle
				} else {
					return rightColumnStyle
				}
			}
			if col == 0 {
				return pxStyle.BorderStyle(lipgloss.NormalBorder()).BorderTop(true)
			} else {
				return rightColumnStyle.BorderStyle(lipgloss.NormalBorder()).BorderTop(true)
			}
		}).
		Rows(rows...)

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
			// games = lib.GetGames()
			games = lib.GetUpcomingGames()
		} else {
			query := strings.ToLower(args[0])
			for _, t := range teams {
				if strings.ToLower(t.Name) == query || strings.ToLower(t.City) == query || strings.ToLower(t.Abbreviation) == query {
					// games = lib.GetGamesForTeam(t)
					games = lib.GetUpcomingGamesForTeam(t)
				}
			}
		}

		if len(games) == 0 {
			fmt.Println("No recent games found")
			return
		}

		var rows [][]string
		for _, g := range games {
			rows = append(rows, renderGame(g))
		}

		renderGamesTable(rows)
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
