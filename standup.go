package standup

import (
	"fmt"
	"github.com/tbruyelle/hipchat-go/hipchat"
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

type Message struct {
	Name      string
	Yesterday string
	Today     string
	Blocked   string
	IsBlocked bool
}

func (s *Message) Notify() {
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

func (s *Message) message() string {
	return fmt.Sprintf(MESSAGE, s.Name, s.Yesterday, s.Today, s.Blocked)
}

func (s *Message) color() string {
	if s.IsBlocked == true {
		return "red"
	}
	return "green"
}
