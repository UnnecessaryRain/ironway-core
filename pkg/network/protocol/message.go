package protocol

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// IncommingMessage defines a single message or command from the client to server
// mulitple messages may be send in a packet
// however, messages are what send and received to and from the external interface
type IncommingMessage string

// MarshalIncommingMessage turns a packet into json
func MarshalIncommingMessage(m IncommingMessage) []byte {
	b, e := json.Marshal(m)
	if e != nil {
		log.Errorln(e)
		return nil
	}
	return b
}

// UnmarshalIncommingMessage turns a json []byte into a packet
func UnmarshalIncommingMessage(b []byte) IncommingMessage {
	var m IncommingMessage
	e := json.Unmarshal(b, &m)
	if e != nil {
		log.Errorln(e)
	}
	return m
}

// OutgoingMessage defines a message to a client from the server
type OutgoingMessage struct {
	Frame   string `json:"frame,omitempty"`
	Content string `json:"content,omitempty"`
	Mode    string `json:"mode,omitempty"`
}

// MarshalOutgoingMessage turns a packet into json
func MarshalOutgoingMessage(m OutgoingMessage) []byte {
	b, e := json.Marshal(m)
	if e != nil {
		log.Errorln(e)
		return nil
	}
	return b
}

// UnmarshalOutgoingMessage turns a json []byte into a packet
func UnmarshalOutgoingMessage(b []byte) OutgoingMessage {
	var m OutgoingMessage
	e := json.Unmarshal(b, &m)
	if e != nil {
		log.Errorln(e)
	}
	return m
}
