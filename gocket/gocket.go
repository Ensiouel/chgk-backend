package gocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Gocket struct {
	upgrader websocket.Upgrader
	sockets  map[*Socket]bool
	rooms    map[string]*Room
	events   map[string]func(*Socket)
	m        sync.Mutex
}

func New() *Gocket {
	return &Gocket{
		rooms:   map[string]*Room{},
		events:  map[string]func(*Socket){},
		sockets: map[*Socket]bool{},
	}
}

func (g *Gocket) OnConnection(f func(socket *Socket)) {
	g.events["connection"] = f
}

func (g *Gocket) OnDisconnecting(f func(socket *Socket)) {
	g.events["disconnect"] = f
}

func (g *Gocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := g.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	socket := NewSocket(conn, g)
	g.m.Lock()
	g.sockets[socket] = true
	g.m.Unlock()

	go socket.read()
	go socket.write()

	if f, ok := g.events["connection"]; ok {
		f(socket)
	}
}

func (g *Gocket) disconnect(socket *Socket) {
	if _, ok := g.sockets[socket]; ok {
		if f, ok := g.events["disconnect"]; ok {
			f(socket)
		}
		delete(g.sockets, socket)
	}
}

func (g *Gocket) join(name string, socket *Socket) {
	if socket == nil {
		fmt.Println("something wrong", socket.room)
		return
	}

	if socket.room != nil {
		socket.room.leave <- socket
	}

	if room, ok := g.rooms[name]; ok {
		room.join <- socket
	} else {
		room := NewRoom(name)
		go room.Run()
		room.join <- socket
		g.rooms[name] = room
	}

}

func (g *Gocket) GetRoom(name string) *Room {
	if room, ok := g.rooms[name]; ok {
		return room
	}
	return &Room{}
}
