package main

import (
	"log"
	"net/http"

	"github.com/danriti/standup"
	"github.com/danriti/standup/httputils"
)

func main() {
	http.HandleFunc("/", indexHandler)
	log.Println("starting server")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	nr := &notifyResponse{Success: false, Failure: false}
	if r.Method == "POST" {
		r.ParseForm()
		msg := &standup.Message{
			Name:      httputils.PostFormValue(r, "name", ""),
			Yesterday: httputils.PostFormValue(r, "yesterday", ""),
			Today:     httputils.PostFormValue(r, "today", ""),
			Blocked:   httputils.PostFormValue(r, "blocked", "Nope"),
			IsBlocked: httputils.PostFormBoolean(r, "is_blocked"),
		}
		success, err := msg.Notify()
		nr.Success = success
		nr.Failure = !success || err != nil
	}
	httputils.RenderTemplate(w, "standup.html", nr)
}

type notifyResponse struct {
	Success bool
	Failure bool
}
