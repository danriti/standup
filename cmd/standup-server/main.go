package main

import (
	"github.com/danriti/standup"
	"github.com/danriti/standup/httputils"
	"net/http"
)

func main() {
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":3000", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		msg := &standup.Message{
			Name:      httputils.PostFormValue(r, "name", ""),
			Yesterday: httputils.PostFormValue(r, "yesterday", ""),
			Today:     httputils.PostFormValue(r, "today", ""),
			Blocked:   httputils.PostFormValue(r, "blocked", "Nope"),
			IsBlocked: httputils.PostFormBoolean(r, "is_blocked"),
		}
		msg.Notify()
	}
	httputils.RenderTemplate(w, "standup.html")
}
