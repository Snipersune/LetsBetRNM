package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/snipersune/LetsBetRNM/src/renderers"
	"golang.org/x/crypto/bcrypt"
)

// Register page handler
func (h *AppHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("server: registerHandler requested")
	fmt.Printf("server: request type: %s\n", r.Method)

	session, err := h.store.Get(r, "session-name")
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		userID, _ := session.Values["user_id"].(int)
		fmt.Printf("User %d already logged in! Cannot register.\n", userID)
		http.Redirect(w, r, "html/static/tabhome.html", http.StatusSeeOther)
		return
	}

	errorMessage := r.URL.Query().Get("error")
	renderers.RenderTemplate(w, "html/static/register.html", map[string]string{"ErrorMessage": errorMessage})
}

// Handler to process register attempts
func (h *AppHandler) ProcessRegisterHandler(w http.ResponseWriter, r *http.Request) {

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
	_, err = h.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, string(hashedPassword))
	if err != nil {
		http.Redirect(w, r, "/register?error=Username%20already%20exists", http.StatusSeeOther)
		return
	}

	fmt.Println("Registration completed. Redirecting to login.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
