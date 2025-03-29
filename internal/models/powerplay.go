package models

type Powerplay struct {
	ID          int       `json:"id" db:"id"`
	HomeTeams   []string  `json:"teams_home_array"`
	AwayTeams   []string  `json:"teams_away_array"`
	Percentages []float32 `json:"opt1_percs_array"`
	Odds        []float32 `json:"opt1_odds_array"`
}
