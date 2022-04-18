package gocket

type IRoom interface {
	Emit(string, EmitterData)
	On(string, EmitterFunc)
}

type room struct {
	name    string
	events  map[string]EmitterFunc
	join    chan *Socket
	leave   chan *Socket
	sockets []*Socket
}

func Room(name string) *room {
	return &room{
		name:   name,
		events: map[string]EmitterFunc{},
	}
}

func (room *room) Emit(event string, data EmitterData) {
	if f, ok := room.events[event]; ok {
		f(data)
	}
}

func (room *room) On(event string, f EmitterFunc) {
	room.events[event] = f
}

func (room *room) Run() {

}
