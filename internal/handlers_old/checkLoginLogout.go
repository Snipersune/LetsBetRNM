package handlers

import (
	"fmt"
	"net/http"
)

func (h *AppHandler) isAuthenticated(r *http.Request) bool {
	session, _ := h.store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}

func (h *AppHandler) CheckAuthenticationHandle(w http.ResponseWriter, r *http.Request) {
	if h.isAuthenticated(r) {
		fmt.Fprintln(w, "Authenticated")
	} else {
		fmt.Fprintln(w, "Unauthorized")
	}
}
