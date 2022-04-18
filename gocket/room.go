package gocket

type IRoom interface {
	Emit(string, EmitterData)
	On(string, EmitterFunc)
}

type room struct {
	events map[string]EmitterFunc
}

func Room() *room {
	return &room{
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
