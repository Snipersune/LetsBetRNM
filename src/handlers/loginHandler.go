package handlers

import (
	"net/http"

	"github.com/snipersune/LetsBetRNM/src/renderers"
	"golang.org/x/crypto/bcrypt"
)

// Login page handler
func (h *AppHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage := r.URL.Query().Get("error")
	renderers.RenderTemplate(w, "html/static/login.html", map[string]string{"ErrorMessage": errorMessage})
}

// Handler to process log in attempts
func (h *AppHandler) ProcessLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Retrieve user from database
	var hashedPassword string
	err := h.db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
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
	err = h.db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		http.Redirect(w, r, "/login?error=Invalid%20username%20or%20password", http.StatusSeeOther)
		return
	}

	session, _ := h.store.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Values["user_id"] = userID
	session.Values["username"] = username
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
