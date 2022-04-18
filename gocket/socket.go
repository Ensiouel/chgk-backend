package gocket

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Socket struct {
	IRoom
	id     uuid.UUID
	room   *room
	conn   *websocket.Conn
	gocket *Gocket
	events map[string]EmitterFunc
}

func NewSocket(conn *websocket.Conn, gocket *Gocket) *Socket {
	return &Socket{
		id:     uuid.New(),
		conn:   conn,
		gocket: gocket,
	}
}

func (socket *Socket) Emit(event string, data EmitterData) {
	if f, ok := socket.events[event]; ok {
		f(data)
	}
}

func (socket *Socket) On(event string, f EmitterFunc) {
	socket.events[event] = f
}

func (socket *Socket) read() {
	defer func() {
		socket.conn.Close()
		socket.gocket.disconnect(socket)
	}()
	for {
	}
}

func (socket *Socket) write() {
	defer func() {
		socket.conn.Close()
		socket.gocket.disconnect(socket)
	}()
	for {
	}
}

func (socket *Socket) Join(name string) {
	socket.gocket.join(name, socket)
}
