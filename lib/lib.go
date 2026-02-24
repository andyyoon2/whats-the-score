package lib

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

//go:embed teams.json
var teams []byte

// NOTE: balldontlie free tier gives access to Teams, Players, Games,
// at 5 req/min. https://www.balldontlie.io/#pricing

type Team struct {
	Id           int    `json:"id"`
	Conference   string `json:"conference"`
	Division     string `json:"division"`
	City         string `json:"city"`
	Name         string `json:"name"`
	FullName     string `json:"full_name"`
	Abbreviation string `json:"abbreviation"`
}

// https://docs.balldontlie.io/#pagination
type ListMetadata struct {
	NextCursor int `json:"next_cursor"`
	PerPage    int `json:"per_page"`
}

// https://docs.balldontlie.io/#get-all-games
type ListResponse[T any] struct {
	Data     []T          `json:"data"`
	Metadata ListMetadata `json:"meta"`
}

// https://docs.balldontlie.io/#games
type Game struct {
	Id                       int    `json:"id"`
	Date                     string `json:"date"`
	Season                   int    `json:"season"`
	Status                   string `json:"status"`
	Period                   int    `json:"period"`
	Time                     string `json:"time"`
	Postseason               bool   `json:"postseason"`
	Postponed                bool   `json:"postponed"`
	HomeTeamScore            int    `json:"home_team_score"`
	VisitorTeamScore         int    `json:"visitor_team_score"`
	Datetime                 string `json:"datetime"`
	HomeQ1                   int    `json:"home_q1"`
	HomeQ2                   int    `json:"home_q2"`
	HomeQ3                   int    `json:"home_q3"`
	HomeQ4                   int    `json:"home_q4"`
	HomeOt1                  int    `json:"home_ot1"`
	HomeOt2                  int    `json:"home_ot2"`
	HomeOt3                  int    `json:"home_ot3"`
	HomeTimeoutsRemaining    int    `json:"home_timeouts_remaining"`
	HomeInBonus              bool   `json:"home_in_bonus"`
	VisitorQ1                int    `json:"visitor_q1"`
	VisitorQ2                int    `json:"visitor_q2"`
	VisitorQ3                int    `json:"visitor_q3"`
	VisitorQ4                int    `json:"visitor_q4"`
	VisitorOt1               int    `json:"visitor_ot1"`
	VisitorOt2               int    `json:"visitor_ot2"`
	VisitorOt3               int    `json:"visitor_ot3"`
	VisitorTimeoutsRemaining int    `json:"visitor_timeouts_remaining"`
	VisitorInBonus           bool   `json:"visitor_in_bonus"`
	IstStage                 string `json:"ist_stage"`
	HomeTeam                 Team   `json:"home_team"`
	VisitorTeam              Team   `json:"visitor_team"`
}

func GetTeams() []Team {
	var t []Team
	err := json.Unmarshal(teams, &t)
	if err != nil {
		log.Fatal(err)
	}

	return t
}

// Build query parameters for API request
// TODO: This doesn't quite work always to get the latest games, but it's good enough for now
func buildDateRanges(lookback int) (string, string, string) {
	// NOTE: When formatting strings, you need to describe the reference date
	// https://pkg.go.dev/time#example-Time.Format
	//	Jan 2 15:04:05 2006 MST
	// An easy way to remember this value is that it holds, when presented
	// in this order, the values (lined up with the elements above):
	//	  1 2  3  4  5    6  -7
	today := time.Now()
	lastWeek := today.AddDate(0, 0, -lookback)
	startDate := fmt.Sprintf("start_date=%s", lastWeek.Format("2006-01-02"))
	endDate := fmt.Sprintf("end_date=%s", today.Format("2006-01-02"))

	// Handle season param: Season typically starts in Oct and ends in June
	todayYear, todayMonth, _ := today.Date()
	if todayMonth < time.August {
		todayYear -= 1
	}
	season := fmt.Sprintf("seasons[]=%d", todayYear)

	return startDate, endDate, season
}

// Return today in YYYY-MM-DD format
func getTodayDate() string {
	return time.Now().Format("2006-01-02")
}

// Returns today+offsetDays in YYYY-MM-DD format
func getTodayPlusOffsetDate(offsetDays int) string {
	today := time.Now()
	offsetDate := today.AddDate(0, 0, offsetDays)

	return offsetDate.Format("2006-01-02")
}

// Handle season param: Season typically starts in Oct and ends in June
func getSeason() int {
	today := time.Now()
	todayYear, todayMonth, _ := today.Date()
	if todayMonth < time.August {
		todayYear -= 1
	}
	return todayYear
}

func fetchGames(path string) []Game {
	data := Get(path)

	var g ListResponse[Game]
	err := json.Unmarshal(data, &g)
	if err != nil {
		log.Fatal(err)
	}

	return g.Data
}

// Fetch today's games for all teams
func GetUpcomingGames() []Game {
	path := fmt.Sprintf("/games?dates[]=%s", getTodayDate())
	return fetchGames(path)
}

func GetGames() []Game {
	startDate, endDate, season := buildDateRanges(1)

	path := fmt.Sprintf("/games?%s&%s&%s", startDate, endDate, season)
	return fetchGames(path)
}

func GetGamesForTeam(team Team) []Game {
	teamIds := fmt.Sprintf("team_ids[]=%d", team.Id)
	startDate, endDate, season := buildDateRanges(7)

	path := fmt.Sprintf("/games?%s&%s&%s&%s", teamIds, startDate, endDate, season)
	return fetchGames(path)
}
