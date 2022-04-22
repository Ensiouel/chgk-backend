package gocket

import (
	"fmt"
)

type IRoom interface {
	Emit(string, EmitterData)
	On(string, EmitterFunc)
}

type room struct {
	name      string
	events    map[string]EmitterFunc
	join      chan *Socket
	leave     chan *Socket
	sockets   map[*Socket]bool
	broadcast chan []byte
}

func Room(name string) *room {
	return &room{
		name:      name,
		events:    map[string]EmitterFunc{},
		sockets:   map[*Socket]bool{},
		join:      make(chan *Socket),
		leave:     make(chan *Socket),
		broadcast: make(chan []byte),
	}
}

func (room *room) Emit(event string, data *EmitterData) {
	for socket := range room.sockets {
		socket.Emit(event, data)
	}
}

func (room *room) Run() {
	defer func() {
		close(room.join)
		close(room.leave)
		close(room.broadcast)
	}()
	for {
		select {
		case socket := <-room.join:
			socket.room = room
			room.sockets[socket] = true
			fmt.Printf("Room '%s' socket join:\tsocketId = %s\n", room.name, socket.id)
		case socket := <-room.leave:
			if _, ok := room.sockets[socket]; ok {
				delete(room.sockets, socket)
			}
			fmt.Printf("Room '%s' socket leave:\tsocketId = %s\n", room.name, socket.id)
		case message := <-room.broadcast:
			for socket := range room.sockets {
				select {
				case socket.send <- message:
				default:
					close(socket.send)
					delete(room.sockets, socket)
				}
			}
		}
	}
}

func (room *room) GetSockets() map[*Socket]bool {
	return room.sockets
}

func (room *room) Broadcast(data *EmitterData) {

}
