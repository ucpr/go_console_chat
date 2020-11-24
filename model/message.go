package model

import (
	"fmt"
	"time"
)

type Message struct {
	Username  string
	Msg       string
	Timestamp time.Time
}

func NewMessage(username, message string) *Message {
	return &Message{
		Username:  username,
		Msg:       message,
		Timestamp: time.Now(),
	}
}

func (msg *Message) ToText() string {
	timestamp := msg.Timestamp.Format("2006/01/02 15:04:05")
	return fmt.Sprintf("%s | %s | %s\n", timestamp, msg.Username, msg.Msg)
}
