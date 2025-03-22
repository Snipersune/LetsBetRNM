package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/snipersune/LetsBetRNM/src/auth"
	"github.com/snipersune/LetsBetRNM/src/renderers"
)

// History handler
func (h *AppHandler) HistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	fmt.Printf("%d\n", userID)

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

	// type GameDetails struct {
	// 	GameData string `json:"GAMEDATA"` // Not finished
	// 	Bets     string `json:"BETS"`
	// }
	//
	// type GameData struct {
	// 	BasicInfo BasicBetInfo
	// 	Details   GameDetails
	// }

	var basicBetInfoArr []BasicBetInfo
	//var gameDetailsArr []GameDetails
	//var gameDataArr []GameData

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

		/*
			var jsonGameDetails string
			err = db.QueryRow("SELECT data FROM powerplays WHERE id = ?", gameID).Scan(&jsonGameDetails)
			if err != nil {
				log.Fatal(err)
				return
			}

			var gameData GameData
			gameData.BasicInfo = basicBetInfo

			// Convert JSON string to struct
			err = json.Unmarshal([]byte(jsonGameDetails), &gameData.Details)
			if err != nil {
				log.Fatal(err)
			}

			gameDetailsArr = append(gameDetailsArr, gameData.Details)
			gameDataArr = append(gameDataArr, gameData)
		*/

		basicBetInfoArr = append(basicBetInfoArr, basicBetInfo)
	}

	renderers.RenderTemplate(w, "html/static/history.html", basicBetInfoArr)
}
