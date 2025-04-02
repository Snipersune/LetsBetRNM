package main

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/snipersune/LetsBetRNM/internal/auth"
	"github.com/snipersune/LetsBetRNM/internal/fetch"
	"github.com/snipersune/LetsBetRNM/internal/handlers"

	_ "github.com/mattn/go-sqlite3"
)

/*
----------- TODOS -----------

Betting math
Repo/service/handler logic

*/

// Create a session store (store session data in memory)
var store = sessions.NewCookieStore([]byte("super-secret-key"))

// Database name
var dbFname = "biggestassdatabase.db"

func main() {

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,  // Session lasts for 1 hour (3600 seconds)
		HttpOnly: true,  // Prevent JavaScript access (security best practice)
		Secure:   false, // Require HTTPS
	}

	// Initialize database instance
	db, err := database.InitializeSqlDB(dbFname)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Initialize app handler
	h := handlers.NewAppHandler(store, db)

	go func() {
		http.HandleFunc("/", h.DefaultHandler)

		http.HandleFunc("/login", http.HandlerFunc(h.LoginHandler))
		http.HandleFunc("/register", http.HandlerFunc(h.RegisterHandler))
		http.HandleFunc("/logout", http.HandlerFunc(h.LogoutHandler))

		http.HandleFunc("/process-login", h.ProcessLoginHandler)
		http.HandleFunc("/process-register", h.ProcessRegisterHandler)
		http.HandleFunc("/check-authentication", h.CheckAuthenticationHandle)

		http.Handle("/process-powerplay", auth.AuthMiddleware(http.HandlerFunc(h.ProcessPowerplayHandler), store))

		http.Handle("/tabhome", auth.AuthMiddleware(http.HandlerFunc(h.HomeTabHandler), store))

		log.Println("Server running on http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	go fetch.FetchPowerplay()

	select {}
}
