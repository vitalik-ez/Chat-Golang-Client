package room

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vitalik-ez/Chat-Golang-Client/auth"
)

const (
	createRoomAddress = "http://localhost:8000/api/room"
	logInRoom         = "ws://localhost:8000/api/room/ws/"
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
			connectToRoom(user)
		default:
			fmt.Println("Exit")
			return
		}
	}
}

type Message struct {
	Room     string    `json:"room" binding:"required"`
	Author   string    `json:"author" binding:"required"`
	Text     string    `json:"text"   binding:"required"`
	CreateAt time.Time `json:"time"   binding:"required"`
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

func connectToRoom(user *auth.User) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	var bearer = "Bearer " + user.Token
	c, _, err := websocket.DefaultDialer.Dial(logInRoom, http.Header{"Authorization": []string{bearer}})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	var listOfRooms []string
	err = c.ReadJSON(&listOfRooms)
	if err != nil {
		log.Println("error while receive list of rooms", err)
		return
	}
	var pointOfMenu uint
	if len(listOfRooms) == 0 {
		fmt.Println("Rooms aren't existed")
		time.Sleep(2 * time.Second)
		return
	} else {
		for index, value := range listOfRooms {
			fmt.Println(index+1, value)
		}
		fmt.Print("Enter number of room: ")
		_, err = fmt.Scanf("%d", &pointOfMenu)
		if err != nil {
			fmt.Println("Invalid pointOfMenu:", err)
			return
		}
		c.WriteJSON(&ServerCommand{
			Command: "join",
			Data:    listOfRooms[pointOfMenu-1],
			Author:  user.Name,
		})
	}

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
			fmt.Println("recv:", message.CreateAt.Format("01-02-2006 15:04:05 Monday"), message.Author, message.Text)
		}
	}()

	inputMessage := make(chan string)
	go readInput(inputMessage)

	for {
		select {
		case <-done:
			return
		case message := <-inputMessage:
			msg := Message{Text: message, CreateAt: time.Now(), Author: user.Name, Room: listOfRooms[pointOfMenu]}
			fmt.Println("send message", msg)
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
					return*/
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
	req, err := http.NewRequest(http.MethodPost, createRoomAddress, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	room := ResponseIdRoom{}
	if resp.StatusCode == http.StatusOK {
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
