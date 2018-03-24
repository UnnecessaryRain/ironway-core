package protocol

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type Message string

func MarshalMessage(m Message) []byte {
	b, e := json.Marshal(m)
	if e != nil {
		log.Errorln(e)
		return nil
	}
	return b
}

func UnmarshalMessage(b []byte) Message {
	var m Message
	e := json.Unmarshal(b, &m)
	if e != nil {
		log.Errorln(e)
	}
	return m
}
