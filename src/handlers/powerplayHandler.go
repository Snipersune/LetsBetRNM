package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/snipersune/LetsBetRNM/src/auth"
	"github.com/snipersune/LetsBetRNM/src/renderers"
)

// Powerplay page handler
func (h *AppHandler) PowerplayHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("server: powerplay handler started")
	defer fmt.Println("server: powerplay handler ended")

	renderers.RenderTemplate(w, "html/rendered/powerplay.html", nil)
}

// Handler to process powerplay game submission
func (h *AppHandler) ProcessPowerplayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	err := r.ParseForm() // Parse form data

	if err != nil {
		log.Fatal(err)
	}

	results := make(map[string][]int)

	for i := 1; i <= 8; i++ {
		rowID := fmt.Sprintf("Row-%d", i)
		for _, option := range []string{"1", "X", "2"} {
			key := fmt.Sprintf("r%d-%s", i, option)

			val, err := strconv.Atoi(r.FormValue(key))
			if err != nil {
				log.Fatal(err)
			}

			results[rowID] = append(results[rowID], val) // "1" (selected) or "0" (not selected)
		}
	}

	fmt.Printf("Final User Selections: %v\n", results)

	// Convert map to JSON
	jsonResults, err := json.Marshal(results)
	if err != nil {
		fmt.Println("Failed to parse results to json:", err)
		return
	}

	// Insert user into database
	_, err = h.db.Exec("INSERT INTO bets (powerplay_id, user_id, data) VALUES (?, ?, ?)", 1, userID, jsonResults)
	if err != nil {
		fmt.Println("Failed to insert json results data:", err)
		return
	}

	http.Redirect(w, r, "/powerplay", http.StatusSeeOther)

}
