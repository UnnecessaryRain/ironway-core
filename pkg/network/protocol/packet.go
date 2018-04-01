package protocol

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// Packet is used to send multiple messages to and from the server
// Packet can also wrap token in the upper level to save data
type Packet struct {
	Messages []Message `json:"messages,omitempty"`
}

// MarshalPacket turns a packet into json
func MarshalPacket(p Packet) []byte {
	b, e := json.Marshal(p)
	if e != nil {
		log.Errorln(e)
		return nil
	}
	return b
}

// UnmarshalPacket turns a json []byte into a packet
func UnmarshalPacket(b []byte) Packet {
	var p Packet
	e := json.Unmarshal(b, &p)
	if e != nil {
		log.Errorln(e)
	}
	return p
}
