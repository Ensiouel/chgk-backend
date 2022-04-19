package main

import (
	"chgk/gocket"
	"fmt"
	"log"
	"net/http"
)

func main() {
	server := gocket.New()

	server.OnConnection(func(socket *gocket.Socket) {
		for s := range server.GetRoom("test chat").GetSockets() {
			socket.Emit("new user", gocket.EmitterData{
				"id": s.GetID().String(),
			})
		}

		socket.Join("test chat")

		socket.Emit("your id", gocket.EmitterData{
			"id": socket.GetID().String(),
		})

		socket.To("test chat").Emit("new user", gocket.EmitterData{
			"id": socket.GetID().String(),
		})
	})

	server.OnDisconnecting(func(socket *gocket.Socket) {
		for s := range server.GetRoom("test chat").GetSockets() {
			s.Emit("leave user", gocket.EmitterData{
				"id": socket.GetID().String(),
			})
		}
	})

	fmt.Println("The server is running on port :8080...")

	http.Handle("/ws", server)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
