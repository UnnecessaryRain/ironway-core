package protocol

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// Metadata about the player/client who sent this packet
type Metadata struct {
	Username  string `json:"username,omitempty"`
	Token     string `json:"token,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

// Packet is used to send multiple messages to and from the server
// Packet can also wrap token in the upper level to save data
type Packet struct {
	Metadata       Metadata           `json:"metadata,omitempty"`
	ClientMessages []OutgoingMessage  `json:"client_messages,omitempty"`
	ServerMessages []IncommingMessage `json:"server_messages,omitempty"`
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
