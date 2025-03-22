package renderers

import (
	"fmt"
	"net/http"
	"text/template"
)

// Render templates
func RenderTemplate(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		fmt.Fprintf(w, "Error: %d\n", err)
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
