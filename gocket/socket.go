package gocket

import (
	"fmt"

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
		events: map[string]EmitterFunc{},
	}
}

func (socket *Socket) Emit(event string, data EmitterData) {
	if f, ok := socket.events[event]; ok {
		go f(data)
	}
}

func (socket *Socket) On(event string, f EmitterFunc) {
	socket.events[event] = f
}

func (socket *Socket) GetID() uuid.UUID {
	return socket.id
}

type SocketEvent struct {
	Event string      `json:"event"`
	Data  EmitterData `json:"data"`
}

func (socket *Socket) read() {
	defer func() {
		socket.conn.Close()
		socket.gocket.disconnect(socket)
		socket.room.leave <- socket
	}()
	for {
		event := SocketEvent{}
		if err := socket.conn.ReadJSON(&event); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			break
		}

		if f, ok := socket.events[event.Event]; ok {
			go f(event.Data)
		}
	}
}

func (socket *Socket) write() {
	defer func() {
		socket.conn.Close()
		socket.gocket.disconnect(socket)
		socket.room.leave <- socket
	}()
	for {
	}
}

func (socket *Socket) Join(name string) {
	go socket.gocket.join(name, socket)
}
