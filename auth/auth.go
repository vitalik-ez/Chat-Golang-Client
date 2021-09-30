package auth

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vitalik-ez/Chat-Golang-Client/api"
)

const (
	signUpAddress = "http://localhost:8000/auth/sign-up"
	signInAddress = "http://localhost:8000/auth/sign-in"
)

func getCommandsForSignUp() []string {
	return []string{"name", "email", "password"}
}

func getCommandsForSignIn() []string {
	return []string{"email", "password"}
}

func inputData(commands []string) map[string]string {
	inputData := make(map[string]string)
	reader := bufio.NewReader(os.Stdin)
	for _, command := range commands {
		for {
			fmt.Printf("Enter %s:", command)
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			if len(text) == 0 {
				continue
			}
			inputData[command] = text
			break
		}
	}
	return inputData
}

type ResponseRegister struct {
	Id uint `json:"id" binding:"required"`
}

func signUp() {
	input := inputData(getCommandsForSignUp())
	requestBody, err := json.Marshal(input)
	if err != nil {
		log.Fatal(err)
	}
	body := api.PostRequest(signUpAddress, requestBody)
	register := ResponseRegister{}
	if body != nil {
		err := json.Unmarshal(body, &register)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println("You successfully register")
	} else {
		fmt.Println("Server error. Try again!")
	}
}

type User struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}

func signIn() *User {
	input := inputData(getCommandsForSignIn())
	requestBody, err := json.Marshal(input)
	if err != nil {
		log.Fatal(err)
	}
	body := api.PostRequest(signInAddress, requestBody)
	user := User{}
	if body != nil {
		err := json.Unmarshal(body, &user)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println("You successfully log in")
	} else {
		fmt.Println("Server error. Try again!")
		return nil
	}
	return &user
}

func Menu() *User {
	for {
		fmt.Println("*** MENU ***")
		fmt.Println("1) Sign in.")
		fmt.Println("2) Sign up.")
		fmt.Println("3) Exit.")
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
			fmt.Println("Authorization")
			if user := signIn(); user != nil {
				return user
			}
		case 2:
			fmt.Println("Register")
			signUp()
		default:
			fmt.Println("Exit")
		}
	}

}
