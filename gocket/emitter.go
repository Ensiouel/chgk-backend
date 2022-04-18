package gocket

import "reflect"

type EmitterData map[string]interface{}
type EmitterFunc func(data EmitterData)

func (data *EmitterData) Get(name string) reflect.Value {
	return reflect.ValueOf((*data)[name])
}
