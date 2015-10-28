package standup

import (
	"bytes"
	"fmt"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"os"
	"text/template"
)

var (
	token  = os.Getenv("HIPCHAT_TOKEN")
	roomId = os.Getenv("HIPCHAT_ROOM_ID")
)

const MESSAGE_TMPL = `{{.Name}}:
<ul>
<li><b>Yesterday</b>: {{.Yesterday}}</li>
<li><b>Today</b>: {{.Today}}</li>
<li><b>Blocked</b>: {{.Blocked}}</li>
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
		fmt.Fprintf(os.Stderr, "Error during Notify %q\n", err)
	}
	fmt.Printf("Success %+v\n", resp.StatusCode)
}

func (m *Message) formatted() string {
	var msg bytes.Buffer
	t, _ := template.New("Message").Parse(MESSAGE_TMPL)
	_ = t.Execute(&msg, m)
	return msg.String()
}

func (m *Message) color() string {
	if m.IsBlocked == true {
		return "red"
	}
	return "green"
}
