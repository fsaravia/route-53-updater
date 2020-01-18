package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const apiKeyParam = "x-api-key"

type response struct {
	Ip string `json:"ip_address"`
}

func obtainIPAddress () (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", os.Getenv("URL"), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add(apiKeyParam, os.Getenv("API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	response := &response{}

	err = json.Unmarshal(body, response)
	if err != nil {
		return "", err
	}

	return response.Ip, nil
}

func main() {
	ip, err := obtainIPAddress()
	if err != nil {
		panic(err)
	}

	fmt.Printf("The IP address of this device is: %s", ip)
}
