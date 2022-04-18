package main

import (
	"fmt"
	"gocket/v2/gocket"
)

func main() {
	server := gocket.New()
	server.OnConnection(func(socket *gocket.Socket) {
		socket.On("new user", func(data gocket.EmitterData) {
			socket.Join("chat")
		})
	})

	room := gocket.Room("main")
	room.On("new user", func(data gocket.EmitterData) {
		userID := data.Get("user_id").String()
		userName := data.Get("user_name").String()
		date := data.Get("date").Int()
		fmt.Println(userID, userName, date)
	})

	room.Emit("new user", gocket.EmitterData{
		"user_id":   "2faf2qafsv",
		"user_name": "Rima",
		"date":      12312313234,
	})
}
