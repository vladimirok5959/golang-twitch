package pubsub

import (
	"encoding/json"
	"fmt"
)

type AnswerType string

const (
	Listen    AnswerType = "LISTEN"
	Message   AnswerType = "MESSAGE"
	Ping      AnswerType = "PING"
	Pong      AnswerType = "PONG"
	Reconnect AnswerType = "RECONNECT"
	Response  AnswerType = "RESPONSE"
	Unlisten  AnswerType = "UNLISTEN"
)

func (a AnswerType) String() string {
	return string(a)
}

// -----------------------------------------------------------------------------

type Answer struct {
	processed bool

	Type  AnswerType `json:"type"`
	Data  any        `json:"data,omitempty"`
	Error string     `json:"error,omitempty"`
	Nonce string     `json:"nonce,omitempty"`
}

func (a *Answer) Parse() {
	if a.processed {
		return
	}
	a.processed = true

	data := AnswerDataMessage{}

	switch v := a.Data.(type) {
	case map[string]any:
		for fn, fv := range v {
			if fn == "message" {
				data.Message = fmt.Sprintf("%s", fv)
			} else if fn == "topic" {
				data.Topic = fmt.Sprintf("%s", fv)
			}
		}
	default:
	}

	a.Data = data
}

func (a Answer) HasError() bool {
	return a.Error != ""
}

func (a Answer) JSON() []byte {
	bytes, _ := json.Marshal(a)
	return bytes
}

// -----------------------------------------------------------------------------

type AnswerDataMessage struct {
	Message string `json:"message"`
	Topic   string `json:"topic"`
}

func (a AnswerDataMessage) JSON() []byte {
	bytes, _ := json.Marshal(a)
	return bytes
}

// -----------------------------------------------------------------------------

type AnswerDataTopics struct {
	Topics []string `json:"topics"`
}

func (a AnswerDataTopics) JSON() []byte {
	bytes, _ := json.Marshal(a)
	return bytes
}
