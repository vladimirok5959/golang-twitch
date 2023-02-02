package pubsub

import "encoding/json"

type Response struct {
	Type  string      `json:"type"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
	Nonce string      `json:"nonce,omitempty"`
}

func (r Response) JSON() []byte {
	bytes, _ := json.Marshal(r)
	return bytes
}

type DataTopics struct {
	Topics []string `json:"topics"`
}

func (d DataTopics) JSON() []byte {
	bytes, _ := json.Marshal(d)
	return bytes
}
