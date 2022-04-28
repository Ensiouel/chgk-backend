package gocket

import (
	"encoding/json"
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

type EmitRequest struct {
	Type  string      `json:"type"`
	Data  EmitterData `json:"data"`
	Event string      `json:"event"`
}

type EmitterData map[string]interface{}
type EmitterFunc func(data EmitterData)

func (data *EmitterData) Get(name string) reflect.Value {
	return reflect.ValueOf((*data)[name])
}

func (e *Emitter) Emit(event string, data *EmitterData) {
	request := EmitRequest{
		Type:  "emit",
		Data:  *data,
		Event: event,
	}
	b, _ := json.Marshal(&request)

	for _, socket := range e.sockets {
		go socket.EmitBytes(event, b)
	}
}
