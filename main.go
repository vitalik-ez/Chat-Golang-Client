package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vitalik-ez/Chat-Golang-Client/api"
	"github.com/vitalik-ez/Chat-Golang-Client/room"
)

func initConfig() string {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
	return os.Getenv("SERVER_BASE_PATH")
}

func main() {
	api.StatusServer()
	serverBasePath := initConfig()
	room.Menu(serverBasePath)

}
