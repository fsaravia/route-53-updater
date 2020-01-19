package ipResolver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	apiKeyParam = "x-api-key"
	getVerb     = "GET"
)

type response struct {
	Ip string `json:"ip_address"`
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return "", errors.New(fmt.Sprintf("Please set the %s environment variable", key))
	}

	return value, nil
}

func resolverURL() (string, error) {
	return getEnv("RESOLVER_URL")
}

func resolverAPIKey() (string, error) {
	return getEnv("API_KEY")
}

func handleError(err error) (string, error) {
	return "", err
}

func ResolveIp() (string, error) {
	client := &http.Client{}

	resolverURL, err := resolverURL()
	if err != nil {
		return handleError(err)
	}

	req, err := http.NewRequest(getVerb, resolverURL, nil)
	if err != nil {
		return handleError(err)
	}

	apiKey, err := resolverAPIKey()
	if err != nil {
		return handleError(err)
	}

	req.Header.Add(apiKeyParam, apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return handleError(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return handleError(err)
	}

	response := &response{}

	err = json.Unmarshal(body, response)
	if err != nil {
		handleError(err)
	}

	return response.Ip, nil
}
