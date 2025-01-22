package main

import (
	"fmt"
	"net/http"
	"log"
	"os/exec"
	"html/template"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)


/*
----------- TODOS -----------

Authentication - User-password (or fake)
Sqlite - Save user data
Betting math
Add html sites to view betting history and stats

*/

// Render templates
func renderTemplate(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// Database instance
var db *sql.DB


// Default screen handler
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "template/default.html", nil)
}

// Login page handler
func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage := r.URL.Query().Get("error")
	renderTemplate(w, "template/login.html", map[string]string{"ErrorMessage": errorMessage})
}

// Register page handler
func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("server: registerPageHandler requested")
	fmt.Printf("server: request type: %s\n", r.Method)
	errorMessage := r.URL.Query().Get("error")
	renderTemplate(w, "template/register.html", map[string]string{"ErrorMessage": errorMessage})
}

// Dashboard handler
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "template/dashboard.html", nil)
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

	fmt.Println("Registering completed. Redirecting to login")
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

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}


func powerplay(w http.ResponseWriter, r *http.Request) {
	fmt.Println("server: powerplay handler started")
	defer fmt.Println("server: powerplay handler ended")

	c := exec.Command("/home/dnee/Documents/Kod/LetsBetRNM/.venv/bin/python", "/home/dnee/Documents/Kod/LetsBetRNM/src/update_powerplay_html.py")
	if err := c.Run(); err != nil {
		fmt.Println("server: Error: ", err)
	}

	http.ServeFile(w, r, "output/powerplay_rendered.html")
}

func main() {

	var err error

	// Initialize SQLite database
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create users table if it doesn't exist
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT
	)
	`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", defaultHandler)

	http.HandleFunc("/login", loginPageHandler)
	http.HandleFunc("/register", registerPageHandler)

	http.HandleFunc("/process-login", loginHandler)
	http.HandleFunc("/process-register", registerHandler)

	http.HandleFunc("/dashboard", dashboardHandler)

	http.HandleFunc("/powerplay", powerplay)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
