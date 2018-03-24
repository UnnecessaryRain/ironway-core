package protocol

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type Packet struct {
	Messages []Message `json:"messages,omitempty"`
}

func MarshalPacket(p Packet) []byte {
	b, e := json.Marshal(p)
	if e != nil {
		log.Errorln(e)
		return nil
	}
	return b
}

func UnmarshalPacket(b []byte) Packet {
	var p Packet
	e := json.Unmarshal(b, &p)
	if e != nil {
		log.Errorln(e)
	}
	return p
}
