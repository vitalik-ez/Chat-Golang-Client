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

/*
func clientHttpRequest(request *http.Request) []byte {
	client := &http.Client{
		Timeout: time.Second * 2,
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	return body
}

func GetRequest(address string, token ...string) []byte {
	req, err := http.NewRequest(
		http.MethodGet, address, nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	return clientHttpRequest(req)
}
*/
