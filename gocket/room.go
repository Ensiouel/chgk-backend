package gocket

import (
	"fmt"
	"sync"
)

type IRoom interface {
	Emit(string, EmitterData)
	On(string, EmitterFunc)
}

type Room struct {
	name      string
	events    map[string]EmitterFunc
	join      chan *Socket
	leave     chan *Socket
	sockets   map[*Socket]bool
	broadcast chan []byte
	m         sync.Mutex
}

func NewRoom(name string) *Room {
	return &Room{
		name:      name,
		events:    map[string]EmitterFunc{},
		sockets:   map[*Socket]bool{},
		join:      make(chan *Socket),
		leave:     make(chan *Socket),
		broadcast: make(chan []byte),
	}
}

func (room *Room) Emit(event string, data *EmitterData) {
	for socket := range room.sockets {
		go socket.Emit(event, data)
	}
}

func (room *Room) Run() {
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

func (room *Room) GetSockets() map[*Socket]bool {
	return room.sockets
}
