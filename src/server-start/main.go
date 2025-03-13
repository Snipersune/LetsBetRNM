package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/gorilla/sessions"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

/*
----------- TODOS -----------

Authentication - User-password (or fake)
Sqlite - Save user data
Betting math
Add html sites to view betting history and stats

*/

type contextKey string

const userIDKey contextKey = "user_id"
const userNameKey contextKey = "username"

// Create a session store (store session data in memory)
var store = sessions.NewCookieStore([]byte("super-secret-key"))

// Database instance
var db *sql.DB

// Middleware to protect routes
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")

		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			// http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, ok := session.Values["user_id"].(int)
		if !ok || userID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userName := session.Values["username"]

		r = r.WithContext(context.WithValue(r.Context(), userIDKey, userID))
		r = r.WithContext(context.WithValue(r.Context(), userNameKey, userName))

		next.ServeHTTP(w, r)
	})
}

// Render templates
func renderTemplate(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// Default screen handler
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// Home screen handler
func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "html/static/home.html", nil)
}

// Home screen handler
func historyHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	fmt.Printf("%d\n", userID)

	queryGameRows, err := db.Query(
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

	type BasicGameInfo struct {
		Participant string
		PlacedAt    string
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

		err := queryGameRows.Scan(&gameID, &basicGameInfo.Participant, &basicGameInfo.PlacedAt, &jsonBetData)
		if err != nil {
			log.Fatal(err)
		}

		// Parse and format
		t, _ := time.Parse(time.RFC3339, basicGameInfo.PlacedAt)
		basicGameInfo.PlacedAt = t.Format("2006-01-02 15:04:05")

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

// Dashboard handler
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "html/static/dashboard.html", nil)
}

// Login page handler
func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage := r.URL.Query().Get("error")
	renderTemplate(w, "html/static/login.html", map[string]string{"ErrorMessage": errorMessage})
}

// Register page handler
func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("server: registerPageHandler requested")
	fmt.Printf("server: request type: %s\n", r.Method)
	errorMessage := r.URL.Query().Get("error")
	renderTemplate(w, "html/static/register.html", map[string]string{"ErrorMessage": errorMessage})
}

// Handler to register a new user
func registerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Redirect(w, r, "/register?error=Could%20not%20process%20password", http.StatusSeeOther)
		return
	}

	// Insert user into database
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, string(hashedPassword))
	if err != nil {
		http.Redirect(w, r, "/register?error=Username%20already%20exists", http.StatusSeeOther)
		return
	}

	fmt.Println("Registration completed. Redirecting to login.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Handler to log in an existing user
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Retrieve user from database
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		http.Redirect(w, r, "/login?error=Invalid%20username%20or%20password", http.StatusSeeOther)
		return
	}

	// Compare hashed passwords
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		http.Redirect(w, r, "/login?error=Invalid%20username%20or%20password", http.StatusSeeOther)
		return
	}

	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		http.Redirect(w, r, "/login?error=Invalid%20username%20or%20password", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Values["user_id"] = userID
	session.Values["username"] = username
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func powerplayPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("server: powerplay handler started")
	defer fmt.Println("server: powerplay handler ended")

	c := exec.Command("/home/dnee/Documents/Kod/LetsBetRNM/.venv/bin/python", "/home/dnee/Documents/Kod/LetsBetRNM/src/update_powerplay_html.py")
	if err := c.Run(); err != nil {
		fmt.Println("server: Error: ", err)
	}
	// renderTemplate(w, "html/rendered/powerplay.html")
	http.ServeFile(w, r, "html/rendered/powerplay.html")
}

func powerplayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(userIDKey).(int)
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
	_, err = db.Exec("INSERT INTO bets (powerplay_id, user_id, data) VALUES (?, ?, ?)", 1, userID, jsonResults)
	if err != nil {
		fmt.Println("Failed to insert json results data:", err)
		return
	}

	http.Redirect(w, r, "/powerplay", http.StatusSeeOther)

}

func main() {

	var err error

	// Initialize SQLite database
	db, err = sql.Open("sqlite3", "./biggestassdatabase.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create users table if it doesn't exist
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`)

	if err != nil {
		log.Fatal(err)
	}

	// Create teams table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS teams (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
    	name TEXT UNIQUE NOT NULL,
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`)

	if err != nil {
		log.Fatal(err)
	}

	// Create team members table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS team_members (
		user_id INTEGER,
		team_id INTEGER,
		role TEXT DEFAULT 'member',
		joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, team_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
	)
	`)

	if err != nil {
		log.Fatal(err)
	}

	// Create powerplay table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS powerplays (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		data TEXT NOT NULL
	)
	`)

	if err != nil {
		log.Fatal(err)
	}

	// Create powerplay table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		powerplay_id INTEGER NOT NULL,
    	team_id INTEGER NULL,	-- NULL if it's an individual bet
    	user_id INTEGER NULL,	-- NULL if it's a team-wide bet
		data TEXT NOT NULL,		-- stores bet placement data
		placed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
		FOREIGN KEY (powerplay_id) REFERENCES powerplays(id) ON DELETE CASCADE
	)
	`)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", defaultHandler)
	http.Handle("/home", authMiddleware(http.HandlerFunc(homeHandler)))

	http.HandleFunc("/login", loginPageHandler)
	http.HandleFunc("/register", registerPageHandler)

	http.HandleFunc("/process-login", loginHandler)
	http.HandleFunc("/process-register", registerHandler)

	http.Handle("/powerplay", authMiddleware(http.HandlerFunc(powerplayPageHandler)))
	http.Handle("/process-powerplay", authMiddleware(http.HandlerFunc(powerplayHandler)))

	http.Handle("/history", authMiddleware(http.HandlerFunc(historyHandler)))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
