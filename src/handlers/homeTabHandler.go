package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/snipersune/LetsBetRNM/src/auth"
	"github.com/snipersune/LetsBetRNM/src/renderers"
)

func (h *AppHandler) HomeTabHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	fmt.Printf("Request from user ID: %d\n", userID)

	queryGameRows, err := h.db.Query(
		`SELECT 
		bets.powerplay_id,
		COALESCE(teams.name, users.username) AS participant,
		bets.placed_at, 
		bets.data
		FROM bets
		LEFT JOIN teams ON bets.team_id = teams.id
		LEFT JOIN users ON bets.user_id = users.id
		WHERE bets.user_id = ? OR bets.team_id IN (
    	SELECT team_id FROM team_members WHERE user_id = ?
		)
		ORDER BY bets.placed_at DESC;
		`, userID, userID,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer queryGameRows.Close()

	type BasicBetInfo struct {
		Participant string
		PlacedAt    string
	}

	var basicBetInfoArr []BasicBetInfo

	for queryGameRows.Next() {
		var gameID int
		var basicBetInfo BasicBetInfo
		var jsonBetData string

		err := queryGameRows.Scan(&gameID, &basicBetInfo.Participant, &basicBetInfo.PlacedAt, &jsonBetData)
		if err != nil {
			log.Fatal(err)
		}

		// Parse and format
		t, _ := time.Parse(time.RFC3339, basicBetInfo.PlacedAt)
		basicBetInfo.PlacedAt = t.Format("2006-01-02 15:04:05")

		basicBetInfoArr = append(basicBetInfoArr, basicBetInfo)
	}

	type GameInfo struct {
		HomeTeams       []string `json:"teams_home_array"`
		AwayTeams       []string `json:"teams_away_array"`
		PercentagesOpt1 []string `json:"opt1_percs_array"`
		PercentagesOptX []string `json:"optX_percs_array"`
		PercentagesOpt2 []string `json:"opt2_percs_array"`
		OddsOpt1        []string `json:"opt1_odds_array"`
		OddsOptX        []string `json:"optX_odds_array"`
		OddsOpt2        []string `json:"opt2_odds_array"`
	}

	gameJsonFile, err := os.Open("output/powerplay_data.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer gameJsonFile.Close()

	decoder := json.NewDecoder(gameJsonFile)

	// Decode the JSON into struct
	var gameInfo GameInfo
	err = decoder.Decode(&gameInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	var jsonGameInfo []byte
	jsonGameInfo, err = json.Marshal(gameInfo)
	if err != nil {
		fmt.Println("Error encoding back to JSON:", err)
		return
	}

	type FileData struct {
		GameInfo   string
		BetHistory []BasicBetInfo
	}

	var fileData = FileData{
		string(jsonGameInfo),
		basicBetInfoArr,
	}

	// gameJsonData, err := os.ReadFile("output/powerplay_data.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// type FileData struct {
	// 	GameInfo   string
	// 	BetHistory []BasicBetInfo
	// }
	//
	// var fileData = FileData{
	// 	string(gameJsonData),
	// 	basicBetInfoArr,
	// }

	renderers.RenderTemplate(w, "html/static/test_tab_powerplay.html", fileData)
}
