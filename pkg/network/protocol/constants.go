package protocol

import "time"

// Websocket constants
const (
	// WriteWait is the Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second

	// PongWait is the Time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second

	// PingPeriod is the pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * 9) / 10

	// MaxMessageSize message size allowed from peer.
	MaxMessageSize = 512
)

// Message constants
const (
	AppendMode  = "APPEND"
	ReplaceMode = "REPLACE"

	ChatFrame   = "chat"
	OnlineFrame = "online"
)
