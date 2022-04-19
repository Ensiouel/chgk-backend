package gocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Gocket struct {
	upgrader websocket.Upgrader
	sockets  map[*Socket]bool
	rooms    map[string]*room
	events   map[string]func(*Socket)
}

func New() *Gocket {
	return &Gocket{
		rooms:   map[string]*room{},
		events:  map[string]func(*Socket){},
		sockets: map[*Socket]bool{},
	}
}

func (g *Gocket) OnConnection(f func(socket *Socket)) {
	g.events["connection"] = f
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
	g.sockets[socket] = true

	if f, ok := g.events["connection"]; ok {
		go f(socket)
	}

	fmt.Printf("New connection:\tsocketId = %s, connections = %d\n", socket.id, len(g.sockets))

	go socket.read()
	go socket.write()
}

func (g *Gocket) disconnect(socket *Socket) {
	if _, ok := g.sockets[socket]; ok {
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
		room := Room(name)
		go room.Run()
		room.join <- socket
		g.rooms[name] = room
	}

}

func (g *Gocket) GetRoom(name string) *room {
	if room, ok := g.rooms[name]; ok {
		return room
	}
	return &room{}
}
