package httputils

import (
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func PostFormBoolean(r *http.Request, name string) bool {
	if r.PostFormValue(name) != "" {
		return true
	}
	return false
}

func PostFormValue(r *http.Request, name string, fallback string) string {
	value := r.PostFormValue(name)
	if value != "" {
		return value
	}
	return fallback
}
