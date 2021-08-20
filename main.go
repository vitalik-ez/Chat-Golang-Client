package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	signUp = "http://localhost:8000/auth/sign-up"
	signIn = "http://localhost:8000/auth/sign-un"
)

type InputCredentials struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func InputData() *InputCredentials {
	return &InputCredentials{
		Name:     "Vitaliy1",
		Email:    "vetalyeshor345@gmail.com",
		Password: "qwerty",
	}
}

func main() {
	fmt.Println("Start client side")

	input := InputData()

	requestBody, err := json.Marshal(input)

	client := &http.Client{}
	req, err := http.NewRequest(
		"POST", signUp, bytes.NewBuffer(requestBody),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("You successfully register")
	} else {
		io.Copy(os.Stdout, resp.Body)
	}

}
