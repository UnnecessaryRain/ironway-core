package client

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/UnnecessaryRain/ironway-core/pkg/network/protocol"
	"github.com/gorilla/websocket"
)

// Message is a bundle of Client and protocol.Message
// Used for sending along the receivedChan and identifying the sender
type Message struct {
	Sender   protocol.Sender
	Metadata *protocol.Metadata
	Message  *protocol.IncommingMessage
}

// Client is an instance of a websocket client
// A client can send and receive messages through the SendChan and receiveChannel
type Client struct {
	Username string

	conn *websocket.Conn

	SendChan chan protocol.OutgoingMessage

	receivedChan             chan<- Message
	unregisterFromServerChan chan<- *Client
	registerOnServerChan     chan<- *Client
}

// NewClient creates a new client object
func NewClient(conn *websocket.Conn, receiveChannel chan<- Message, registerChan, unregisterChannel chan<- *Client) *Client {
	return &Client{
		conn:                     conn,
		SendChan:                 make(chan protocol.OutgoingMessage, 256),
		receivedChan:             receiveChannel,
		registerOnServerChan:     registerChan,
		unregisterFromServerChan: unregisterChannel,
	}
}

// Send sends the message to the send channel and then to the client
func (c *Client) Send(m protocol.OutgoingMessage) {
	c.SendChan <- m
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

	timeout := time.NewTimer(5 * time.Second)

	go func() {
		<-timeout.C
		log.Errorf("timeout for addr %v. did not receive username in time. disconnecting", c.conn.RemoteAddr())
		c.conn.Close()
	}()

	for {
		_, packetBytes, err := c.conn.ReadMessage()
		timeout.Stop()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Debugln(err)
			}
			break
		}

		packet := protocol.UnmarshalPacket(packetBytes)
		log.Debugf("message: %#v", packet)
		for i, message := range packet.ServerMessages {

			// FIXME(#9): the first messaage from the client must be a one word string with no spaces
			// this will be their user name for now
			if len(c.Username) == 0 && i == 0 {
				if len(strings.Split(strings.TrimSpace(string(message)), "")) == 0 {
					log.Warningln("username expected as first message. got ", message)
					return
				}
				c.Username = strings.TrimSpace(string(message))
				c.registerOnServerChan <- c
				continue
			}
			c.receivedChan <- Message{c, &packet.Metadata, &message}
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
			nextPacket.ClientMessages = append(nextPacket.ClientMessages, message)

			// read the current queued packets
			n := len(c.SendChan)
			for i := 0; i < n; i++ {
				nextPacket.ClientMessages = append(nextPacket.ClientMessages, <-c.SendChan)
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
