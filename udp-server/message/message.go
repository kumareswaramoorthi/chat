package message

import (
	"strings"
	"time"
)

type Message struct {
	Action    string
	UserName  string
	MessageID string
	Content   string
	TimeStamp string
}

// ParseMessage function decodes the incoming
// raw string and creates a Message struct variable.
func ParseMessage(msg string) (m Message) {
	msgsReceived := strings.Split(msg, "|")

	switch msgsReceived[0] {
	case "CONNECT":
		m.Action = msgsReceived[0]
		m.UserName = msgsReceived[1]
		m.TimeStamp = time.Now().Format("2006-01-02 3:4:5 pm")
	case "DISCONNECT":
		m.Action = msgsReceived[0]
		m.UserName = msgsReceived[1]
		m.TimeStamp = time.Now().Format("2006-01-02 3:4:5 pm")
	case "BROADCAST":
		m.Action = msgsReceived[0]
		m.UserName = msgsReceived[1]
		m.MessageID = msgsReceived[2]
		m.Content = msgsReceived[3]
		m.TimeStamp = time.Now().Format("2006-01-02 3:4:5 pm")
	case "DELETE":
		m.Action = msgsReceived[0]
		m.UserName = msgsReceived[1]
		m.MessageID = msgsReceived[2]
		m.TimeStamp = time.Now().Format("2006-01-02 3:4:5 pm")
	}
	return
}
