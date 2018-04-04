package server

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/UnnecessaryRain/ironway-core/pkg/network/client"
	"github.com/UnnecessaryRain/ironway-core/pkg/network/protocol"
	log "github.com/sirupsen/logrus"
)

// Options defines config for the server to use
// passed in from the serveCommand
type Options struct {
	Addr string
}

// Server creates an endpoint for upgrading clients to websockets
// Also handles the connections between client ws
// receives messages from clients to be used by game core
type Server struct {
	Options

	Clients map[*client.Client]struct{}

	OnMessageHandler func(client.Message)

	receivedChan chan client.Message

	registerChan   chan *client.Client
	unregisterChan chan *client.Client

	broadcastChan chan protocol.OutgoingMessage
}

// NewServer creates a new Server object
func NewServer(options Options) *Server {
	server := &Server{
		Options:        options,
		Clients:        make(map[*client.Client]struct{}),
		registerChan:   make(chan *client.Client),
		unregisterChan: make(chan *client.Client),
		receivedChan:   make(chan client.Message, 256),
		broadcastChan:  make(chan protocol.OutgoingMessage, 256),
	}

	return server
}

// Send broadcasts the message to all connected clients
func (s *Server) Send(m protocol.OutgoingMessage) {
	s.broadcastChan <- m
}

// OnMessage callback for a client message received
func (s *Server) OnMessage(f func(client.Message)) {
	s.OnMessageHandler = f
}

func (s *Server) run(stopChan chan struct{}) {
	sendOnline := func() {
		var users bytes.Buffer

		users.WriteString("<u>[ ONLINE (")
		users.WriteString(fmt.Sprint(len(s.Clients)))
		users.WriteString(") ]</u>")
		for c := range s.Clients {
			users.WriteString("<br>")
			users.WriteString(c.Username)
			users.WriteString("\n")
		}

		s.broadcastChan <- protocol.OutgoingMessage{
			Frame:   protocol.OnlineFrame,
			Content: users.String(),
			Mode:    protocol.ReplaceMode,
		}
	}

	for {
		select {
		case client := <-s.registerChan:
			s.Clients[client] = struct{}{}
			log.Infoln("new client registered with username", client.Username)
			sendOnline()

		case client := <-s.unregisterChan:
			if _, ok := s.Clients[client]; ok {
				delete(s.Clients, client)
				close(client.SendChan)
			}
			sendOnline()
			log.Infoln("client disconnected")

		case clientMessage := <-s.receivedChan:
			s.OnMessageHandler(clientMessage)

		case globalMessage := <-s.broadcastChan:
			for c := range s.Clients {
				c.Send(globalMessage)
			}

		case <-stopChan:
			log.Infoln("stopping server")
			return
		}
	}
}

// ServeForever starts the server and blocks forever
func (s *Server) ServeForever(stopChan chan struct{}) {
	go s.run(stopChan)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveSocket(s, w, r)
	})
	err := http.ListenAndServe(s.Options.Addr, nil)
	if err != nil {
		log.Fatal("listenAndServe: ", err)
	}

	log.Info("server shutdown")
}
