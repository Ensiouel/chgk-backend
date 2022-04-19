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
		socket.On("new user", func(data gocket.EmitterData) {
			socketID := data.Get("id").String()
			fmt.Println(socketID)
		})

		for s := range server.GetRoom("test chat").GetSockets() {
			fmt.Println(socket.GetID())
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
		fmt.Println("мужик ушел")
	})

	fmt.Println("The server is running on port :8080...")

	http.Handle("/ws", server)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
