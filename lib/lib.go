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

func GetGamesForTeam(team Team) []Game {
	// NOTE: When formatting strings, you need to describe the reference date
	// https://pkg.go.dev/time#example-Time.Format
	//	Jan 2 15:04:05 2006 MST
	// An easy way to remember this value is that it holds, when presented
	// in this order, the values (lined up with the elements above):
	//	  1 2  3  4  5    6  -7
	teamIds := fmt.Sprintf("team_ids[]=%d", team.Id)

	// TODO: This doesn't quite work always to get the latest games, but it's good enough for now
	today := time.Now()
	lastWeek := today.AddDate(0, 0, -7)
	startDate := fmt.Sprintf("start_date=%s", lastWeek.Format("2006-01-02"))
	endDate := fmt.Sprintf("end_date=%s", today.Format("2006-01-02"))

	// Handle season param: Season typically starts in Oct and ends in June
	todayYear, todayMonth, _ := today.Date()
	if todayMonth < time.August {
		todayYear -= 1
	}
	season := fmt.Sprintf("seasons[]=%d", todayYear)

	path := fmt.Sprintf("/games?%s&%s&%s&%s", teamIds, season, startDate, endDate)
	data := Get(path)

	var g ListResponse[Game]
	err := json.Unmarshal(data, &g)
	if err != nil {
		log.Fatal(err)
	}

	return g.Data
}
