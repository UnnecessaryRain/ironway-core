package server

import (
	"net/http"

	"github.com/UnnecessaryRain/ironway-core/pkg/client"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveSocket(server *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := client.NewClient(conn, server.receivedChan, server.unregisterChan)
	server.registerChan <- client
	go client.StartWriter()
	go client.StartReader()
}
