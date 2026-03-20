/*
Copyright © 2026 Andy Yoon
*/
package cmd

import (
	"fmt"
	"log"
	"os"
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
	homeScore := g.DisplayScore("home")
	visitorScore := g.DisplayScore("visitor")

	home := teamCellStyle.Render(g.GetHomeTeamName()) + homeScore
	visitor := teamCellStyle.Render(g.GetVisitorTeamName()) + visitorScore

	// Bold the winner if the game has ended.
	if g.CompletionStatus() == lib.Final {
		if g.GetHomeTeamScore() > g.GetVisitorTeamScore() {
			home = boldStyle.Render(home)
		} else {
			visitor = boldStyle.Render(visitor)
		}
	}

	// Display the game date.
	var dateDisplay string
	datetime, err := time.Parse(time.RFC3339, g.GetDatetime())
	if err != nil {
		log.Printf("[warning] Unable to parse date %s, %v", g.GetDatetime(), err)
		dateDisplay = g.GetDatetime()
	} else {
		dateDisplay = datetime.Local().Format("Mon Jan 02")
	}

	// Display the game time.
	timeDisplay := g.DisplayTime()

	// Return strings to be rendered in a single table row
	return []string{visitor + "\n" + home, timeDisplay + "\n" + dateDisplay}
}

func renderGamesTable(rows [][]string) {
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
	Use:     "list [TEAM]",
	Aliases: []string{"ls"},
	Short:   "List scores by league or team",
	Long: `List scores by league or team.

If no team is given, lists today's games for all teams in the league.

If a team is given, lists the upcoming schedule for that team. You can specify
the team name by name, location, or abbreviation (case-insensitive).
For example, all of these commands will resolve to Los Angeles Lakers:

  wts ls lakers
  wts ls "los angeles"
  wts ls LAL

If no league is given, we default to NBA.
Currently supported leagues are: NBA, MLB. (case-insensitive)`,
	Run: func(cmd *cobra.Command, args []string) {
		league, err := cmd.Flags().GetString("league")
		if err != nil {
			log.Fatalf("Parsing error: %v", err)
		}
		history, err := cmd.Flags().GetBool("history")
		if err != nil {
			log.Fatalf("Parsing error: %v", err)
		}

		provider, err := lib.NewProvider(league)
		if err != nil {
			log.Fatalf("Error loading %s: %v", league, err)
		}

		teams, err := provider.Teams()
		if err != nil {
			log.Fatalf("Error loading teams for league %s: %v", league, err)
		}

		var games []lib.Game
		if len(args) == 0 {
			if history {
				games, err = provider.HistoricalGames()
			} else {
				games, err = provider.UpcomingGames()
			}
			if err != nil {
				log.Fatalf("Error loading games: %v", err)
			}
		} else {
			// TODO: Support multiple teams in args. If 2 teams are given and they are playing each other,
			// we could show the expanded game details. Or, no, just use a flag that's better.
			query := strings.ToLower(args[0])
			for _, t := range teams {
				if strings.ToLower(t.GetName()) == query || strings.ToLower(t.GetLocation()) == query || strings.ToLower(t.GetAbbreviation()) == query {
					if history {
						games, err = provider.HistoricalGamesForTeam(t)
					} else {
						games, err = provider.UpcomingGamesForTeam(t)
					}
					if err != nil {
						log.Fatalf("Error loading games: %v", err)
					}
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

	// TODO: Support searching in all leagues if not given. We default somewhat arbitrarily to NBA.
	listCmd.Flags().StringP("league", "l", "nba", "Filter games by league")

	// TODO: Implement these
	listCmd.Flags().BoolP("history", "H", false, "List historical games")
	// listCmd.Flags().StringP("limit", "n", "", "Number of games to list (default all for league, 3 for team)")
	// some sort of date filter. Either number of days or a date range.
}
