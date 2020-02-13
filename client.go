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

type Client struct {
	Config    *Config
	SessionId string
}

func NewClient(config *Config) *Client {
	client := Client{
		Config:    config,
		SessionId: getSessionId(config),
	}

	return &client
}

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

func (self *Client) UpdateSessionId() {
	self.SessionId = getSessionId(self.Config)
}

func (self *Client) rpc(requestBody RequestBody) http.Response {
	client := &http.Client{}

	url := getUrl(self.Config)

	jsonData, err := json.Marshal(requestBody)
	panicOnError(err)

	request, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	panicOnError(err)

	login := self.Config.Login
	request.SetBasicAuth(login.Username, login.Password)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Transmission-Session-Id", self.SessionId)

	response, err := client.Do(request)
	panicOnError(err)

	return *response
}

func (self *Client) AddTorrent(link string) {
	var requestBody RequestBody
	requestBody.Method = "torrent-add"
	requestBody.Arguments.Filename = link

	if self.Config.Paused {
		requestBody.Arguments.Paused = true
	}

	self.rpc(requestBody)
}
