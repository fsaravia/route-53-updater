package ipresolver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	apiKeyParam = "x-api-key"
	getVerb     = "GET"
)

type response struct {
	Ip string `json:"ip_address"`
}

func handleError(err error) (string, error) {
	return "", err
}

func ResolveIp(resolverURL string, resolverAPIKey string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest(getVerb, resolverURL, nil)
	if err != nil {
		return handleError(err)
	}

	req.Header.Add(apiKeyParam, resolverAPIKey)

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
