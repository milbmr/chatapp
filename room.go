package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/milbmr/chatapp/trace"
)

type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	t       trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan *message, messageBufferSize),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		t:       trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case msg := <-r.forward:
			r.t.Trace("message recieved")
			for client := range r.clients {
				client.send <- msg
				r.t.Trace("message --sent to client")
			}
		case client := <-r.join:
			r.clients[client] = true
			r.t.Trace("new user joined")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.t.Trace("user left")
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("socket serveHttp:", err)
		return
	}

	var user user
	cookie, err := req.Cookie("userdata")
	if err != nil {
		log.Printf("error getting cookie %s", err.Error())
	}
	decodeFromString(cookie.Value, &user)

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: user,
	}

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
