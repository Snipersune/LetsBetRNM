package main

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/snipersune/LetsBetRNM/src/auth"
	"github.com/snipersune/LetsBetRNM/src/database"
	"github.com/snipersune/LetsBetRNM/src/fetch"
	"github.com/snipersune/LetsBetRNM/src/handlers"

	_ "github.com/mattn/go-sqlite3"
)

/*
----------- TODOS -----------

Betting math
Add html sites to view betting history and stats

*/

// Create a session store (store session data in memory)
var store = sessions.NewCookieStore([]byte("super-secret-key"))

// Database name
var dbFname = "biggestassdatabase.db"

func main() {

	// Initialize database instance
	db, err := database.SqlInitDb(dbFname)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Initialize app handler
	h := handlers.NewAppHandler(store, db)

	go func() {
		http.HandleFunc("/", h.DefaultHandler)
		http.Handle("/home", auth.AuthMiddleware(http.HandlerFunc(h.HomeHandler), store))

		http.HandleFunc("/login", h.LoginHandler)
		http.HandleFunc("/register", h.RegisterHandler)

		http.HandleFunc("/process-login", h.ProcessLoginHandler)
		http.HandleFunc("/process-register", h.ProcessRegisterHandler)

		http.Handle("/powerplay", auth.AuthMiddleware(http.HandlerFunc(h.PowerplayHandler), store))
		http.Handle("/process-powerplay", auth.AuthMiddleware(http.HandlerFunc(h.ProcessPowerplayHandler), store))

		http.Handle("/history", auth.AuthMiddleware(http.HandlerFunc(h.HistoryHandler), store))
		http.Handle("/tabhome", auth.AuthMiddleware(http.HandlerFunc(h.HomeTabHandler), store))

		log.Println("Server running on http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	go fetch.FetchPowerplay()

	select {}
}
