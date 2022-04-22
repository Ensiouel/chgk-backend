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
		socket.On("user connected", func(data gocket.EmitterData) {
			roomID := data.Get("room_id").String()

			socket.Storage = map[string]string{
				"local_user_id": data.Get("local_user_id").String(),
			}

			for s := range server.GetRoom(roomID).GetSockets() {
				socket.Emit("join user", &gocket.EmitterData{
					"id":            s.GetID().String(),
					"local_user_id": s.Storage["local_user_id"],
				})
			}

			socket.To(roomID).Emit("join user", &gocket.EmitterData{
				"id":            socket.GetID().String(),
				"local_user_id": socket.Storage["local_user_id"],
			})

			socket.Join(roomID)
		})

		socket.Emit("your id", &gocket.EmitterData{
			"id": socket.GetID().String(),
		})
	})

	server.OnDisconnecting(func(socket *gocket.Socket) {
		if socket.GetRoom() != nil {
			socket.GetRoom().Emit("leave user", &gocket.EmitterData{
				"id":            socket.GetID().String(),
				"local_user_id": socket.Storage["local_user_id"],
			})
		}
	})

	fmt.Println("The server is running on port :4221...")

	http.Handle("/ws", server)
	log.Fatal(http.ListenAndServe(":4221", nil))
}
