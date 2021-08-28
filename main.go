package main

import (
	"fmt"

	"github.com/vitalik-ez/Chat-Golang-Client/auth"
)

func Menu() {
	auth.Menu()
}

func main() {
	fmt.Println("Start client side")
	Menu()

}
