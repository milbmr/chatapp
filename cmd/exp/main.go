package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/gorilla/websocket"
)

type client struct {
	send chan []byte
	room *room
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{WriteBufferSize: socketBufferSize, ReadBufferSize: socketBufferSize}

type room struct {
	mutex      sync.Mutex
	msgForward chan []byte
	clients    map[*client]bool
}

func (room *room) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("socket:", err)
		return
	}

	client := &client{
		send: make(chan []byte),
		room: room,
	}

	room.mutex.Lock()
	room.clients[client] = true
	room.mutex.Unlock()

	go func() {
		defer socket.Close()
		for msg := range client.send {
			err := socket.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		}
	}()

	defer func() {
		room.mutex.Lock()
		delete(room.clients, client)
		room.mutex.Unlock()
		close(client.send)
	}()

	for {
		_, msg, err := socket.ReadMessage()
		if err != nil {
			return
		}

		room.msgForward <- msg
	}
}

type temp struct {
	once sync.Once
	tlp  *template.Template
}

func (t *temp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.tlp = template.Must(template.ParseFiles(filepath.Join("templates/chat.html")))
	})

	t.tlp.Execute(w, nil)
}

func main() {
	temp := temp{}
	room := &room{
		msgForward: make(chan []byte),
		clients:    make(map[*client]bool),
	}

	http.Handle("/", &temp)
	http.Handle("/room", room)

	go func() {
		for msg := range room.msgForward {
			for client := range room.clients {
				client.send <- msg
			}
		}
	}()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("listenandserve:", err)
	}
}
