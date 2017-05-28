package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

func (this *client) read() {
	defer this.socket.Close()
	for {
		_, msg, err := this.socket.ReadMessage()
		if err != nil {
			return
		}
		this.room.forward <- msg
	}
}

func (this *client) write() {
	defer this.socket.Close()
	for msg := range this.send {
		err := this.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
