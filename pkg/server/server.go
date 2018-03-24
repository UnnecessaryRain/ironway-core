package server

import (
	"net/http"

	"github.com/UnnecessaryRain/ironway-core/pkg/protocol"

	"github.com/UnnecessaryRain/ironway-core/pkg/client"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	Addr string
}

type Server struct {
	Options

	Clients map[*client.Client]struct{}

	receivedChan chan client.Message

	registerChan   chan *client.Client
	unregisterChan chan *client.Client

	stopChan chan struct{}
}

func NewServer(options Options, stop chan struct{}) *Server {
	server := &Server{
		Options:        options,
		Clients:        make(map[*client.Client]struct{}),
		registerChan:   make(chan *client.Client),
		unregisterChan: make(chan *client.Client),
		receivedChan:   make(chan client.Message, 256),
		stopChan:       stop,
	}

	return server
}

func (s *Server) run() {
	for {
		select {
		case client := <-s.registerChan:
			s.Clients[client] = struct{}{}
			log.Infoln("new client registered")
		case client := <-s.unregisterChan:
			if _, ok := s.Clients[client]; ok {
				delete(s.Clients, client)
				close(client.SendChan)
			}
			log.Infoln("client disconnected")
		case clientMessage := <-s.receivedChan:
			log.Printf("%s\n", *clientMessage.Message)
			clientMessage.Client.SendChan <- protocol.Message("Omg message back from server!")
		case <-s.stopChan:
			log.Infoln("stopping server")
			return
		}
	}
}

func (s *Server) ServeForever() {
	go s.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveSocket(s, w, r)
	})
	err := http.ListenAndServe(s.Options.Addr, nil)
	if err != nil {
		log.Fatal("listenAndServe: ", err)
	}

	log.Info("server shutdown")
}
