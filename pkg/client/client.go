package client

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/UnnecessaryRain/ironway-core/pkg/protocol"
	"github.com/gorilla/websocket"
)

// Message is a bundle of Client and protocol.Message
// Used for sending along the receivedChan and identifying the sender
type Message struct {
	Client  *Client
	Message *protocol.Message
}

// Client is an instance of a websocket client
// A client can send and receive messages through the SendChan and receiveChannel
type Client struct {
	conn *websocket.Conn

	SendChan chan protocol.Message

	receivedChan             chan<- Message
	unregisterFromServerChan chan<- *Client
}

// NewClient creates a new client object
func NewClient(conn *websocket.Conn, receiveChannel chan<- Message, unregisterChannel chan<- *Client) *Client {
	return &Client{
		conn:                     conn,
		SendChan:                 make(chan protocol.Message, 256),
		receivedChan:             receiveChannel,
		unregisterFromServerChan: unregisterChannel,
	}
}

// StartReader starts the reading pump from the websocket
// messages read will be sent to the passed in receivedChan
func (c *Client) StartReader() {
	defer func() {
		c.unregisterFromServerChan <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(protocol.MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(protocol.PongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(protocol.PongWait))
		return nil
	})

	for {
		_, packetBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warningln(err)
			}
			break
		}

		packet := protocol.UnmarshalPacket(packetBytes)
		log.Printf("message: %#v", packet)
		for _, message := range packet.Messages {
			c.receivedChan <- Message{c, &message}
		}
	}
}

// StartWriter starts a writer pump, writing messages to the websocket
// messages can be given to the client to write through the SendChan
func (c *Client) StartWriter() {
	ticker := time.NewTicker(protocol.PingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.SendChan:
			c.conn.SetWriteDeadline(time.Now().Add(protocol.WriteWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Errorln("failed to get writer", err)
				return
			}

			var nextPacket protocol.Packet
			nextPacket.Messages = append(nextPacket.Messages, message)

			// read the current queued packets
			n := len(c.SendChan)
			for i := 0; i < n; i++ {
				nextPacket.Messages = append(nextPacket.Messages, <-c.SendChan)
			}
			w.Write(protocol.MarshalPacket(nextPacket))
			if err := w.Close(); err != nil {
				log.Errorln("failed to close writer", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(protocol.WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
