package protocol

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// Message defines a single message or command from the client and server
// mulitple messages may be send in a packet
// however, messages are what send and received to and from the external interface
type Message string

// MarshalMessage turns a packet into json
func MarshalMessage(m Message) []byte {
	b, e := json.Marshal(m)
	if e != nil {
		log.Errorln(e)
		return nil
	}
	return b
}

// UnmarshalMessage turns a json []byte into a packet
func UnmarshalMessage(b []byte) Message {
	var m Message
	e := json.Unmarshal(b, &m)
	if e != nil {
		log.Errorln(e)
	}
	return m
}
