/*
Copyright © 2026 Andy Yoon
*/
package cmd

import (
	"fmt"
	"log/slog"
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
		slog.Warn(fmt.Sprintf("Unable to parse date %s, %v", g.GetDatetime(), err))
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

List today's games for all teams in the league:
  wts ls

List the upcoming schedule for a team:
  wts ls lakers
  wts ls --league mlb dodgers

You can specify the team by name, location, or abbreviation:
  wts ls spurs  --> San Antonio Spurs
  wts ls denver --> Denver Nuggets
  wts ls bos    --> Boston Celtics

If no league is given, we default to NBA. Currently supported leagues: NBA, MLB.
All args are case-insensitive.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		league, err := cmd.Flags().GetString("league")
		if err != nil {
			slog.Error(fmt.Sprintf("Parsing error: %v", err))
			return err
		}
		history, err := cmd.Flags().GetBool("history")
		if err != nil {
			slog.Error(fmt.Sprintf("Parsing error: %v", err))
			return err
		}

		provider, err := lib.NewProvider(league)
		if err != nil {
			slog.Error(fmt.Sprintf("Error loading %s: %v", league, err))
			return err
		}

		teams, err := provider.Teams()
		if err != nil {
			slog.Error(fmt.Sprintf("Error loading teams for league %s: %v", league, err))
			return err
		}

		var games []lib.Game
		if len(args) == 0 {
			if history {
				games, err = provider.HistoricalGames()
			} else {
				games, err = provider.UpcomingGames()
			}
			if err != nil {
				slog.Error(fmt.Sprintf("Error loading games: %v", err))
				return err
			}
		} else {
			query := strings.ToLower(args[0])
			for _, t := range teams {
				if strings.ToLower(t.GetName()) == query || strings.ToLower(t.GetLocation()) == query || strings.ToLower(t.GetAbbreviation()) == query {
					if history {
						games, err = provider.HistoricalGamesForTeam(t)
					} else {
						games, err = provider.UpcomingGamesForTeam(t)
					}
					if err != nil {
						slog.Error(fmt.Sprintf("Error loading games: %v", err))
						return err
					}
				}
			}
		}

		if len(games) == 0 {
			fmt.Println("No recent games found")
			return nil
		}

		var rows [][]string
		for _, g := range games {
			rows = append(rows, renderGame(g))
		}

		renderGamesTable(rows)
		return nil
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	// TODO: Support searching by team in all leagues if not given.
	listCmd.Flags().StringP("league", "l", "nba", "Filter by league")
	listCmd.Flags().BoolP("history", "H", false, "List historical games")
}
