package main

import (
	"chgk/gocket"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var Port = flag.String("port", "4221", "port")

func main() {
	flag.Parse()
	tables := map[string]*Table{}

	server := gocket.New()
	server.OnConnection(func(socket *gocket.Socket) {
		socket.On("init", func(data gocket.EmitterData) {
			tableID := data.Get("table_id").String()
			userID := data.Get("user_id").String()
			userName := data.Get("user_name").String()

			socket.Storage["user_id"] = userID
			socket.Storage["table_id"] = tableID

			if _, ok := tables[tableID]; !ok {
				tables[tableID] = NewTable()
			}

			table := tables[tableID]

			var user *User
			if table.ContainsUser(userID) == false {
				user = NewUser(userID, userName, UserRoleSpectator, socket, true)
				table.Users = append(table.Users, user)
			} else {
				user = table.GetUser(userID)
				user.Online = true
			}

			fmt.Println(tableID, userID, userName)

			socket.Join(tableID)
			socket.Emit("state", table.State())
			socket.To(tableID).Emit("user:join", user.State())

			// for client := range server.GetRoom(tableID).GetSockets() {
			// 	if client.Storage["user_id"] == userID {
			// 		client.Close( /*message*/ )
			// 	}
			// }

			/* TODO: Тут нужно отключать все дубликаты соединений */
		})
	})

	server.OnDisconnecting(func(socket *gocket.Socket) {
		tableID := socket.Storage["table_id"]
		userID := socket.Storage["user_id"]

		user := tables[tableID].GetUser(userID)
		if user != nil {
			user.Online = false
			/* Отправлять emit если все сокеты с userID отключены */
			socket.To(tableID).Emit("user:leave", user.State())
		}
	})

	fmt.Printf("The server is running on port :%s...\n", *Port)

	http.Handle("/ws", server)
	log.Fatal(http.ListenAndServe(":"+*Port, nil))
}
