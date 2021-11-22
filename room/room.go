package room

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vitalik-ez/Chat-Golang-Client/entity"
)

type Command string

const (
	join      Command = "join"
	leave     Command = "leave"
	broadcast Command = "broadcast"
	leaveRoom Command = "leave_room"
)

const (
	joinRoomAddress string = "api/room/ws/"
)

func httpToWS(serverBasePath string) string {
	return "ws" + serverBasePath[len("http"):]
}

func Menu(serverBasePath string) {

	serverBasePath = httpToWS(serverBasePath)
	client := NewClient(serverBasePath)

	fmt.Println("Enter your name:")
	client.Name = inputData()

	for {
		fmt.Print("Enter name of room: ")
		room := inputData()
		client.connectToRoom(room)
	}
}

func inputData() string {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()
		if len(text) == 0 {
			fmt.Println("Enter the required information!")
			continue
		}
		return text
	}
}

type Server struct {
	Command  Command `json:"command"`
	Message  string  `json:"message"`
	UserName string  `json:"userName"`
	Room     string  `json:"room"`
}

func readInput(inputMessage chan<- string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, _ := reader.ReadString('\n')
		message = strings.Replace(message, "\n", "", -1)
		if len(message) == 0 {
			continue
		}
		inputMessage <- message
	}
}

func readMessage(c *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		message := entity.NewEmptyMessage()
		err := c.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error: %v", err)
			}
			break
		}
		fmt.Println("recv:", message.CreateAt.Format("01-02-2006 15:04:05 Monday"), message.UserName, ":", message.Text)
	}
}

func (c *Client) connectToRoom(room string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	address := c.Config.ServerBasePath + joinRoomAddress
	conn, _, err := websocket.DefaultDialer.Dial(address, http.Header{})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = conn.WriteJSON(&Server{
		Command:  join,
		Room:     room,
		UserName: c.Name,
	})
	if err != nil {
		log.Println("Error with command to join to the room.", err)
		return
	}

	done := make(chan struct{})
	go readMessage(conn, done)

	inputMessage := make(chan string)
	go readInput(inputMessage)

	for {
		select {
		case <-done:
			return
		case message := <-inputMessage:
			var server *Server
			server = &Server{
				Room:     room,
				UserName: c.Name,
				Message:  message,
			}
			switch message {
			case string(leaveRoom):
				server.Command = leave
			default:
				server.Command = broadcast
			}

			err := conn.WriteJSON(server)
			if err != nil {
				log.Println("write:", err)
				return
			}

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
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
