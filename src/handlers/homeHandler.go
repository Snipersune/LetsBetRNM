package handlers

import (
	"net/http"

	"github.com/snipersune/LetsBetRNM/src/renderers"
)

// Home screen handler
func (h *AppHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	renderers.RenderTemplate(w, "html/static/home.html", nil)
}
