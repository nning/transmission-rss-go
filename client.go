package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RequestBody struct {
	Method    string `json:"method"`
	Arguments struct {
		Filename string `json:"filename"`
		Paused   bool   `json:"paused"`
	} `json:"arguments"`
}

// TODO "class"

func getSessionId(config *Config) string {
	client := &http.Client{}

	url := getUrl(config)

	request, err := http.NewRequest("GET", url, nil)
	panicOnError(err)

	request.SetBasicAuth(config.Login.Username, config.Login.Password)
	response, err := client.Do(request)
	panicOnError(err)

	_, err = ioutil.ReadAll(response.Body)
	panicOnError(err)

	sessionId := response.Header["X-Transmission-Session-Id"][0]

	fmt.Println("SESSION_ID " + sessionId)

	return sessionId
}

func rpc(config *Config, sessionId string, requestBody RequestBody) http.Response {
	client := &http.Client{}

	url := getUrl(config)

	if sessionId == "" {
		sessionId = getSessionId(config)
	}

	jsonData, err := json.Marshal(requestBody)
	panicOnError(err)

	request, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	panicOnError(err)

	request.SetBasicAuth(config.Login.Username, config.Login.Password)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Transmission-Session-Id", sessionId)

	response, err := client.Do(request)
	panicOnError(err)

	return *response
}

func addTorrent(config *Config, sessionId string, link string) {
	var requestBody RequestBody
	requestBody.Method = "torrent-add"
	requestBody.Arguments.Filename = link

	if config.Paused {
		requestBody.Arguments.Paused = true
	}

	rpc(config, sessionId, requestBody)
}
