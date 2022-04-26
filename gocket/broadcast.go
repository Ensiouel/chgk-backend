package gocket

func NewBroadcast() *Broadcast {
	return &Broadcast{
		buffer: make(chan []byte),
	}
}

type Broadcast struct {
	buffer chan []byte
}

func (b *Broadcast) Write(message []byte) {
	b.buffer <- message
}

func (b *Broadcast) Read() <-chan []byte {
	return (*b).buffer
}

func (b *Broadcast) Emit(event string, data *EmitterData) {

}
