package main

import (
	"fmt"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"os"
)

var (
	token  = os.Getenv("HIPCHAT_TOKEN")
	roomId = os.Getenv("HIPCHAT_ROOM_ID")
)

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
	msg := `%v:
<ul>
<li><b>Yesterday</b>: %v</li>
<li><b>Today</b>: %v</li>
<li><b>Blocked</b>: %v</li>
</ul>`
	return fmt.Sprintf(msg, s.name, s.yesterday, s.today, s.blocked)
}

func (s *Standup) color() string {
	if s.isBlocked == true {
		return "red"
	}
	return "green"
}

func main() {
	s := &Standup{
		name:      "Dan",
		yesterday: "Foo",
		today:     "Bar",
		blocked:   "Nope",
		isBlocked: false,
	}
	s.Notify()
}
