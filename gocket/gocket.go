package gocket

import "net/http"

type Gocket struct {
	sockets []*socket
}

func New() *Gocket {
	return &Gocket{}
}

func (g *Gocket) OnConnection(f func(socket *socket)) {

}

func (g *Gocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
