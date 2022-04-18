package main

import (
	"fmt"
	"gocket/v2/gocket"
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

		server.To("test chat").Emit("new user", gocket.EmitterData{
			"id": socket.GetID().String(),
		})

		socket.Join("test chat")
	})

	fmt.Println("The server is running...")

	http.Handle("/ws", server)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
