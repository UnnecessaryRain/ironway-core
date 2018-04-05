package frame

import (
	"fmt"
	"time"

	"github.com/UnnecessaryRain/ironway-core/pkg/network/protocol"
)

// ChatFormat pretty formats a string message to outgoing message
func ChatFormat(sender string, message string, t int64) protocol.OutgoingMessage {
	return protocol.OutgoingMessage{
		Frame:   protocol.ChatFrame,
		Content: fmt.Sprintf("\\green{%v} %v: %v", time.Unix(t, 0).Format(time.Kitchen), sender, message),
		Mode:    protocol.AppendMode,
	}
}
