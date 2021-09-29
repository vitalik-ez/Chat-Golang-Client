package api

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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
	client := &http.Client{
		Timeout: time.Second * 2,
	}
	resp, err := client.Do(req)
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
