package room

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vitalik-ez/Chat-Golang-Client/api"
	"github.com/vitalik-ez/Chat-Golang-Client/auth"
	"github.com/vitalik-ez/Chat-Golang-Client/entity"
)

const (
	createRoomAddress     = "http://localhost:8000/api/room"
	getListOfRoomsAddress = "http://localhost:8000/api/room"
	logInRoomAddress      = "ws://localhost:8000/api/room/ws/"
	join                  = "join"
	leave                 = "leave"
	broadcast             = "broadcast"
)

func RoomMenu(user *auth.User) {

	for {
		fmt.Println("Room Menu")
		fmt.Println("1. Create room")
		fmt.Println("2. Enter in exist room")
		fmt.Println("3. Exit")
		fmt.Print("Enter a number of point: ")
		var pointOfMenu uint8
		for {
			_, err := fmt.Scanf("%d", &pointOfMenu)
			if err != nil {
				fmt.Println("Invalid pointOfMenu:", err)
				continue
			}
			break
		}

		switch pointOfMenu {
		case 1:
			fmt.Println("Create Room:")
			room := createRoom(user.Token)
			if room != nil {
				fmt.Println("The id of created room: ", room.Id)
			}
			//connectToRoom(user)
			//return roomId
		case 2:
			fmt.Println("Enter in exist room")
			if room := getListOfRooms(user.Token); room != "" {
				connectToRoom(user, room)
			}
		default:
			fmt.Println("Exit")
			return
		}
	}
}

type listOfRoom struct {
	List []string `json:"list" binding:"required"`
}

func getListOfRooms(token string) string {
	body := api.GetRequest(getListOfRoomsAddress, "Bearer "+token)
	rooms := listOfRoom{}
	if body != nil {
		err := json.Unmarshal(body, &rooms)
		fmt.Println("Rooms", rooms)
		if err != nil {
			log.Fatal(err.Error())
		}
		if len(rooms.List) > 0 {
			for index, room := range rooms.List {
				fmt.Printf("%d. %s \n", index+1, room)
			}
			//
			var numberOfRoom uint
			fmt.Print("Enter number of room: ")
			_, err = fmt.Scanf("%d", &numberOfRoom)
			if err != nil {
				fmt.Println("Invalid number of room:", err)
				// handle error
			}
			return rooms.List[numberOfRoom-1]
		} else {
			fmt.Println("Rooms aren't existed")
			time.Sleep(2 * time.Second)
		}
	} else {
		fmt.Println("Server error. Try aganin!")
	}
	return ""
}

func readInput(inputMessage chan<- string) {
	for {
		var message string
		_, err := fmt.Scanln(&message)
		if err != nil {
			panic(err)
		}
		inputMessage <- message
	}

}

type ServerCommand struct {
	Command string `json:"command"`
	Data    string `json:"data"`
	Author  string `json:"author"`
}

type ListRoom struct {
	List string `json:"list" binding:"required"`
}

func readMessage(c *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		message := entity.NewEmptyMessage()
		err := c.ReadJSON(&message)
		if err != nil {
			log.Println("read:", err)
			return
		}
		fmt.Println("recv:", message.CreateAt.Format("01-02-2006 15:04:05 Monday"), message.Author, message.Text)
	}
}

func connectToRoom(user *auth.User, room string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	var bearer = "Bearer " + user.Token
	c, _, err := websocket.DefaultDialer.Dial(logInRoomAddress, http.Header{"Authorization": []string{bearer}})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	err = c.WriteJSON(&ServerCommand{
		Command: join,
		Data:    room,
		Author:  user.Name,
	})
	if err != nil {
		log.Println("Error with command to join to the room.", err)
		return
	}

	done := make(chan struct{})
	go readMessage(c, done)
	/*go func() {
		defer close(done)
		for {
			message := entity.NewEmptyMessage()
			err := c.ReadJSON(&message)
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Println("recv:", message.CreateAt.Format("01-02-2006 15:04:05 Monday"), message.Author, message.Text)
		}
	}()*/

	inputMessage := make(chan string)
	go readInput(inputMessage)

	for {
		select {
		case <-done:
			return
		case message := <-inputMessage:
			msg := entity.NewMessage(room, user.Name, message)
			err := c.WriteJSON(msg)
			if err != nil {
				log.Println("write:", err)
				return
			}
			/*
				case <-interrupt:
					log.Println("interrupt")

					// Cleanly close the connection by sending a close message and then
					// waiting (with timeout) for the server to close the connection.
					err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					if err != nil {
						log.Println("write close:", err)
						return
					}
					select {
					case <-done:
					case <-time.After(time.Second):
					}
					return
			*/
		}
	}
}

type InputRoomData struct {
	Name string `json:"name"`
}

type ResponseIdRoom struct {
	Id int `json:"id"`
}

func inputData() *InputRoomData {
	fmt.Print("Enter name room: ")
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if len(text) == 0 {
			continue
		}
		return &InputRoomData{
			Name: text,
		}
	}
}

func createRoom(userToken string) *ResponseIdRoom {
	input := inputData()
	requestBody, err := json.Marshal(input)
	if err != nil {
		log.Fatal(err)
	}
	var bearer = "Bearer " + userToken
	body := api.PostRequest(createRoomAddress, requestBody, bearer)
	room := ResponseIdRoom{}
	if body != nil {
		err := json.Unmarshal(body, &room)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println("You successfully create room")
	} else {
		fmt.Println("Room wasn't created. Try aganin!")
		return nil
	}
	return &room
}
