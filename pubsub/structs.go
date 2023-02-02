package pubsub

import "encoding/json"

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

type Answer struct {
	Type  AnswerType  `json:"type"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
	Nonce string      `json:"nonce,omitempty"`
}

func (r Answer) JSON() []byte {
	bytes, _ := json.Marshal(r)
	return bytes
}

type AnswerDataTopics struct {
	Topics []string `json:"topics"`
}

func (d AnswerDataTopics) JSON() []byte {
	bytes, _ := json.Marshal(d)
	return bytes
}
