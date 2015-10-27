package main

import (
	"fmt"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"html/template"
	"net/http"
	"os"
)

var (
	token  = os.Getenv("HIPCHAT_TOKEN")
	roomId = os.Getenv("HIPCHAT_ROOM_ID")
)

const MESSAGE = `%v:
<ul>
<li><b>Yesterday</b>: %v</li>
<li><b>Today</b>: %v</li>
<li><b>Blocked</b>: %v</li>
</ul>`

type Standup struct {
	name      string
	yesterday string
	today     string
	blocked   string
	isBlocked bool
}

func (s *Standup) Notify() {
	c := hipchat.NewClient(token)
	nr := &hipchat.NotificationRequest{
		Message:       s.message(),
		MessageFormat: "html",
		Color:         s.color(),
	}
	resp, err := c.Room.Notification(roomId, nr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during room notification %q\n", err)
	}
	fmt.Printf("Success %+v\n", resp.StatusCode)
}

func (s *Standup) message() string {
	return fmt.Sprintf(MESSAGE, s.name, s.yesterday, s.today, s.blocked)
}

func (s *Standup) color() string {
	if s.isBlocked == true {
		return "red"
	}
	return "green"
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		s := &Standup{
			name:      postFormValue(r, "name", ""),
			yesterday: postFormValue(r, "yesterday", ""),
			today:     postFormValue(r, "today", ""),
			blocked:   postFormValue(r, "blocked", "Nope"),
			isBlocked: postFormBoolean(r, "is_blocked"),
		}
		s.Notify()
	}
	renderTemplate(w, "standup.html")
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":3000", nil)
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
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

func postFormBoolean(r *http.Request, name string) bool {
	if r.PostFormValue(name) != "" {
		return true
	}
	return false
}

func postFormValue(r *http.Request, name string, fallback string) string {
	value := r.PostFormValue(name)
	if value != "" {
		return value
	}
	return fallback
}
