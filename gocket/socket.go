package gocket

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Socket struct {
	IRoom
	id        uuid.UUID
	room      *Room
	conn      *websocket.Conn
	gocket    *Gocket
	events    map[string]EmitterFunc
	close     chan struct{}
	send      chan []byte
	Storage   map[string]string
	Broadcast *Broadcast
}

func NewSocket(conn *websocket.Conn, gocket *Gocket) *Socket {
	return &Socket{
		id:        uuid.New(),
		conn:      conn,
		gocket:    gocket,
		events:    map[string]EmitterFunc{},
		close:     make(chan struct{}),
		send:      make(chan []byte),
		Storage:   map[string]string{},
		Broadcast: NewBroadcast(),
	}
}

func (socket *Socket) Emit(event string, data *EmitterData) {
	request := EmitRequest{}

	request.Type = "emit"
	request.Data = *data
	request.Event = event
	b, _ := json.Marshal(&request)

	socket.send <- b
}

func (socket *Socket) EmitBytes(event string, b []byte) {
	socket.send <- b
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
		if socket.room != nil {
			socket.room.leave <- socket
		}
		socket.gocket.disconnect(socket)
	}()
	for {
		event := SocketEvent{}
		if err := socket.conn.ReadJSON(&event); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v\n", err)
			}
			break
		}

		if f, ok := socket.events[event.Event]; ok {
			f(event.Data)
		}
	}
}

func (socket *Socket) write() {
	defer func() {
		socket.conn.Close()
	}()
	for {
		select {
		case message, ok := <-socket.send:
			if !ok {
				socket.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := socket.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-socket.close:
			return
		}
	}
}

func (socket *Socket) Join(name string) {
	go socket.gocket.join(name, socket)
}

func (socket *Socket) To(name string) *Emitter {
	if room, ok := socket.gocket.rooms[name]; ok {
		var sockets []*Socket
		for s := range room.sockets {
			if s == socket {
				continue
			}
			sockets = append(sockets, s)
		}
		return &Emitter{
			Type:    EmitExceptSender,
			sockets: sockets,
		}
	}
	return &Emitter{}
}

func (socket *Socket) GetRoom() *Room {
	if socket.room != nil {
		return socket.room
	}
	return nil
}

func (socket *Socket) Close() {
	socket.conn.Close()
	if socket.room != nil {
		socket.room.leave <- socket
	}
	socket.gocket.disconnect(socket)
}
