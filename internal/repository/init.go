package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func InitializeSqlDB(dbName string) (*sql.DB, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dbFname := fmt.Sprintf("%s/%s", cwd, dbName)

	// Initialize SQLite database
	db, err := sql.Open("sqlite3", dbFname)
	if err != nil {
		log.Fatal(err)
	}

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
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	// Create powerplay table if it doesn't exist
	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS powerplays (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			data TEXT NOT NULL
		)
		`)

	if err != nil {
		return nil, err
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
		return nil, err
	}
	return db, nil
}
