package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

var db *sql.DB

type contextKey string

const userIDKey contextKey = "user_id"

// Home screen handler
func historyHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	queryGameRows, err := db.Query(
		`SELECT 
		bets.game_id,
		bets.placed_at, 
		COALESCE(teams.name, users.username) AS participant,
		bets.data,
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

	type BasicGameInfo struct {
		PlacedAt    string
		Participant string
	}

	type GameDetails struct {
		GameData string `json:"GAMEDATA"` // Not finished
		Bets     string `json:"BETS"`
	}

	type GameData struct {
		BasicInfo BasicGameInfo
		Details   GameDetails
	}

	var basicGameInfoArr []BasicGameInfo
	//var gameDetailsArr []GameDetails
	//var gameDataArr []GameData

	for queryGameRows.Next() {
		var gameID int
		var basicGameInfo BasicGameInfo
		var jsonBetData string

		err := queryGameRows.Scan(&gameID, &basicGameInfo, &jsonBetData)
		if err != nil {
			log.Fatal(err)
		}

		/*
			var jsonGameDetails string
			err = db.QueryRow("SELECT data FROM powerplays WHERE id = ?", gameID).Scan(&jsonGameDetails)
			if err != nil {
				log.Fatal(err)
				return
			}

			var gameData GameData
			gameData.BasicInfo = basicGameInfo

			// Convert JSON string to struct
			err = json.Unmarshal([]byte(jsonGameDetails), &gameData.Details)
			if err != nil {
				log.Fatal(err)
			}

			gameDetailsArr = append(gameDetailsArr, gameData.Details)
			gameDataArr = append(gameDataArr, gameData)
		*/

		basicGameInfoArr = append(basicGameInfoArr, basicGameInfo)
	}

	renderTemplate(w, "html/static/history.html", basicGameInfoArr)
}
