package handlers

import "net/http"

// Default screen handler
func (h *AppHandler) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/tabhome", http.StatusSeeOther)
}
