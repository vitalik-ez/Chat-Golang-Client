package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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
	fmt.Println(body, resp.StatusCode)
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
	if len(token) > 0 {
		req.Header.Add("Authorization", token[0])
	}
	return clientHttpRequest(req)
}

func PostRequest(address string, requestBody []byte, token ...string) []byte {
	req, err := http.NewRequest(
		http.MethodPost, address, bytes.NewBuffer(requestBody),
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if len(token) > 0 {
		req.Header.Add("Authorization", token[0])
	}
	return clientHttpRequest(req)
}
