package main

import (
	"fmt"

	"github.com/vitalik-ez/Chat-Golang-Client/auth"
	"github.com/vitalik-ez/Chat-Golang-Client/room"
)

func Menu() {
	user := auth.Menu()
	room.RoomMenu(user)

}

func main() {
	fmt.Println("Start client side")
	Menu()

}
