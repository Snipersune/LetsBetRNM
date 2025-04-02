package handlers

import (
	"fmt"
	"net/http"
)

func (h *AppHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")

	fmt.Println("Logout")

	// Remove session values
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1 // Expire immediately

	// Save the session to clear it
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
