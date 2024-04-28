package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type message struct {
	Name    string
	Message string
	Time    time.Time
}

type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData user
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}

		msg.Name = c.userData.Name
		msg.Time = time.Now()

		c.room.forward <- msg
	}
}
