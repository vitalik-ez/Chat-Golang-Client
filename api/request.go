package api

import (
	"log"
	"net/http"
	"time"
)

const (
	statusServerAddress = "http://localhost:8000/status-server"
)

func StatusServer() {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(statusServerAddress)
	if err != nil {
		log.Fatal("Server doesn't work. ", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Status code isn't correct.")
	}
}
