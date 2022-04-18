package gocket

import "net/http"

type Gocket struct {
	sockets []*Socket
}

func New() *Gocket {
	return &Gocket{}
}

func (g *Gocket) OnConnection(f func(socket *Socket)) {

}

func (g *Gocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
