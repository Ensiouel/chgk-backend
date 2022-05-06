package main

import (
	chgkdb "chgk/chgk-db"
	"chgk/gocket"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var Port = flag.String("port", "4221", "port")

func main() {
	flag.Parse()

	session := chgkdb.NewSession()

	var res struct {
		Tours []struct {
			Questions []struct {
				ID           string `json:"id"`
				Number       int64  `json:"number"`
				Type         string `json:"type"`
				Question     string `json:"question"`
				Answer       string `json:"answer"`
				PassCriteria string `json:"passCriteria"`
				Authors      string `json:"authors"`
				Sources      string `json:"sources"`
				Comments     string `json:"comments"`
			} `json:"questions"`
		} `json:"tours"`
	}
	err := session.Call("packages/18plus", nil).Execute(&res)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("r", len(res.Tours[1].Questions))

	tables := map[string]*Table{}

	server := gocket.New()
	server.OnConnection(func(socket *gocket.Socket) {
		socket.On("init", func(data gocket.EmitterData) {
			tableID := data.Get("table_id").String()
			userID := data.Get("user_id").String()
			userName := data.Get("user_name").String()

			socket.Storage["user_id"] = userID
			socket.Storage["table_id"] = tableID

			var newUserIsMaster bool
			if _, ok := tables[tableID]; !ok {
				tables[tableID] = NewTable()
				newUserIsMaster = true
			}

			table := tables[tableID]

			var user *User
			if table.ContainsUser(userID) == false {
				userRole := UserRoleSpectator
				if newUserIsMaster {
					userRole = UserRoleMaster
				}
				user = NewUser(userID, userName, userRole, socket, true)

				table.Users[user] = true
			} else {
				user = table.GetUser(userID)
				user.Online = true
			}

			socket.Join(tableID)
			socket.Emit("state", table.State())
			server.GetRoom(tableID).Emit("state", table.State())
		})

		socket.On("pack:changeId", func(data gocket.EmitterData) {
			tableID := socket.Storage["table_id"]
			table := tables[tableID]

			table.ranges = ranges(12)
			table.QuestionsPlayed = []int{}
			table.PackID = data.Get("pack_id").String()

			server.GetRoom(tableID).Emit("state", table.State())
		})

		socket.On("role:change", func(data gocket.EmitterData) {
			tableID := socket.Storage["table_id"]
			table := tables[tableID]
			userID := data.Get("user_id").String()
			newRole := data.Get("new_role").Float()

			if table.Open {
				table.GetUser(userID).Role = UserRole(newRole)
				server.GetRoom(tableID).Emit("state", table.State())
			}
		})

		socket.On("question:choice", func(data gocket.EmitterData) {
			tableID := socket.Storage["table_id"]
			table := tables[tableID]

			questionID := data.Get("question").Float()

			table.SelectedQuestion = int(questionID)
			table.QuestionsPlayed = append(table.QuestionsPlayed, int(questionID))

			socket.To(tableID).Emit("state", table.State())
		})

		socket.On("wheel:start", func(data gocket.EmitterData) {
			tableID := socket.Storage["table_id"]
			table := tables[tableID]

			server.Blogger.Info.Println(tableID, table.ranges)

			questionID := 0
			questionID, table.ranges = randPop(table.ranges)

			table.SelectedQuestion = int(questionID)
			table.QuestionsPlayed = append(table.QuestionsPlayed, int(questionID))

			fmt.Println("questionID:", questionID)

			server.GetRoom(tableID).Emit("wheel:spin", &gocket.EmitterData{
				"question":   questionID,
				"spin_count": 3,
			})
		})

		socket.On("timer:start", func(data gocket.EmitterData) {
			tableID := socket.Storage["table_id"]
			table := tables[tableID]

			table.Duration = int(data.Get("duration").Float())
			table.timerStart = time.Now()
			table.TimerRunning = true

			timer := time.NewTimer(time.Duration(table.Duration) * time.Second)
			go func() {
				for {
					select {
					case <-timer.C:
						table.TimerRunning = false
						timer.Stop()
					}
				}
			}()

			server.GetRoom(tableID).Emit("state", table.State())
		})
	})

	server.OnDisconnecting(func(socket *gocket.Socket) {
		if len(socket.Storage) == 0 {
			return
		}

		tableID := socket.Storage["table_id"]
		userID := socket.Storage["user_id"]

		table := tables[tableID]

		user := table.GetUser(userID)
		if user != nil {
			count := 0

			for client := range server.GetRoom(tableID).GetSockets() {
				if client.Storage["user_id"] == userID {
					count++
				}
			}

			if count == 1 {
				user.Online = false
				socket.To(tableID).Emit("state", table.State())
			}
		}
	})

	fmt.Printf("The server is running on port :%s...\n", *Port)

	http.Handle("/ws", server)
	log.Fatal(http.ListenAndServe(":"+*Port, nil))
}
