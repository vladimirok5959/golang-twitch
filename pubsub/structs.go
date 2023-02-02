package pubsub

import (
	"encoding/json"
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

type Answer struct {
	Type  AnswerType  `json:"type"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
	Nonce string      `json:"nonce,omitempty"`
}

func (a Answer) HasError() bool {
	return a.Error != ""
}

func (a Answer) JSON() []byte {
	bytes, _ := json.Marshal(a)
	return bytes
}

type AnswerDataTopics struct {
	Topics []string `json:"topics"`
}

func (a AnswerDataTopics) JSON() []byte {
	bytes, _ := json.Marshal(a)
	return bytes
}
