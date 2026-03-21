package lib

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

//go:embed nba.json
var nba []byte

//go:embed mlb.json
var mlb []byte

// NOTE: balldontlie free tier gives access to Teams, Players, Games,
// at 5 req/min. https://www.balldontlie.io/#pricing

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

// -----------------------------------------------------------------------------
// Teams
// -----------------------------------------------------------------------------

type Team interface {
	GetId() int
	GetName() string
	GetFullName() string
	GetAbbreviation() string
	GetLocation() string
}

type NbaTeam struct {
	Id           int    `json:"id"`
	Conference   string `json:"conference"`
	Division     string `json:"division"`
	City         string `json:"city"`
	Name         string `json:"name"`
	FullName     string `json:"full_name"`
	Abbreviation string `json:"abbreviation"`
}

func (t NbaTeam) GetId() int {
	return t.Id
}
func (t NbaTeam) GetName() string {
	return t.Name
}
func (t NbaTeam) GetFullName() string {
	return t.FullName
}
func (t NbaTeam) GetAbbreviation() string {
	return t.Abbreviation
}
func (t NbaTeam) GetLocation() string {
	return t.City
}

type MlbTeam struct {
	Id               int    `json:"id"`
	Slug             string `json:"slug"`
	Abbreviation     string `json:"abbreviation"`
	DisplayName      string `json:"display_name"`
	ShortDisplayName string `json:"short_display_name"`
	Name             string `json:"name"`
	Location         string `json:"location"`
	League           string `json:"league"`
	Division         string `json:"division"`
}

func (t MlbTeam) GetId() int {
	return t.Id
}
func (t MlbTeam) GetName() string {
	return t.Name
}
func (t MlbTeam) GetFullName() string {
	return t.DisplayName
}
func (t MlbTeam) GetAbbreviation() string {
	return t.Abbreviation
}
func (t MlbTeam) GetLocation() string {
	return t.Location
}

// -----------------------------------------------------------------------------
// Games
// -----------------------------------------------------------------------------

type GameCompletionStatus string

const (
	NotStarted GameCompletionStatus = "Not Started"
	InProgress GameCompletionStatus = "In Progress"
	Final      GameCompletionStatus = "Final"
)

type Game interface {
	GetId() int
	GetDatetime() string
	GetSeason() int
	GetStatus() string
	GetPeriod() int
	GetHomeTeamName() string
	GetVisitorTeamName() string
	GetHomeTeamScore() int
	GetVisitorTeamScore() int
	CompletionStatus() GameCompletionStatus
	DisplayEndStatus() string
	DisplayTime() string
	DisplayScore(t string) string
}

// https://docs.balldontlie.io/#games
type NbaGame struct {
	Id                       int     `json:"id"`
	Date                     string  `json:"date"`
	Season                   int     `json:"season"`
	Status                   string  `json:"status"`
	Period                   int     `json:"period"`
	Time                     string  `json:"time"`
	Postseason               bool    `json:"postseason"`
	Postponed                bool    `json:"postponed"`
	HomeTeamScore            int     `json:"home_team_score"`
	VisitorTeamScore         int     `json:"visitor_team_score"`
	Datetime                 string  `json:"datetime"`
	HomeQ1                   int     `json:"home_q1"`
	HomeQ2                   int     `json:"home_q2"`
	HomeQ3                   int     `json:"home_q3"`
	HomeQ4                   int     `json:"home_q4"`
	HomeOt1                  int     `json:"home_ot1"`
	HomeOt2                  int     `json:"home_ot2"`
	HomeOt3                  int     `json:"home_ot3"`
	HomeTimeoutsRemaining    int     `json:"home_timeouts_remaining"`
	HomeInBonus              bool    `json:"home_in_bonus"`
	VisitorQ1                int     `json:"visitor_q1"`
	VisitorQ2                int     `json:"visitor_q2"`
	VisitorQ3                int     `json:"visitor_q3"`
	VisitorQ4                int     `json:"visitor_q4"`
	VisitorOt1               int     `json:"visitor_ot1"`
	VisitorOt2               int     `json:"visitor_ot2"`
	VisitorOt3               int     `json:"visitor_ot3"`
	VisitorTimeoutsRemaining int     `json:"visitor_timeouts_remaining"`
	VisitorInBonus           bool    `json:"visitor_in_bonus"`
	IstStage                 string  `json:"ist_stage"`
	HomeTeam                 NbaTeam `json:"home_team"`
	VisitorTeam              NbaTeam `json:"visitor_team"`
}

func (g NbaGame) GetId() int {
	return g.Id
}
func (g NbaGame) GetDatetime() string {
	return g.Datetime
}
func (g NbaGame) GetSeason() int {
	return g.Season
}
func (g NbaGame) GetStatus() string {
	return g.Status
}
func (g NbaGame) GetPeriod() int {
	return g.Period
}
func (g NbaGame) GetHomeTeamName() string {
	return g.HomeTeam.GetFullName()
}
func (g NbaGame) GetVisitorTeamName() string {
	return g.VisitorTeam.GetFullName()
}
func (g NbaGame) GetHomeTeamScore() int {
	return g.HomeTeamScore
}
func (g NbaGame) GetVisitorTeamScore() int {
	return g.VisitorTeamScore
}
func (g NbaGame) CompletionStatus() GameCompletionStatus {
	if g.Status == "Final" {
		return Final
	} else if g.Period == 0 {
		return NotStarted
	} else {
		return InProgress
	}
}
func (g NbaGame) DisplayEndStatus() string {
	timeDisplay := "Final"
	if g.GetPeriod() > 4 {
		timeDisplay += "/OT"
	}
	if g.GetPeriod() > 5 {
		timeDisplay += strconv.Itoa(g.GetPeriod() - 4)
	}
	return timeDisplay
}

func (g NbaGame) DisplayTime() string {
	if g.CompletionStatus() == Final {
		return g.DisplayEndStatus()
	} else if g.CompletionStatus() == InProgress {
		return g.Time
	} else {
		datetime, err := time.Parse(time.RFC3339, g.GetDatetime())
		if err != nil {
			log.Printf("[warning] Unable to parse date %s, %v", g.GetDatetime(), err)
			return g.Date
		}
		return datetime.Local().Format("03:04 PM")
	}
}

// Display scores padded to 3-digits
func (g NbaGame) DisplayScore(t string) string {
	if g.CompletionStatus() == NotStarted {
		return "   "
	}
	if t == "home" {
		return fmt.Sprintf("%3d", g.GetHomeTeamScore())
	} else {
		return fmt.Sprintf("%3d", g.GetVisitorTeamScore())
	}
}

type MlbGameTeamData struct {
	Runs         int   `json:"runs"`
	Hits         int   `json:"hits"`
	Errors       int   `json:"errors"`
	InningScores []int `json:"inning_scores"`
}

type MlbGameScoringPlay struct {
	Play      string `json:"play"`
	Inning    string `json:"inning"`
	Period    string `json:"period"`
	AwayScore int    `json:"away_score"`
	HomeScore int    `json:"home_score"`
}

type MlbGame struct {
	Id             int                  `json:"id"`
	HomeTeamName   string               `json:"home_team_name"`
	AwayTeamName   string               `json:"away_team_name"`
	HomeTeam       MlbTeam              `json:"home_team"`
	AwayTeam       MlbTeam              `json:"away_team"`
	Season         int                  `json:"season"`
	Postseason     bool                 `json:"postseason"`
	SeasonType     string               `json:"season_type"`
	Date           string               `json:"date"`
	HomeTeamData   MlbGameTeamData      `json:"home_team_data"`
	AwayTeamData   MlbGameTeamData      `json:"away_team_data"`
	Venue          string               `json:"venue"`
	Attendance     int                  `json:"attendance"`
	ConferencePlay bool                 `json:"conference_play"`
	Status         string               `json:"status"`
	Period         int                  `json:"period"`
	Clock          int                  `json:"clock"`
	DisplayClock   string               `json:"display_clock"`
	ScoringSummary []MlbGameScoringPlay `json:"scoring_summary"`
}

func (g MlbGame) GetId() int {
	return g.Id
}
func (g MlbGame) GetDatetime() string {
	return g.Date
}
func (g MlbGame) GetSeason() int {
	return g.Season
}
func (g MlbGame) GetStatus() string {
	return g.Status
}
func (g MlbGame) GetPeriod() int {
	return g.Period
}
func (g MlbGame) GetHomeTeamName() string {
	return g.HomeTeam.DisplayName
}
func (g MlbGame) GetVisitorTeamName() string {
	return g.AwayTeam.DisplayName
}
func (g MlbGame) GetHomeTeamScore() int {
	return g.HomeTeamData.Runs
}
func (g MlbGame) GetVisitorTeamScore() int {
	return g.AwayTeamData.Runs
}
func (g MlbGame) CompletionStatus() GameCompletionStatus {
	if g.Status == "STATUS_FINAL" {
		return Final
	} else if g.Status == "STATUS_SCHEDULED" {
		return NotStarted
	} else {
		return InProgress
	}
}
func (g MlbGame) DisplayEndStatus() string {
	if g.Period > 9 {
		return fmt.Sprintf("Final/%d", g.Period)
	}
	return "Final"
}

func (g MlbGame) DisplayTime() string {
	if g.CompletionStatus() == Final {
		return g.DisplayEndStatus()
	} else if g.CompletionStatus() == InProgress {
		// Get current inning based on number of inning scores
		homeScores := len(g.HomeTeamData.InningScores)
		awayScores := len(g.AwayTeamData.InningScores)
		if awayScores > homeScores {
			return fmt.Sprintf("Top %d", g.Period)
		} else {
			return fmt.Sprintf("Bot %d", g.Period)
		}
	} else {
		datetime, err := time.Parse(time.RFC3339, g.GetDatetime())
		if err != nil {
			log.Printf("[warning] Unable to parse date %s, %v", g.GetDatetime(), err)
			return g.Date
		}
		return datetime.Local().Format("03:04 PM")
	}
}

// Display scores padded to 2-digits
func (g MlbGame) DisplayScore(t string) string {
	if g.CompletionStatus() == NotStarted {
		return "  "
	}
	if t == "home" {
		return fmt.Sprintf("%2d", g.GetHomeTeamScore())
	} else {
		return fmt.Sprintf("%2d", g.GetVisitorTeamScore())
	}
}

// -----------------------------------------------------------------------------
// Providers
// -----------------------------------------------------------------------------

type LeagueProvider interface {
	Teams() ([]Team, error)
	CurrentSeason() int
	UpcomingGames() ([]Game, error)
	UpcomingGamesForTeam(team Team) ([]Game, error)
	HistoricalGames() ([]Game, error)
	HistoricalGamesForTeam(team Team) ([]Game, error)
}

type NbaProvider struct{}

func (p NbaProvider) Teams() ([]Team, error) {
	var nbaTeams []NbaTeam
	if err := json.Unmarshal(nba, &nbaTeams); err != nil {
		return nil, err
	}
	teams := make([]Team, len(nbaTeams))
	for i, t := range nbaTeams {
		teams[i] = t
	}
	return teams, nil
}

// Handle season param: Season typically starts in Oct and ends in June
func (p NbaProvider) CurrentSeason() int {
	today := time.Now()
	todayYear, todayMonth, _ := today.Date()
	if todayMonth < time.August {
		todayYear -= 1
	}
	return todayYear
}

func (p NbaProvider) UpcomingGames() ([]Game, error) {
	path := fmt.Sprintf("/v1/games?dates[]=%s", getTodayDate())
	gs, err := fetchNbaGames(path)
	if err != nil {
		return nil, err
	}
	games := make([]Game, len(gs))
	for i, g := range gs {
		games[i] = g
	}
	return games, nil
}

func (p NbaProvider) UpcomingGamesForTeam(team Team) ([]Game, error) {
	path := fmt.Sprintf("/v1/games?team_ids[]=%d&start_date=%s&per_page=3", team.GetId(), getTodayDate())
	gs, err := fetchNbaGames(path)
	if err != nil {
		return nil, err
	}
	games := make([]Game, len(gs))
	for i, g := range gs {
		games[i] = g
	}
	return games, nil
}

func (p NbaProvider) HistoricalGames() ([]Game, error) {
	// BALLDONTLIE doesn't let us sort by date desc. So we just get the last 2 days as a rough estimate.
	startDate := getTodayPlusOffsetDate(-2)
	endDate := getTodayDate()
	season := p.CurrentSeason()

	path := fmt.Sprintf("/v1/games?start_date=%s&end_date=%s&seasons[]=%d", startDate, endDate, season)
	gs, err := fetchNbaGames(path)
	if err != nil {
		return nil, err
	}
	games := make([]Game, len(gs))
	for i, g := range gs {
		games[i] = g
	}
	return games, nil
}

func (p NbaProvider) HistoricalGamesForTeam(t Team) ([]Game, error) {
	// BALLDONTLIE doesn't let us sort by date desc. So we just get the last 7 days as a rough estimate.
	startDate := getTodayPlusOffsetDate(-7)
	endDate := getTodayDate()
	season := p.CurrentSeason()

	path := fmt.Sprintf("/v1/games?team_ids[]=%d&start_date=%s&end_date=%s&seasons[]=%d", t.GetId(), startDate, endDate, season)
	gs, err := fetchNbaGames(path)
	if err != nil {
		return nil, err
	}
	games := make([]Game, len(gs))
	for i, g := range gs {
		games[i] = g
	}
	return games, nil
}

func fetchNbaGames(path string) ([]NbaGame, error) {
	data, err := Get(path)
	if err != nil {
		return nil, err
	}

	var g ListResponse[NbaGame]
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, err
	}

	return g.Data, nil
}

type MlbProvider struct{}

func (p MlbProvider) Teams() ([]Team, error) {
	var mlbTeams []MlbTeam
	if err := json.Unmarshal(mlb, &mlbTeams); err != nil {
		return nil, err
	}
	teams := make([]Team, len(mlbTeams))
	for i, t := range mlbTeams {
		teams[i] = t
	}
	return teams, nil
}

func (p MlbProvider) CurrentSeason() int {
	today := time.Now()
	todayYear, _, _ := today.Date()
	return todayYear
}

func (p MlbProvider) UpcomingGames() ([]Game, error) {
	path := fmt.Sprintf("/mlb/v1/games?dates[]=%s", getTodayDate())
	gs, err := fetchMlbGames(path)
	if err != nil {
		return nil, err
	}
	games := make([]Game, len(gs))
	for i, g := range gs {
		games[i] = g
	}
	return games, nil
}

func (p MlbProvider) UpcomingGamesForTeam(team Team) ([]Game, error) {
	nDays := 3
	days := make([]string, 0, nDays)
	days = append(days, getTodayDate())
	for i := 1; i < nDays; i++ {
		days = append(days, getTodayPlusOffsetDate(i))
	}
	daysParam := "dates[]=" + strings.Join(days, "&dates[]=")

	path := fmt.Sprintf("/mlb/v1/games?team_ids[]=%d&%s", team.GetId(), daysParam)
	fmt.Println("Requesting path: " + path)
	gs, err := fetchMlbGames(path)
	if err != nil {
		return nil, err
	}
	games := make([]Game, len(gs))
	for i, g := range gs {
		games[i] = g
	}
	return games, nil
}

func (p MlbProvider) HistoricalGames() ([]Game, error) {
	// BALLDONTLIE doesn't let us sort by date desc. So we just get the last 2 days as a rough estimate.
	dateQueryParam := makeMlbDateRangeQueryParam(-2, 2)
	season := p.CurrentSeason()

	path := fmt.Sprintf("/mlb/v1/games?%s&seasons[]=%d", dateQueryParam, season)
	gs, err := fetchMlbGames(path)
	if err != nil {
		return nil, err
	}
	games := make([]Game, len(gs))
	for i, g := range gs {
		games[i] = g
	}
	return games, nil
}

func (p MlbProvider) HistoricalGamesForTeam(t Team) ([]Game, error) {
	// BALLDONTLIE doesn't let us sort by date desc. So we just get the last 7 days as a rough estimate.
	dateQueryParam := makeMlbDateRangeQueryParam(-7, 7)
	season := p.CurrentSeason()

	path := fmt.Sprintf("/mlb/v1/games?team_ids[]=%d&%s&seasons[]=%d", t.GetId(), dateQueryParam, season)
	gs, err := fetchMlbGames(path)
	if err != nil {
		return nil, err
	}
	games := make([]Game, len(gs))
	for i, g := range gs {
		games[i] = g
	}
	return games, nil
}

func fetchMlbGames(path string) ([]MlbGame, error) {
	data, err := Get(path)
	if err != nil {
		return nil, err
	}

	var g ListResponse[MlbGame]
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, err
	}

	return g.Data, nil
}

func NewProvider(league string) (LeagueProvider, error) {
	switch strings.ToLower(league) {
	case "nba":
		return NbaProvider{}, nil
	case "mlb":
		return MlbProvider{}, nil
	default:
		return nil, fmt.Errorf("Sorry, %s is not supported yet.", league)
	}
}

// -----------------------------------------------------------------------------
// Functions
// -----------------------------------------------------------------------------

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

// To query MLB date ranges, we have to pass each date one at a time.
func makeMlbDateRangeQueryParam(startOffsetFromToday int, numDays int) string {
	days := make([]string, 0, numDays)
	days = append(days, getTodayPlusOffsetDate(startOffsetFromToday))
	for i := 1; i < numDays; i++ {
		days = append(days, getTodayPlusOffsetDate(startOffsetFromToday+i))
	}
	return "dates[]=" + strings.Join(days, "&dates[]=")
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

// func GetGames() []Game {
// 	startDate, endDate, season := buildDateRanges(1)

// 	path := fmt.Sprintf("/games?%s&%s&%s", startDate, endDate, season)
// 	return fetchGames(path)
// }

// func GetGamesForTeam(team NbaTeam) []Game {
// 	teamIds := fmt.Sprintf("team_ids[]=%d", team.Id)
// 	startDate, endDate, season := buildDateRanges(7)

// 	path := fmt.Sprintf("/games?%s&%s&%s&%s", teamIds, startDate, endDate, season)
// 	return fetchGames(path)
// }
