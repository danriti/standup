package standup

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/tbruyelle/hipchat-go/hipchat"
)

var (
	token  = os.Getenv("HIPCHAT_TOKEN")
	roomId = os.Getenv("HIPCHAT_ROOM_ID")
)

type Message struct {
	Name      string
	Yesterday string
	Today     string
	Blocked   string
	IsBlocked bool
}

func (m *Message) Notify() (success bool, err error) {
	c := hipchat.NewClient(token)
	nr := &hipchat.NotificationRequest{
		Message:       m.formatted(),
		MessageFormat: "html",
		Color:         m.color(),
	}
	resp, err := c.Room.Notification(roomId, nr)
	if err != nil {
		log.Println("error notifying room", err)
		return false, err
	}
	log.Println("success notifying room", resp.StatusCode)
	return resp.StatusCode == http.StatusNoContent, err
}

func (m *Message) formatted() string {
	var msg bytes.Buffer
	t, _ := template.New("Message").Parse(messageTemplate)
	t.Execute(&msg, m)
	return msg.String()
}

func (m *Message) color() string {
	if m.IsBlocked == true {
		return "red"
	}
	return "green"
}

const messageTemplate = `{{.Name}}:
<ul>
<li><b>Yesterday</b>: {{.Yesterday}}</li>
<li><b>Today</b>: {{.Today}}</li>
<li><b>Blocked</b>: {{.Blocked}}</li>
</ul>`
