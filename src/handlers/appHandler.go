package handlers

import (
	"database/sql"

	"github.com/gorilla/sessions"
)

type AppHandler struct {
	store *sessions.CookieStore
	db    *sql.DB
}

func NewAppHandler(store *sessions.CookieStore, db *sql.DB) *AppHandler {
	if store == nil {
		panic("Session cookiestore cannot be nil")
	}
	if db == nil {
		panic("Sql database cannot be nil")
	}
	return &AppHandler{store: store, db: db}
}
