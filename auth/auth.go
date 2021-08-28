package auth

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	signUpAddress = "http://localhost:8000/auth/sign-up"
	signInAddress = "http://localhost:8000/auth/sign-un"
)

type InputCredentials struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func inputDataForRegistration() *InputCredentials {
	commands := []string{"Enter a name:", "Enter a email:", "Enter a password:"}
	var inputData []string
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < len(commands); i++ {
		for {
			fmt.Print(commands[i])
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			//validate
			if len(text) == 0 {
				continue
			}
			inputData = append(inputData, text)
			break

		}
	}

	return &InputCredentials{
		Name:     inputData[0],
		Email:    inputData[1],
		Password: inputData[2],
	}
}

func inputDataForLogIn() {

}

type ResponseRegister struct {
	Id uint `json:"id" binding:"required"`
}

func signUp() {
	input := inputDataForRegistration()
	requestBody, err := json.Marshal(input)

	client := &http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(
		http.MethodPost, signUpAddress, bytes.NewBuffer(requestBody),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	register := ResponseRegister{}
	if resp.StatusCode == 200 {
		fmt.Println("You successfully register")
		err := json.Unmarshal(body, &register)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(register.Id)
	} else {
		io.Copy(os.Stdout, resp.Body)
	}

}

func signIn() {

}

func Menu() {
	for {
		fmt.Println("*** MENU ***")
		fmt.Println("1) Sign in.")
		fmt.Println("2) Sign up.")
		fmt.Println("3) Exit.")
		fmt.Print("Enter a number of point:")

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
			signIn()
		case 2:
			fmt.Println("Register")
			signUp()
		default:
			fmt.Println("Exit")
			return
		}
	}

}
