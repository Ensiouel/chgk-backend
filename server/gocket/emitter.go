package gocket

import (
	"reflect"
)

const (
	EmitExceptSender = EmitterType(iota)
)

type EmitterType int

type Emitter struct {
	Type    EmitterType
	sender  *Socket
	sockets []*Socket
}

type EmitterData map[string]interface{}
type EmitterFunc func(data EmitterData)

func (data *EmitterData) Get(name string) reflect.Value {
	return reflect.ValueOf((*data)[name])
}

func (e *Emitter) Emit(event string, data EmitterData) {
	for _, socket := range e.sockets {
		socket.Emit(event, data)
	}
}
