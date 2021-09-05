package room

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vitalik-ez/Chat-Golang-Client/auth"
)

const (
	createRoomAddress = "http://localhost:8000/api/room"
	logInRoom         = "ws://localhost:8000/api/room/ws/"
)

func RoomMenu(user *auth.User) uint {
	roomId := createRoom(user.Token)
	roomAddress := fmt.Sprintf("%s%d", logInRoom, roomId)
	connectToRoom(user, roomAddress)
	return roomId
}

type Message struct {
	Data     string    `json:"data"`
	CreateAt time.Time `json:"time"`
	Author   string    `json:"author"`
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

func connectToRoom(user *auth.User, roomAddress string) {
	fmt.Println(user.Name, " connect to the room")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var bearer = "Bearer " + user.Token
	c, _, err := websocket.DefaultDialer.Dial(roomAddress, http.Header{"Authorization": []string{bearer}})
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			message := Message{}
			err := c.ReadJSON(&message)
			if err != nil {
				log.Println("read:", err)
				return
			}
			//log.Println("recv:", message.CreateAt, message.Author, message.Data)
			fmt.Println("recv:", message.CreateAt.Format("01-02-2006 15:04:05 Monday"), message.Author, message.Data)
		}
	}()

	inputMessage := make(chan string)
	go readInput(inputMessage)

	for {
		select {
		case <-done:
			return
		case message := <-inputMessage:
			msg := Message{Data: message, CreateAt: time.Now(), Author: user.Name}
			err := c.WriteJSON(msg)
			if err != nil {
				log.Println("write:", err)
				return
			}
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
		}
	}
}

type InputRoomData struct {
	Name string `json:"name"`
}

type ResponseIdRoom struct {
	Id uint `json:"id"`
}

func inputData() *InputRoomData {
	return &InputRoomData{
		Name: "KPI114654",
	}
}

func createRoom(userToken string) uint {
	input := inputData()
	requestBody, err := json.Marshal(input)

	if err != nil {
		log.Fatal(err)
	}

	var bearer = "Bearer " + userToken

	// Create a new request using http
	req, err := http.NewRequest(http.MethodPost, createRoomAddress, bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatal(err)
	}
	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	room := ResponseIdRoom{}

	if resp.StatusCode == 200 {
		fmt.Println("You successfully create or log in room")
		err := json.Unmarshal(body, &room)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		log.Fatal("Error Room id !!!")
		io.Copy(os.Stdout, resp.Body)
	}
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println(string([]byte(body)))

	return room.Id

}
