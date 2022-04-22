package gocket

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Socket struct {
	IRoom
	id      uuid.UUID
	room    *room
	conn    *websocket.Conn
	gocket  *Gocket
	events  map[string]EmitterFunc
	close   chan struct{}
	send    chan []byte
	Storage map[string]string
}

func NewSocket(conn *websocket.Conn, gocket *Gocket) *Socket {
	return &Socket{
		id:      uuid.New(),
		conn:    conn,
		gocket:  gocket,
		events:  map[string]EmitterFunc{},
		close:   make(chan struct{}),
		send:    make(chan []byte),
		Storage: map[string]string{},
	}
}

func (socket *Socket) Emit(event string, data *EmitterData) {
	var emitRequest struct {
		Type  string      `json:"type"`
		Data  EmitterData `json:"data"`
		Event string      `json:"event"`
	}

	emitRequest.Type = "emit"
	emitRequest.Data = *data
	emitRequest.Event = event

	b, _ := json.Marshal(&emitRequest)

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

func (socket *Socket) GetRoom() *room {
	if socket.room != nil {
		return socket.room
	}
	return nil
}
