package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func (this *room) run() {
	for {
		select {
		case client := <-this.join:
			this.clients[client] = true
		case client := <-this.leave:
			delete(this.clients, client)
			close(client.send)
		case msg := <-this.forward:
			for client := range this.clients {
				client.send <- msg
			}
		}
	}
}

func (this *room) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   this,
	}
	this.join <- client
	defer func() {
		this.leave <- client
	}()
	go client.write()
	client.read()
}
