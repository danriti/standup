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

const MESSAGE_TMPL = `%v:
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

func (m *Message) Notify() {
	c := hipchat.NewClient(token)
	nr := &hipchat.NotificationRequest{
		Message:       m.formatted(),
		MessageFormat: "html",
		Color:         m.color(),
	}
	resp, err := c.Room.Notification(roomId, nr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during room notification %q\n", err)
	}
	fmt.Printf("Success %+v\n", resp.StatusCode)
}

func (m *Message) formatted() string {
	return fmt.Sprintf(MESSAGE_TMPL, m.Name, m.Yesterday, m.Today, m.Blocked)
}

func (m *Message) color() string {
	if m.IsBlocked == true {
		return "red"
	}
	return "green"
}
