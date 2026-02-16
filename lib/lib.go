package lib

import (
	_ "embed"
	"encoding/json"
	"log"
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

func GetTeams() []Team {
	var t []Team
	err := json.Unmarshal(teams, &t)
	if err != nil {
		log.Fatal(err)
	}

	return t
}
